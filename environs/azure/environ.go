// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package azure

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"launchpad.net/gwacl"

	agenttools "launchpad.net/juju-core/agent/tools"
	"launchpad.net/juju-core/constraints"
	"launchpad.net/juju-core/environs"
	"launchpad.net/juju-core/environs/cloudinit"
	"launchpad.net/juju-core/environs/config"
	"launchpad.net/juju-core/environs/instances"
	"launchpad.net/juju-core/environs/tools"
	"launchpad.net/juju-core/instance"
	"launchpad.net/juju-core/state"
	"launchpad.net/juju-core/state/api"
	"launchpad.net/juju-core/utils/parallel"
)

const (
	// In our initial implementation, each instance gets its own hosted
	// service, deployment and role in Azure.  The role always gets this
	// hostname (instance==service).
	roleHostname = "default"

	// deploymentSlot says in which slot to deploy instances.  Azure
	// supports 'Production' or 'Staging'.
	// This provider always deploys to Production.  Think twice about
	// changing that: DNS names in the staging slot work differently from
	// those in the production slot.  In Staging, Azure assigns an
	// arbitrary hostname that we can then extract from the deployment's
	// URL.  In Production, the hostname in the deployment URL does not
	// actually seem to resolve; instead, the service name is used as the
	// DNS name, with ".cloudapp.net" appended.
	deploymentSlot = "Production"

	// Address space of the virtual network used by the nodes in this
	// environement, in CIDR notation. This is the network used for
	// machine-to-machine communication.
	networkDefinition = "10.0.0.0/8"
)

type azureEnviron struct {
	// Except where indicated otherwise, all fields in this object should
	// only be accessed using a lock or a snapshot.
	sync.Mutex

	// name is immutable; it does not need locking.
	name string

	// ecfg is the environment's Azure-specific configuration.
	ecfg *azureEnvironConfig

	// storage is this environ's own private storage.
	storage environs.Storage

	// publicStorage is the public storage that this environ uses.
	publicStorage environs.StorageReader

	// storageAccountKey holds an access key to this environment's
	// private storage.  This is automatically queried from Azure on
	// startup.
	storageAccountKey string
}

// azureEnviron implements Environ.
var _ environs.Environ = (*azureEnviron)(nil)

// NewEnviron creates a new azureEnviron.
func NewEnviron(cfg *config.Config) (*azureEnviron, error) {
	env := azureEnviron{name: cfg.Name()}
	err := env.SetConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Set up storage.
	env.storage = &azureStorage{
		storageContext: &environStorageContext{environ: &env},
	}

	// Set up public storage.
	publicContext := publicEnvironStorageContext{environ: &env}
	if publicContext.getContainer() == "" {
		// No public storage configured.  Use EmptyStorage.
		env.publicStorage = environs.EmptyStorage
	} else {
		// Set up real public storage.
		env.publicStorage = &azureStorage{storageContext: &publicContext}
	}

	return &env, nil
}

// extractStorageKey returns the primary account key from a gwacl
// StorageAccountKeys struct, or if there is none, the secondary one.
func extractStorageKey(keys *gwacl.StorageAccountKeys) string {
	if keys.Primary != "" {
		return keys.Primary
	}
	return keys.Secondary
}

// queryStorageAccountKey retrieves the storage account's key from Azure.
func (env *azureEnviron) queryStorageAccountKey() (string, error) {
	azure, err := env.getManagementAPI()
	if err != nil {
		return "", err
	}
	defer env.releaseManagementAPI(azure)

	accountName := env.getSnapshot().ecfg.storageAccountName()
	keys, err := azure.GetStorageAccountKeys(accountName)
	if err != nil {
		return "", fmt.Errorf("cannot obtain storage account keys: %v", err)
	}

	key := extractStorageKey(keys)
	if key == "" {
		return "", fmt.Errorf("no keys available for storage account")
	}

	return key, nil
}

// Name is specified in the Environ interface.
func (env *azureEnviron) Name() string {
	return env.name
}

// getSnapshot produces an atomic shallow copy of the environment object.
// Whenever you need to access the environment object's fields without
// modifying them, get a snapshot and read its fields instead.  You will
// get a consistent view of the fields without any further locking.
// If you do need to modify the environment's fields, do not get a snapshot
// but lock the object throughout the critical section.
func (env *azureEnviron) getSnapshot() *azureEnviron {
	env.Lock()
	defer env.Unlock()

	// Copy the environment.  (Not the pointer, the environment itself.)
	// This is a shallow copy.
	snap := *env
	// Reset the snapshot's mutex, because we just copied it while we
	// were holding it.  The snapshot will have a "clean," unlocked mutex.
	snap.Mutex = sync.Mutex{}
	return &snap
}

// startBootstrapInstance starts the bootstrap instance for this environment.
func (env *azureEnviron) startBootstrapInstance(cons constraints.Value) (instance.Instance, error) {
	// The bootstrap instance gets machine id "0".  This is not related to
	// instance ids or anything in Azure.  Juju assigns the machine ID.
	const machineID = "0"

	// Create an empty bootstrap state file so we can get its URL.
	// It will be updated with the instance id and hardware characteristics
	// after the bootstrap instance is started.
	stateFileURL, err := environs.CreateStateFile(env.Storage())
	if err != nil {
		return nil, err
	}
	machineConfig := environs.NewBootstrapMachineConfig(machineID, stateFileURL)

	logger.Debugf("bootstrapping environment %q", env.Name())
	possibleTools, err := tools.FindBootstrapTools(env, cons)
	if err != nil {
		return nil, err
	}
	inst, err := env.internalStartInstance(cons, possibleTools, machineConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot start bootstrap instance: %v", err)
	}
	return inst, nil
}

// getAffinityGroupName returns the name of the affinity group used by all
// the Services in this environment.
func (env *azureEnviron) getAffinityGroupName() string {
	return env.getEnvPrefix() + "ag"
}

func (env *azureEnviron) createAffinityGroup() error {
	affinityGroupName := env.getAffinityGroupName()
	azure, err := env.getManagementAPI()
	if err != nil {
		return err
	}
	defer env.releaseManagementAPI(azure)
	snap := env.getSnapshot()
	location := snap.ecfg.location()
	cag := gwacl.NewCreateAffinityGroup(affinityGroupName, affinityGroupName, affinityGroupName, location)
	return azure.CreateAffinityGroup(&gwacl.CreateAffinityGroupRequest{
		CreateAffinityGroup: cag})
}

func (env *azureEnviron) deleteAffinityGroup() error {
	affinityGroupName := env.getAffinityGroupName()
	azure, err := env.getManagementAPI()
	if err != nil {
		return err
	}
	defer env.releaseManagementAPI(azure)
	return azure.DeleteAffinityGroup(&gwacl.DeleteAffinityGroupRequest{
		Name: affinityGroupName})
}

// getVirtualNetworkName returns the name of the virtual network used by all
// the VMs in this environment.
func (env *azureEnviron) getVirtualNetworkName() string {
	return env.getEnvPrefix() + "vnet"
}

func (env *azureEnviron) createVirtualNetwork() error {
	vnetName := env.getVirtualNetworkName()
	affinityGroupName := env.getAffinityGroupName()
	azure, err := env.getManagementAPI()
	if err != nil {
		return err
	}
	defer env.releaseManagementAPI(azure)
	virtualNetwork := gwacl.VirtualNetworkSite{
		Name:          vnetName,
		AffinityGroup: affinityGroupName,
		AddressSpacePrefixes: []string{
			networkDefinition,
		},
	}
	return azure.AddVirtualNetworkSite(&virtualNetwork)
}

func (env *azureEnviron) deleteVirtualNetwork() error {
	azure, err := env.getManagementAPI()
	if err != nil {
		return err
	}
	defer env.releaseManagementAPI(azure)
	vnetName := env.getVirtualNetworkName()
	return azure.RemoveVirtualNetworkSite(vnetName)
}

// getContainerName returns the name of the private storage account container
// that this environment is using.
func (env *azureEnviron) getContainerName() string {
	return env.getEnvPrefix() + "private"
}

// Bootstrap is specified in the Environ interface.
// TODO(bug 1199847): This work can be shared between providers.
func (env *azureEnviron) Bootstrap(cons constraints.Value) (err error) {
	// TODO(bug 1199847). The creation of the affinity group and the
	// virtual network is specific to the Azure provider.
	err = env.createAffinityGroup()
	if err != nil {
		return err
	}
	// If we fail after this point, clean up the affinity group.
	defer func() {
		if err != nil {
			env.deleteAffinityGroup()
		}
	}()
	err = env.createVirtualNetwork()
	if err != nil {
		return err
	}
	// If we fail after this point, clean up the virtual network.
	defer func() {
		if err != nil {
			env.deleteVirtualNetwork()
		}
	}()

	inst, err := env.startBootstrapInstance(cons)
	if err != nil {
		return err
	}
	// TODO(wallyworld) - save hardware characteristics
	err = environs.SaveState(
		env.Storage(),
		&environs.BootstrapState{StateInstances: []instance.Id{inst.Id()}})
	if err != nil {
		err2 := env.StopInstances([]instance.Instance{inst})
		if err2 != nil {
			// Failure upon failure.  Log it, but return the
			// original error.
			logger.Errorf("cannot release failed bootstrap instance: %v", err2)
		}
		return fmt.Errorf("cannot save state: %v", err)
	}

	// TODO make safe in the case of racing Bootstraps
	// If two Bootstraps are called concurrently, there's
	// no way to make sure that only one succeeds.
	return nil
}

// StateInfo is specified in the Environ interface.
func (env *azureEnviron) StateInfo() (*state.Info, *api.Info, error) {
	return environs.StateInfo(env)
}

// Config is specified in the Environ interface.
func (env *azureEnviron) Config() *config.Config {
	snap := env.getSnapshot()
	return snap.ecfg.Config
}

// SetConfig is specified in the Environ interface.
func (env *azureEnviron) SetConfig(cfg *config.Config) error {
	ecfg, err := azureEnvironProvider{}.newConfig(cfg)
	if err != nil {
		return err
	}

	env.Lock()
	defer env.Unlock()

	if env.ecfg != nil {
		_, err = azureEnvironProvider{}.Validate(cfg, env.ecfg.Config)
		if err != nil {
			return err
		}
	}

	env.ecfg = ecfg

	// Reset storage account key.  Even if we had one before, it may not
	// be appropriate for the new config.
	env.storageAccountKey = ""

	return nil
}

// attemptCreateService tries to create a new hosted service on Azure, with a
// name it chooses (based on the given prefix), but recognizes that the name
// may not be available.  If the name is not available, it does not treat that
// as an error but just returns nil.
func attemptCreateService(azure *gwacl.ManagementAPI, prefix string, affinityGroupName string, location string) (*gwacl.CreateHostedService, error) {
	var err error
	name := gwacl.MakeRandomHostedServiceName(prefix)
	err = azure.CheckHostedServiceNameAvailability(name)
	if err != nil {
		// The calling function should retry.
		return nil, nil
	}
	req := gwacl.NewCreateHostedServiceWithLocation(name, name, location)
	req.AffinityGroup = affinityGroupName
	err = azure.AddHostedService(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// architectures lists the CPU architectures supported by Azure.
var architectures = []string{"amd64", "i386"}

// newHostedService creates a hosted service.  It will make up a unique name,
// starting with the given prefix.
func newHostedService(azure *gwacl.ManagementAPI, prefix string, affinityGroupName string, location string) (*gwacl.CreateHostedService, error) {
	var err error
	var svc *gwacl.CreateHostedService
	for tries := 10; tries > 0 && err == nil && svc == nil; tries-- {
		svc, err = attemptCreateService(azure, prefix, affinityGroupName, location)
	}
	if err != nil {
		return nil, fmt.Errorf("could not create hosted service: %v", err)
	}
	if svc == nil {
		return nil, fmt.Errorf("could not come up with a unique hosted service name - is your randomizer initialized?")
	}
	return svc, nil
}

// selectInstanceTypeAndImage returns the appropriate instance-type name and
// the OS image name for launching a virtual machine with the given parameters.
func (env *azureEnviron) selectInstanceTypeAndImage(cons constraints.Value, series, location string) (string, string, error) {
	ecfg := env.getSnapshot().ecfg
	sourceImageName := ecfg.forceImageName()
	if sourceImageName != "" {
		// Configuration forces us to use a specific image.  There may
		// not be a suitable image in the simplestreams database.
		// This means we can't use Juju's normal selection mechanism,
		// because it combines instance-type and image selection: if
		// there are no images we can use, it won't offer us an
		// instance type either.
		//
		// Select the instance type using simple, Azure-specific code.
		machineType, err := selectMachineType(gwacl.RoleSizes, defaultToBaselineSpec(cons))
		if err != nil {
			return "", "", err
		}
		return machineType.Name, sourceImageName, nil
	}

	// Choose the most suitable instance type and OS image, based on
	// simplestreams information.
	//
	// This should be the normal execution path.  The user is not expected
	// to configure a source image name in normal use.
	constraint := instances.InstanceConstraint{
		Region:      location,
		Series:      series,
		Arches:      architectures,
		Constraints: cons,
	}
	spec, err := findInstanceSpec(env, ecfg.imageStream(), constraint)
	if err != nil {
		return "", "", err
	}
	return spec.InstanceType.Id, spec.Image.Id, nil
}

// internalStartInstance does the provider-specific work of starting an
// instance.  The code in StartInstance is actually largely agnostic across
// the EC2/OpenStack/MAAS/Azure providers.
// The instance will be set up for the same series for which you pass tools.
// All tools in possibleTools must be for the same series.
// machineConfig will be filled out with further details, but should contain
// MachineID, MachineNonce, StateInfo, and APIInfo.
// TODO(bug 1199847): Some of this work can be shared between providers.
func (env *azureEnviron) internalStartInstance(cons constraints.Value, possibleTools agenttools.List, machineConfig *cloudinit.MachineConfig) (_ instance.Instance, err error) {
	// Declaring "err" in the function signature so that we can "defer"
	// any cleanup that needs to run during error returns.

	series := possibleTools.Series()
	if len(series) != 1 {
		panic(fmt.Errorf("should have gotten tools for one series, got %v", series))
	}

	err = environs.FinishMachineConfig(machineConfig, env.Config(), cons)
	if err != nil {
		return nil, err
	}

	// Pick tools.  Needed for the custom data (which is what we normally
	// call userdata).
	machineConfig.Tools = possibleTools[0]
	logger.Infof("picked tools %q", machineConfig.Tools)

	// Compose userdata.
	userData, err := makeCustomData(machineConfig)
	if err != nil {
		return nil, fmt.Errorf("custom data: %v", err)
	}

	azure, err := env.getManagementAPI()
	if err != nil {
		return nil, err
	}
	defer env.releaseManagementAPI(azure)

	snap := env.getSnapshot()
	location := snap.ecfg.location()
	service, err := newHostedService(azure.ManagementAPI, env.getEnvPrefix(), env.getAffinityGroupName(), location)
	if err != nil {
		return nil, err
	}
	serviceName := service.ServiceName

	// If we fail after this point, clean up the hosted service.
	defer func() {
		if err != nil {
			azure.DestroyHostedService(
				&gwacl.DestroyHostedServiceRequest{
					ServiceName: serviceName,
				})
		}
	}()

	instanceType, sourceImageName, err := env.selectInstanceTypeAndImage(cons, series[0], location)
	if err != nil {
		return nil, err
	}

	// virtualNetworkName is the virtual network to which all the
	// deployments in this environment belong.
	virtualNetworkName := env.getVirtualNetworkName()

	// 1. Create an OS Disk.
	vhd := env.newOSDisk(sourceImageName)

	// 2. Create a Role for a Linux machine.
	role := env.newRole(instanceType, vhd, userData, roleHostname)

	// 3. Create the Deployment object.
	deployment := env.newDeployment(role, serviceName, serviceName, virtualNetworkName)

	err = azure.AddDeployment(deployment, serviceName)
	if err != nil {
		return nil, err
	}

	var inst instance.Instance

	// From here on, remember to shut down the instance before returning
	// any error.
	defer func() {
		if err != nil && inst != nil {
			err2 := env.StopInstances([]instance.Instance{inst})
			if err2 != nil {
				// Failure upon failure.  Log it, but return
				// the original error.
				logger.Errorf("error releasing failed instance: %v", err)
			}
		}
	}()

	// Assign the returned instance to 'inst' so that the deferred method
	// above can perform its check.
	inst, err = env.getInstance(serviceName)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

// getInstance returns an up-to-date version of the instance with the given
// name.
func (env *azureEnviron) getInstance(instanceName string) (instance.Instance, error) {
	context, err := env.getManagementAPI()
	if err != nil {
		return nil, err
	}
	defer env.releaseManagementAPI(context)
	service, err := context.GetHostedServiceProperties(instanceName, false)
	if err != nil {
		return nil, fmt.Errorf("could not get instance %q: %v", instanceName, err)
	}
	instance := &azureInstance{service.HostedServiceDescriptor, env}
	return instance, nil
}

// newOSDisk creates a gwacl.OSVirtualHardDisk object suitable for an
// Azure Virtual Machine.
func (env *azureEnviron) newOSDisk(sourceImageName string) *gwacl.OSVirtualHardDisk {
	vhdName := gwacl.MakeRandomDiskName("juju")
	vhdPath := fmt.Sprintf("vhds/%s", vhdName)
	snap := env.getSnapshot()
	storageAccount := snap.ecfg.storageAccountName()
	mediaLink := gwacl.CreateVirtualHardDiskMediaLink(storageAccount, vhdPath)
	// The disk label is optional and the disk name can be omitted if
	// mediaLink is provided.
	return gwacl.NewOSVirtualHardDisk("", "", "", mediaLink, sourceImageName, "Linux")
}

// getInitialEndpoints returns a slice of the endpoints every instance should have open
// (ssh port, etc).
func (env *azureEnviron) getInitialEndpoints() []gwacl.InputEndpoint {
	config := env.Config()
	return []gwacl.InputEndpoint{
		{
			LocalPort: 22,
			Name:      "sshport",
			Port:      22,
			Protocol:  "tcp",
		},
		// TODO: Ought to have this only for state servers.
		{
			LocalPort: config.StatePort(),
			Name:      "stateport",
			Port:      config.StatePort(),
			Protocol:  "tcp",
		},
		// TODO: Ought to have this only for API servers.
		{
			LocalPort: config.APIPort(),
			Name:      "apiport",
			Port:      config.APIPort(),
			Protocol:  "tcp",
		}}
}

// newRole creates a gwacl.Role object (an Azure Virtual Machine) which uses
// the given Virtual Hard Drive.
//
// The VM will have:
// - an 'ubuntu' user defined with an unguessable (randomly generated) password
// - its ssh port (TCP 22) open
// - its state port (TCP mongoDB) port open
// - its API port (TCP) open
//
// roleSize is the name of one of Azure's machine types, e.g. ExtraSmall,
// Large, A6 etc.
func (env *azureEnviron) newRole(roleSize string, vhd *gwacl.OSVirtualHardDisk, userData string, roleHostname string) *gwacl.Role {
	// Create a Linux Configuration with the username and the password
	// empty and disable SSH with password authentication.
	hostname := roleHostname
	username := "ubuntu"
	password := gwacl.MakeRandomPassword()
	linuxConfigurationSet := gwacl.NewLinuxProvisioningConfigurationSet(hostname, username, password, userData, "true")
	// Generate a Network Configuration with the initially required ports
	// open.
	networkConfigurationSet := gwacl.NewNetworkConfigurationSet(env.getInitialEndpoints(), nil)
	roleName := gwacl.MakeRandomRoleName("juju")
	// The ordering of these configuration sets is significant for the tests.
	return gwacl.NewRole(
		roleSize, roleName,
		[]gwacl.ConfigurationSet{*linuxConfigurationSet, *networkConfigurationSet},
		[]gwacl.OSVirtualHardDisk{*vhd})
}

// newDeployment creates and returns a gwacl Deployment object.
func (env *azureEnviron) newDeployment(role *gwacl.Role, deploymentName string, deploymentLabel string, virtualNetworkName string) *gwacl.Deployment {
	// Use the service name as the label for the deployment.
	return gwacl.NewDeploymentForCreateVMDeployment(deploymentName, deploymentSlot, deploymentLabel, []gwacl.Role{*role}, virtualNetworkName)
}

// StartInstance is specified in the Environ interface.
// TODO(bug 1199847): This work can be shared between providers.
func (env *azureEnviron) StartInstance(machineID, machineNonce string, series string, cons constraints.Value,
	stateInfo *state.Info, apiInfo *api.Info) (instance.Instance, *instance.HardwareCharacteristics, error) {
	possibleTools, err := tools.FindInstanceTools(env, series, cons)
	if err != nil {
		return nil, nil, err
	}
	err = tools.CheckToolsSeries(possibleTools, series)
	if err != nil {
		return nil, nil, err
	}
	machineConfig := environs.NewMachineConfig(machineID, machineNonce, stateInfo, apiInfo)
	// TODO(bug 1193998) - return instance hardware characteristics as well.
	inst, err := env.internalStartInstance(cons, possibleTools, machineConfig)
	return inst, nil, err
}

// Spawn this many goroutines to issue requests for destroying services.
// TODO: this is currently set to 1 because of a problem in Azure:
// removing Services in the same affinity group concurrently causes a conflict.
// This conflict is wrongly reported by Azure as a BadRequest (400).
// This has been reported to Windows Azure.
var maxConcurrentDeletes = 1

// StopInstances is specified in the Environ interface.
func (env *azureEnviron) StopInstances(instances []instance.Instance) error {
	// Each Juju instance is an Azure Service (instance==service), destroy
	// all the Azure services.
	// Acquire management API object.
	context, err := env.getManagementAPI()
	if err != nil {
		return err
	}
	defer env.releaseManagementAPI(context)

	// Destroy all the services in parallel.
	run := parallel.NewRun(maxConcurrentDeletes)
	for _, instance := range instances {
		serviceName := string(instance.Id())
		run.Do(func() error {
			request := &gwacl.DestroyHostedServiceRequest{ServiceName: serviceName}
			return context.DestroyHostedService(request)
		})
	}
	return run.Wait()
}

// Instances is specified in the Environ interface.
func (env *azureEnviron) Instances(ids []instance.Id) ([]instance.Instance, error) {
	// The instance list is built using the list of all the relevant
	// Azure Services (instance==service).
	// Acquire management API object.
	context, err := env.getManagementAPI()
	if err != nil {
		return nil, err
	}
	defer env.releaseManagementAPI(context)

	// Prepare gwacl request object.
	serviceNames := make([]string, len(ids))
	for i, id := range ids {
		serviceNames[i] = string(id)
	}
	request := &gwacl.ListSpecificHostedServicesRequest{ServiceNames: serviceNames}

	// Issue 'ListSpecificHostedServices' request with gwacl.
	services, err := context.ListSpecificHostedServices(request)
	if err != nil {
		return nil, err
	}

	// If no instances were found, return ErrNoInstances.
	if len(services) == 0 {
		return nil, environs.ErrNoInstances
	}

	instances := convertToInstances(services, env)

	// Check if we got a partial result.
	if len(ids) != len(instances) {
		return instances, environs.ErrPartialInstances
	}
	return instances, nil
}

// AllInstances is specified in the Environ interface.
func (env *azureEnviron) AllInstances() ([]instance.Instance, error) {
	// The instance list is built using the list of all the Azure
	// Services (instance==service).
	// Acquire management API object.
	context, err := env.getManagementAPI()
	if err != nil {
		return nil, err
	}
	defer env.releaseManagementAPI(context)

	request := &gwacl.ListPrefixedHostedServicesRequest{ServiceNamePrefix: env.getEnvPrefix()}
	services, err := context.ListPrefixedHostedServices(request)
	if err != nil {
		return nil, err
	}
	return convertToInstances(services, env), nil
}

// getEnvPrefix returns the prefix used to name the objects specific to this
// environment.
func (env *azureEnviron) getEnvPrefix() string {
	return fmt.Sprintf("juju-%s-", env.Name())
}

// convertToInstances converts a slice of gwacl.HostedServiceDescriptor objects
// into a slice of instance.Instance objects.
func convertToInstances(services []gwacl.HostedServiceDescriptor, env *azureEnviron) []instance.Instance {
	instances := make([]instance.Instance, len(services))
	for i, service := range services {
		instances[i] = &azureInstance{service, env}
	}
	return instances
}

// Storage is specified in the Environ interface.
func (env *azureEnviron) Storage() environs.Storage {
	return env.getSnapshot().storage
}

// PublicStorage is specified in the Environ interface.
func (env *azureEnviron) PublicStorage() environs.StorageReader {
	return env.getSnapshot().publicStorage
}

// Destroy is specified in the Environ interface.
func (env *azureEnviron) Destroy(ensureInsts []instance.Instance) error {
	logger.Debugf("destroying environment %q", env.name)

	// Delete storage.
	err := env.Storage().RemoveAll()
	if err != nil {
		return fmt.Errorf("cannot clean up storage: %v", err)
	}

	// Stop all instances.
	insts, err := env.AllInstances()
	if err != nil {
		return fmt.Errorf("cannot get instances: %v", err)
	}
	found := make(map[instance.Id]bool)
	for _, inst := range insts {
		found[inst.Id()] = true
	}

	// Add any instances we've been told about but haven't yet shown
	// up in the instance list.
	for _, inst := range ensureInsts {
		id := inst.Id()
		if !found[id] {
			insts = append(insts, inst)
			found[id] = true
		}
	}
	err = env.StopInstances(insts)
	if err != nil {
		return fmt.Errorf("cannot stop instances: %v", err)
	}

	// Delete vnet and affinity group.
	err = env.deleteVirtualNetwork()
	if err != nil {
		return fmt.Errorf("cannot delete the environment's virtual network: %v", err)
	}
	err = env.deleteAffinityGroup()
	if err != nil {
		return fmt.Errorf("cannot delete the environment's affinity group: %v", err)
	}
	return nil
}

// OpenPorts is specified in the Environ interface. However, Azure does not
// support the global firewall mode.
func (env *azureEnviron) OpenPorts(ports []instance.Port) error {
	return nil
}

// ClosePorts is specified in the Environ interface. However, Azure does not
// support the global firewall mode.
func (env *azureEnviron) ClosePorts(ports []instance.Port) error {
	return nil
}

// Ports is specified in the Environ interface.
func (env *azureEnviron) Ports() ([]instance.Port, error) {
	// TODO: implement this.
	return []instance.Port{}, nil
}

// Provider is specified in the Environ interface.
func (env *azureEnviron) Provider() environs.EnvironProvider {
	return azureEnvironProvider{}
}

// azureManagementContext wraps two things: a gwacl.ManagementAPI (effectively
// a session on the Azure management API) and a tempCertFile, which keeps track
// of the temporary certificate file that needs to be deleted once we're done
// with this particular session.
// Since it embeds *gwacl.ManagementAPI, you can use it much as if it were a
// pointer to a ManagementAPI object.  Just don't forget to release it after
// use.
type azureManagementContext struct {
	*gwacl.ManagementAPI
	certFile *tempCertFile
}

var (
	retryPolicy = gwacl.RetryPolicy{
		NbRetries: 6,
		HttpStatusCodes: []int{
			http.StatusConflict,
			http.StatusRequestTimeout,
			http.StatusInternalServerError,
			http.StatusServiceUnavailable,
		},
		Delay: 10 * time.Second}
)

// getManagementAPI obtains a context object for interfacing with Azure's
// management API.
// For now, each invocation just returns a separate object.  This is probably
// wasteful (each context gets its own SSL connection) and may need optimizing
// later.
func (env *azureEnviron) getManagementAPI() (*azureManagementContext, error) {
	snap := env.getSnapshot()
	subscription := snap.ecfg.managementSubscriptionId()
	certData := snap.ecfg.managementCertificate()
	certFile, err := newTempCertFile([]byte(certData))
	if err != nil {
		return nil, err
	}
	// After this point, if we need to leave prematurely, we should clean
	// up that certificate file.
	location := snap.ecfg.location()
	mgtAPI, err := gwacl.NewManagementAPIWithRetryPolicy(subscription, certFile.Path(), location, retryPolicy)
	if err != nil {
		certFile.Delete()
		return nil, err
	}
	context := azureManagementContext{
		ManagementAPI: mgtAPI,
		certFile:      certFile,
	}
	return &context, nil
}

// releaseManagementAPI frees up a context object obtained through
// getManagementAPI.
func (env *azureEnviron) releaseManagementAPI(context *azureManagementContext) {
	// Be tolerant to incomplete context objects, in case we ever get
	// called during cleanup of a failed attempt to create one.
	if context == nil || context.certFile == nil {
		return
	}
	// For now, all that needs doing is to delete the temporary certificate
	// file.  We may do cleverer things later, such as connection pooling
	// where this method returns a context to the pool.
	context.certFile.Delete()
}

// updateStorageAccountKey queries the storage account key, and updates the
// version cached in env.storageAccountKey.
//
// It takes a snapshot in order to preserve transactional integrity relative
// to the snapshot's starting state, without having to lock the environment
// for the duration.  If there is a conflicting change to env relative to the
// state recorded in the snapshot, this function will fail.
func (env *azureEnviron) updateStorageAccountKey(snapshot *azureEnviron) (string, error) {
	// This method follows an RCU pattern, an optimistic technique to
	// implement atomic read-update transactions: get a consistent snapshot
	// of state; process data; enter critical section; check for conflicts;
	// write back changes.  The advantage is that there are no long-held
	// locks, in particular while waiting for the request to Azure to
	// complete.
	// "Get a consistent snapshot of state" is the caller's responsibility.
	// The caller can use env.getSnapshot().

	// Process data: get a current account key from Azure.
	key, err := env.queryStorageAccountKey()
	if err != nil {
		return "", err
	}

	// Enter critical section.
	env.Lock()
	defer env.Unlock()

	// Check for conflicts: is the config still what it was?
	if env.ecfg != snapshot.ecfg {
		// The environment has been reconfigured while we were
		// working on this, so the key we just get may not be
		// appropriate any longer.  So fail.
		// Whatever we were doing isn't likely to be right any more
		// anyway.  Otherwise, it might be worth returning the key
		// just in case it still works, and proceed without updating
		// env.storageAccountKey.
		return "", fmt.Errorf("environment was reconfigured")
	}

	// Write back changes.
	env.storageAccountKey = key
	return key, nil
}

// getStorageContext obtains a context object for interfacing with Azure's
// storage API.
// For now, each invocation just returns a separate object.  This is probably
// wasteful (each context gets its own SSL connection) and may need optimizing
// later.
func (env *azureEnviron) getStorageContext() (*gwacl.StorageContext, error) {
	snap := env.getSnapshot()
	key := snap.storageAccountKey
	if key == "" {
		// We don't know the storage-account key yet.  Request it.
		var err error
		key, err = env.updateStorageAccountKey(snap)
		if err != nil {
			return nil, err
		}
	}
	context := gwacl.StorageContext{
		Account:       snap.ecfg.storageAccountName(),
		Key:           key,
		AzureEndpoint: gwacl.GetEndpoint(snap.ecfg.location()),
		RetryPolicy:   retryPolicy,
	}
	return &context, nil
}

// getPublicStorageContext obtains a context object for interfacing with
// Azure's storage API (public storage).
func (env *azureEnviron) getPublicStorageContext() (*gwacl.StorageContext, error) {
	ecfg := env.getSnapshot().ecfg
	context := gwacl.StorageContext{
		Account:       ecfg.publicStorageAccountName(),
		Key:           "", // Empty string means anonymous access.
		AzureEndpoint: gwacl.GetEndpoint(ecfg.location()),
		RetryPolicy:   retryPolicy,
	}
	// There is currently no way for this to fail.
	return &context, nil
}

// baseURLs specifies an Azure specific location where we look for simplestreams information.
// It contains the central databases for the released and daily streams, but this may
// become more configurable.  This variable is here as a placeholder, but also
// as an injection point for tests.
var baseURLs = []string{
	"http://cloud-images.ubuntu.com/daily",
}

// GetImageBaseURLs returns a list of URLs which are used to search for simplestreams image metadata.
func (e *azureEnviron) GetImageBaseURLs() ([]string, error) {
	return baseURLs, nil
}

// getImageStream returns the name of the simplestreams stream from which
// this environment wants its images, e.g. "releases" or "daily", or the
// blank string for the default.
func (env *azureEnviron) getImageStream() string {
	// Hard-coded to the default for now.
	return ""
}

// getImageMetadataSigningRequired returns whether this environment requires
// image metadata from Simplestreams to be signed.
func (env *azureEnviron) getImageMetadataSigningRequired() bool {
	// Hard-coded to true for now.  Once we support custom base URLs,
	// this may have to change.
	return true
}
