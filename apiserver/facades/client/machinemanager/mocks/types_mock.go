// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/machinemanager (interfaces: Backend,StorageInterface,Pool,Model,Machine,Application,Unit,Charm,CharmhubClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	charm "github.com/juju/charm/v12"
	storagecommon "github.com/juju/juju/apiserver/common/storagecommon"
	machinemanager "github.com/juju/juju/apiserver/facades/client/machinemanager"
	instance "github.com/juju/juju/core/instance"
	model "github.com/juju/juju/core/model"
	network "github.com/juju/juju/core/network"
	objectstore "github.com/juju/juju/core/objectstore"
	status "github.com/juju/juju/core/status"
	config "github.com/juju/juju/environs/config"
	charmhub "github.com/juju/juju/internal/charmhub"
	transport "github.com/juju/juju/internal/charmhub/transport"
	state "github.com/juju/juju/state"
	binarystorage "github.com/juju/juju/state/binarystorage"
	names "github.com/juju/names/v4"
	gomock "go.uber.org/mock/gomock"
)

// MockBackend is a mock of Backend interface.
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend.
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance.
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// AddMachineInsideMachine mocks base method.
func (m *MockBackend) AddMachineInsideMachine(arg0 state.MachineTemplate, arg1 string, arg2 instance.ContainerType) (*state.Machine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMachineInsideMachine", arg0, arg1, arg2)
	ret0, _ := ret[0].(*state.Machine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMachineInsideMachine indicates an expected call of AddMachineInsideMachine.
func (mr *MockBackendMockRecorder) AddMachineInsideMachine(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMachineInsideMachine", reflect.TypeOf((*MockBackend)(nil).AddMachineInsideMachine), arg0, arg1, arg2)
}

// AddMachineInsideNewMachine mocks base method.
func (m *MockBackend) AddMachineInsideNewMachine(arg0, arg1 state.MachineTemplate, arg2 instance.ContainerType) (*state.Machine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMachineInsideNewMachine", arg0, arg1, arg2)
	ret0, _ := ret[0].(*state.Machine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMachineInsideNewMachine indicates an expected call of AddMachineInsideNewMachine.
func (mr *MockBackendMockRecorder) AddMachineInsideNewMachine(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMachineInsideNewMachine", reflect.TypeOf((*MockBackend)(nil).AddMachineInsideNewMachine), arg0, arg1, arg2)
}

// AddOneMachine mocks base method.
func (m *MockBackend) AddOneMachine(arg0 state.MachineTemplate) (*state.Machine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOneMachine", arg0)
	ret0, _ := ret[0].(*state.Machine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddOneMachine indicates an expected call of AddOneMachine.
func (mr *MockBackendMockRecorder) AddOneMachine(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOneMachine", reflect.TypeOf((*MockBackend)(nil).AddOneMachine), arg0)
}

// AllMachines mocks base method.
func (m *MockBackend) AllMachines() ([]machinemanager.Machine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllMachines")
	ret0, _ := ret[0].([]machinemanager.Machine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllMachines indicates an expected call of AllMachines.
func (mr *MockBackendMockRecorder) AllMachines() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllMachines", reflect.TypeOf((*MockBackend)(nil).AllMachines))
}

// AllSpaceInfos mocks base method.
func (m *MockBackend) AllSpaceInfos() (network.SpaceInfos, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllSpaceInfos")
	ret0, _ := ret[0].(network.SpaceInfos)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllSpaceInfos indicates an expected call of AllSpaceInfos.
func (mr *MockBackendMockRecorder) AllSpaceInfos() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllSpaceInfos", reflect.TypeOf((*MockBackend)(nil).AllSpaceInfos))
}

// Application mocks base method.
func (m *MockBackend) Application(arg0 string) (machinemanager.Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Application", arg0)
	ret0, _ := ret[0].(machinemanager.Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Application indicates an expected call of Application.
func (mr *MockBackendMockRecorder) Application(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Application", reflect.TypeOf((*MockBackend)(nil).Application), arg0)
}

// GetBlockForType mocks base method.
func (m *MockBackend) GetBlockForType(arg0 state.BlockType) (state.Block, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockForType", arg0)
	ret0, _ := ret[0].(state.Block)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetBlockForType indicates an expected call of GetBlockForType.
func (mr *MockBackendMockRecorder) GetBlockForType(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockForType", reflect.TypeOf((*MockBackend)(nil).GetBlockForType), arg0)
}

// Machine mocks base method.
func (m *MockBackend) Machine(arg0 string) (machinemanager.Machine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Machine", arg0)
	ret0, _ := ret[0].(machinemanager.Machine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Machine indicates an expected call of Machine.
func (mr *MockBackendMockRecorder) Machine(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Machine", reflect.TypeOf((*MockBackend)(nil).Machine), arg0)
}

// Model mocks base method.
func (m *MockBackend) Model() (machinemanager.Model, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Model")
	ret0, _ := ret[0].(machinemanager.Model)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Model indicates an expected call of Model.
func (mr *MockBackendMockRecorder) Model() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Model", reflect.TypeOf((*MockBackend)(nil).Model))
}

// ToolsStorage mocks base method.
func (m *MockBackend) ToolsStorage(arg0 objectstore.ObjectStore) (binarystorage.StorageCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToolsStorage", arg0)
	ret0, _ := ret[0].(binarystorage.StorageCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToolsStorage indicates an expected call of ToolsStorage.
func (mr *MockBackendMockRecorder) ToolsStorage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToolsStorage", reflect.TypeOf((*MockBackend)(nil).ToolsStorage), arg0)
}

// Unit mocks base method.
func (m *MockBackend) Unit(arg0 string) (machinemanager.Unit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unit", arg0)
	ret0, _ := ret[0].(machinemanager.Unit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unit indicates an expected call of Unit.
func (mr *MockBackendMockRecorder) Unit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unit", reflect.TypeOf((*MockBackend)(nil).Unit), arg0)
}

// MockStorageInterface is a mock of StorageInterface interface.
type MockStorageInterface struct {
	ctrl     *gomock.Controller
	recorder *MockStorageInterfaceMockRecorder
}

// MockStorageInterfaceMockRecorder is the mock recorder for MockStorageInterface.
type MockStorageInterfaceMockRecorder struct {
	mock *MockStorageInterface
}

// NewMockStorageInterface creates a new mock instance.
func NewMockStorageInterface(ctrl *gomock.Controller) *MockStorageInterface {
	mock := &MockStorageInterface{ctrl: ctrl}
	mock.recorder = &MockStorageInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageInterface) EXPECT() *MockStorageInterfaceMockRecorder {
	return m.recorder
}

// FilesystemAccess mocks base method.
func (m *MockStorageInterface) FilesystemAccess() storagecommon.FilesystemAccess {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilesystemAccess")
	ret0, _ := ret[0].(storagecommon.FilesystemAccess)
	return ret0
}

// FilesystemAccess indicates an expected call of FilesystemAccess.
func (mr *MockStorageInterfaceMockRecorder) FilesystemAccess() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilesystemAccess", reflect.TypeOf((*MockStorageInterface)(nil).FilesystemAccess))
}

// StorageInstance mocks base method.
func (m *MockStorageInterface) StorageInstance(arg0 names.StorageTag) (state.StorageInstance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StorageInstance", arg0)
	ret0, _ := ret[0].(state.StorageInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StorageInstance indicates an expected call of StorageInstance.
func (mr *MockStorageInterfaceMockRecorder) StorageInstance(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StorageInstance", reflect.TypeOf((*MockStorageInterface)(nil).StorageInstance), arg0)
}

// UnitStorageAttachments mocks base method.
func (m *MockStorageInterface) UnitStorageAttachments(arg0 names.UnitTag) ([]state.StorageAttachment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnitStorageAttachments", arg0)
	ret0, _ := ret[0].([]state.StorageAttachment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnitStorageAttachments indicates an expected call of UnitStorageAttachments.
func (mr *MockStorageInterfaceMockRecorder) UnitStorageAttachments(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnitStorageAttachments", reflect.TypeOf((*MockStorageInterface)(nil).UnitStorageAttachments), arg0)
}

// VolumeAccess mocks base method.
func (m *MockStorageInterface) VolumeAccess() storagecommon.VolumeAccess {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VolumeAccess")
	ret0, _ := ret[0].(storagecommon.VolumeAccess)
	return ret0
}

// VolumeAccess indicates an expected call of VolumeAccess.
func (mr *MockStorageInterfaceMockRecorder) VolumeAccess() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeAccess", reflect.TypeOf((*MockStorageInterface)(nil).VolumeAccess))
}

// MockPool is a mock of Pool interface.
type MockPool struct {
	ctrl     *gomock.Controller
	recorder *MockPoolMockRecorder
}

// MockPoolMockRecorder is the mock recorder for MockPool.
type MockPoolMockRecorder struct {
	mock *MockPool
}

// NewMockPool creates a new mock instance.
func NewMockPool(ctrl *gomock.Controller) *MockPool {
	mock := &MockPool{ctrl: ctrl}
	mock.recorder = &MockPoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPool) EXPECT() *MockPoolMockRecorder {
	return m.recorder
}

// GetModel mocks base method.
func (m *MockPool) GetModel(arg0 string) (machinemanager.Model, func(), error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetModel", arg0)
	ret0, _ := ret[0].(machinemanager.Model)
	ret1, _ := ret[1].(func())
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetModel indicates an expected call of GetModel.
func (mr *MockPoolMockRecorder) GetModel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetModel", reflect.TypeOf((*MockPool)(nil).GetModel), arg0)
}

// SystemState mocks base method.
func (m *MockPool) SystemState() (machinemanager.ControllerBackend, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SystemState")
	ret0, _ := ret[0].(machinemanager.ControllerBackend)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SystemState indicates an expected call of SystemState.
func (mr *MockPoolMockRecorder) SystemState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SystemState", reflect.TypeOf((*MockPool)(nil).SystemState))
}

// MockModel is a mock of Model interface.
type MockModel struct {
	ctrl     *gomock.Controller
	recorder *MockModelMockRecorder
}

// MockModelMockRecorder is the mock recorder for MockModel.
type MockModelMockRecorder struct {
	mock *MockModel
}

// NewMockModel creates a new mock instance.
func NewMockModel(ctrl *gomock.Controller) *MockModel {
	mock := &MockModel{ctrl: ctrl}
	mock.recorder = &MockModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModel) EXPECT() *MockModelMockRecorder {
	return m.recorder
}

// CloudCredentialTag mocks base method.
func (m *MockModel) CloudCredentialTag() (names.CloudCredentialTag, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudCredentialTag")
	ret0, _ := ret[0].(names.CloudCredentialTag)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// CloudCredentialTag indicates an expected call of CloudCredentialTag.
func (mr *MockModelMockRecorder) CloudCredentialTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudCredentialTag", reflect.TypeOf((*MockModel)(nil).CloudCredentialTag))
}

// CloudName mocks base method.
func (m *MockModel) CloudName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudName")
	ret0, _ := ret[0].(string)
	return ret0
}

// CloudName indicates an expected call of CloudName.
func (mr *MockModelMockRecorder) CloudName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudName", reflect.TypeOf((*MockModel)(nil).CloudName))
}

// CloudRegion mocks base method.
func (m *MockModel) CloudRegion() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloudRegion")
	ret0, _ := ret[0].(string)
	return ret0
}

// CloudRegion indicates an expected call of CloudRegion.
func (mr *MockModelMockRecorder) CloudRegion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloudRegion", reflect.TypeOf((*MockModel)(nil).CloudRegion))
}

// Config mocks base method.
func (m *MockModel) Config() (*config.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Config")
	ret0, _ := ret[0].(*config.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Config indicates an expected call of Config.
func (mr *MockModelMockRecorder) Config() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockModel)(nil).Config))
}

// ControllerUUID mocks base method.
func (m *MockModel) ControllerUUID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ControllerUUID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ControllerUUID indicates an expected call of ControllerUUID.
func (mr *MockModelMockRecorder) ControllerUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ControllerUUID", reflect.TypeOf((*MockModel)(nil).ControllerUUID))
}

// ModelTag mocks base method.
func (m *MockModel) ModelTag() names.ModelTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModelTag")
	ret0, _ := ret[0].(names.ModelTag)
	return ret0
}

// ModelTag indicates an expected call of ModelTag.
func (mr *MockModelMockRecorder) ModelTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelTag", reflect.TypeOf((*MockModel)(nil).ModelTag))
}

// Name mocks base method.
func (m *MockModel) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockModelMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockModel)(nil).Name))
}

// Type mocks base method.
func (m *MockModel) Type() state.ModelType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(state.ModelType)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockModelMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockModel)(nil).Type))
}

// UUID mocks base method.
func (m *MockModel) UUID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UUID")
	ret0, _ := ret[0].(string)
	return ret0
}

// UUID indicates an expected call of UUID.
func (mr *MockModelMockRecorder) UUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UUID", reflect.TypeOf((*MockModel)(nil).UUID))
}

// MockMachine is a mock of Machine interface.
type MockMachine struct {
	ctrl     *gomock.Controller
	recorder *MockMachineMockRecorder
}

// MockMachineMockRecorder is the mock recorder for MockMachine.
type MockMachineMockRecorder struct {
	mock *MockMachine
}

// NewMockMachine creates a new mock instance.
func NewMockMachine(ctrl *gomock.Controller) *MockMachine {
	mock := &MockMachine{ctrl: ctrl}
	mock.recorder = &MockMachineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMachine) EXPECT() *MockMachineMockRecorder {
	return m.recorder
}

// ApplicationNames mocks base method.
func (m *MockMachine) ApplicationNames() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApplicationNames")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ApplicationNames indicates an expected call of ApplicationNames.
func (mr *MockMachineMockRecorder) ApplicationNames() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApplicationNames", reflect.TypeOf((*MockMachine)(nil).ApplicationNames))
}

// Base mocks base method.
func (m *MockMachine) Base() state.Base {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Base")
	ret0, _ := ret[0].(state.Base)
	return ret0
}

// Base indicates an expected call of Base.
func (mr *MockMachineMockRecorder) Base() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Base", reflect.TypeOf((*MockMachine)(nil).Base))
}

// CompleteUpgradeSeries mocks base method.
func (m *MockMachine) CompleteUpgradeSeries() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompleteUpgradeSeries")
	ret0, _ := ret[0].(error)
	return ret0
}

// CompleteUpgradeSeries indicates an expected call of CompleteUpgradeSeries.
func (mr *MockMachineMockRecorder) CompleteUpgradeSeries() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompleteUpgradeSeries", reflect.TypeOf((*MockMachine)(nil).CompleteUpgradeSeries))
}

// Containers mocks base method.
func (m *MockMachine) Containers() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Containers")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Containers indicates an expected call of Containers.
func (mr *MockMachineMockRecorder) Containers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Containers", reflect.TypeOf((*MockMachine)(nil).Containers))
}

// CreateUpgradeSeriesLock mocks base method.
func (m *MockMachine) CreateUpgradeSeriesLock(arg0 []string, arg1 state.Base) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUpgradeSeriesLock", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUpgradeSeriesLock indicates an expected call of CreateUpgradeSeriesLock.
func (mr *MockMachineMockRecorder) CreateUpgradeSeriesLock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUpgradeSeriesLock", reflect.TypeOf((*MockMachine)(nil).CreateUpgradeSeriesLock), arg0, arg1)
}

// Destroy mocks base method.
func (m *MockMachine) Destroy(arg0 objectstore.ObjectStore) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy.
func (mr *MockMachineMockRecorder) Destroy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockMachine)(nil).Destroy), arg0)
}

// ForceDestroy mocks base method.
func (m *MockMachine) ForceDestroy(arg0 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForceDestroy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForceDestroy indicates an expected call of ForceDestroy.
func (mr *MockMachineMockRecorder) ForceDestroy(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForceDestroy", reflect.TypeOf((*MockMachine)(nil).ForceDestroy), arg0)
}

// GetUpgradeSeriesMessages mocks base method.
func (m *MockMachine) GetUpgradeSeriesMessages() ([]string, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpgradeSeriesMessages")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUpgradeSeriesMessages indicates an expected call of GetUpgradeSeriesMessages.
func (mr *MockMachineMockRecorder) GetUpgradeSeriesMessages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpgradeSeriesMessages", reflect.TypeOf((*MockMachine)(nil).GetUpgradeSeriesMessages))
}

// HardwareCharacteristics mocks base method.
func (m *MockMachine) HardwareCharacteristics() (*instance.HardwareCharacteristics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HardwareCharacteristics")
	ret0, _ := ret[0].(*instance.HardwareCharacteristics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HardwareCharacteristics indicates an expected call of HardwareCharacteristics.
func (mr *MockMachineMockRecorder) HardwareCharacteristics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HardwareCharacteristics", reflect.TypeOf((*MockMachine)(nil).HardwareCharacteristics))
}

// Id mocks base method.
func (m *MockMachine) Id() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(string)
	return ret0
}

// Id indicates an expected call of Id.
func (mr *MockMachineMockRecorder) Id() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*MockMachine)(nil).Id))
}

// InstanceStatus mocks base method.
func (m *MockMachine) InstanceStatus() (status.StatusInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceStatus")
	ret0, _ := ret[0].(status.StatusInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstanceStatus indicates an expected call of InstanceStatus.
func (mr *MockMachineMockRecorder) InstanceStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceStatus", reflect.TypeOf((*MockMachine)(nil).InstanceStatus))
}

// IsLockedForSeriesUpgrade mocks base method.
func (m *MockMachine) IsLockedForSeriesUpgrade() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLockedForSeriesUpgrade")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLockedForSeriesUpgrade indicates an expected call of IsLockedForSeriesUpgrade.
func (mr *MockMachineMockRecorder) IsLockedForSeriesUpgrade() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLockedForSeriesUpgrade", reflect.TypeOf((*MockMachine)(nil).IsLockedForSeriesUpgrade))
}

// IsManager mocks base method.
func (m *MockMachine) IsManager() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsManager")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsManager indicates an expected call of IsManager.
func (mr *MockMachineMockRecorder) IsManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsManager", reflect.TypeOf((*MockMachine)(nil).IsManager))
}

// Principals mocks base method.
func (m *MockMachine) Principals() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Principals")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Principals indicates an expected call of Principals.
func (mr *MockMachineMockRecorder) Principals() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Principals", reflect.TypeOf((*MockMachine)(nil).Principals))
}

// RemoveUpgradeSeriesLock mocks base method.
func (m *MockMachine) RemoveUpgradeSeriesLock() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUpgradeSeriesLock")
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUpgradeSeriesLock indicates an expected call of RemoveUpgradeSeriesLock.
func (mr *MockMachineMockRecorder) RemoveUpgradeSeriesLock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUpgradeSeriesLock", reflect.TypeOf((*MockMachine)(nil).RemoveUpgradeSeriesLock))
}

// SetInstanceStatus mocks base method.
func (m *MockMachine) SetInstanceStatus(arg0 status.StatusInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetInstanceStatus", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetInstanceStatus indicates an expected call of SetInstanceStatus.
func (mr *MockMachineMockRecorder) SetInstanceStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInstanceStatus", reflect.TypeOf((*MockMachine)(nil).SetInstanceStatus), arg0)
}

// SetKeepInstance mocks base method.
func (m *MockMachine) SetKeepInstance(arg0 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetKeepInstance", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetKeepInstance indicates an expected call of SetKeepInstance.
func (mr *MockMachineMockRecorder) SetKeepInstance(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetKeepInstance", reflect.TypeOf((*MockMachine)(nil).SetKeepInstance), arg0)
}

// SetPassword mocks base method.
func (m *MockMachine) SetPassword(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPassword", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPassword indicates an expected call of SetPassword.
func (mr *MockMachineMockRecorder) SetPassword(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPassword", reflect.TypeOf((*MockMachine)(nil).SetPassword), arg0)
}

// SetUpgradeSeriesStatus mocks base method.
func (m *MockMachine) SetUpgradeSeriesStatus(arg0 model.UpgradeSeriesStatus, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetUpgradeSeriesStatus", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUpgradeSeriesStatus indicates an expected call of SetUpgradeSeriesStatus.
func (mr *MockMachineMockRecorder) SetUpgradeSeriesStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUpgradeSeriesStatus", reflect.TypeOf((*MockMachine)(nil).SetUpgradeSeriesStatus), arg0, arg1)
}

// Tag mocks base method.
func (m *MockMachine) Tag() names.Tag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tag")
	ret0, _ := ret[0].(names.Tag)
	return ret0
}

// Tag indicates an expected call of Tag.
func (mr *MockMachineMockRecorder) Tag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockMachine)(nil).Tag))
}

// Units mocks base method.
func (m *MockMachine) Units() ([]machinemanager.Unit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Units")
	ret0, _ := ret[0].([]machinemanager.Unit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Units indicates an expected call of Units.
func (mr *MockMachineMockRecorder) Units() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Units", reflect.TypeOf((*MockMachine)(nil).Units))
}

// UpgradeSeriesStatus mocks base method.
func (m *MockMachine) UpgradeSeriesStatus() (model.UpgradeSeriesStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpgradeSeriesStatus")
	ret0, _ := ret[0].(model.UpgradeSeriesStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpgradeSeriesStatus indicates an expected call of UpgradeSeriesStatus.
func (mr *MockMachineMockRecorder) UpgradeSeriesStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpgradeSeriesStatus", reflect.TypeOf((*MockMachine)(nil).UpgradeSeriesStatus))
}

// WatchUpgradeSeriesNotifications mocks base method.
func (m *MockMachine) WatchUpgradeSeriesNotifications() (state.NotifyWatcher, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchUpgradeSeriesNotifications")
	ret0, _ := ret[0].(state.NotifyWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchUpgradeSeriesNotifications indicates an expected call of WatchUpgradeSeriesNotifications.
func (mr *MockMachineMockRecorder) WatchUpgradeSeriesNotifications() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchUpgradeSeriesNotifications", reflect.TypeOf((*MockMachine)(nil).WatchUpgradeSeriesNotifications))
}

// MockApplication is a mock of Application interface.
type MockApplication struct {
	ctrl     *gomock.Controller
	recorder *MockApplicationMockRecorder
}

// MockApplicationMockRecorder is the mock recorder for MockApplication.
type MockApplicationMockRecorder struct {
	mock *MockApplication
}

// NewMockApplication creates a new mock instance.
func NewMockApplication(ctrl *gomock.Controller) *MockApplication {
	mock := &MockApplication{ctrl: ctrl}
	mock.recorder = &MockApplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder {
	return m.recorder
}

// Charm mocks base method.
func (m *MockApplication) Charm() (machinemanager.Charm, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Charm")
	ret0, _ := ret[0].(machinemanager.Charm)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Charm indicates an expected call of Charm.
func (mr *MockApplicationMockRecorder) Charm() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Charm", reflect.TypeOf((*MockApplication)(nil).Charm))
}

// CharmOrigin mocks base method.
func (m *MockApplication) CharmOrigin() *state.CharmOrigin {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CharmOrigin")
	ret0, _ := ret[0].(*state.CharmOrigin)
	return ret0
}

// CharmOrigin indicates an expected call of CharmOrigin.
func (mr *MockApplicationMockRecorder) CharmOrigin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CharmOrigin", reflect.TypeOf((*MockApplication)(nil).CharmOrigin))
}

// Name mocks base method.
func (m *MockApplication) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockApplicationMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockApplication)(nil).Name))
}

// MockUnit is a mock of Unit interface.
type MockUnit struct {
	ctrl     *gomock.Controller
	recorder *MockUnitMockRecorder
}

// MockUnitMockRecorder is the mock recorder for MockUnit.
type MockUnitMockRecorder struct {
	mock *MockUnit
}

// NewMockUnit creates a new mock instance.
func NewMockUnit(ctrl *gomock.Controller) *MockUnit {
	mock := &MockUnit{ctrl: ctrl}
	mock.recorder = &MockUnitMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnit) EXPECT() *MockUnitMockRecorder {
	return m.recorder
}

// AgentStatus mocks base method.
func (m *MockUnit) AgentStatus() (status.StatusInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AgentStatus")
	ret0, _ := ret[0].(status.StatusInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AgentStatus indicates an expected call of AgentStatus.
func (mr *MockUnitMockRecorder) AgentStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AgentStatus", reflect.TypeOf((*MockUnit)(nil).AgentStatus))
}

// Name mocks base method.
func (m *MockUnit) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockUnitMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockUnit)(nil).Name))
}

// Status mocks base method.
func (m *MockUnit) Status() (status.StatusInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(status.StatusInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Status indicates an expected call of Status.
func (mr *MockUnitMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockUnit)(nil).Status))
}

// UnitTag mocks base method.
func (m *MockUnit) UnitTag() names.UnitTag {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnitTag")
	ret0, _ := ret[0].(names.UnitTag)
	return ret0
}

// UnitTag indicates an expected call of UnitTag.
func (mr *MockUnitMockRecorder) UnitTag() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnitTag", reflect.TypeOf((*MockUnit)(nil).UnitTag))
}

// MockCharm is a mock of Charm interface.
type MockCharm struct {
	ctrl     *gomock.Controller
	recorder *MockCharmMockRecorder
}

// MockCharmMockRecorder is the mock recorder for MockCharm.
type MockCharmMockRecorder struct {
	mock *MockCharm
}

// NewMockCharm creates a new mock instance.
func NewMockCharm(ctrl *gomock.Controller) *MockCharm {
	mock := &MockCharm{ctrl: ctrl}
	mock.recorder = &MockCharmMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCharm) EXPECT() *MockCharmMockRecorder {
	return m.recorder
}

// Manifest mocks base method.
func (m *MockCharm) Manifest() *charm.Manifest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Manifest")
	ret0, _ := ret[0].(*charm.Manifest)
	return ret0
}

// Manifest indicates an expected call of Manifest.
func (mr *MockCharmMockRecorder) Manifest() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Manifest", reflect.TypeOf((*MockCharm)(nil).Manifest))
}

// Meta mocks base method.
func (m *MockCharm) Meta() *charm.Meta {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Meta")
	ret0, _ := ret[0].(*charm.Meta)
	return ret0
}

// Meta indicates an expected call of Meta.
func (mr *MockCharmMockRecorder) Meta() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Meta", reflect.TypeOf((*MockCharm)(nil).Meta))
}

// MockCharmhubClient is a mock of CharmhubClient interface.
type MockCharmhubClient struct {
	ctrl     *gomock.Controller
	recorder *MockCharmhubClientMockRecorder
}

// MockCharmhubClientMockRecorder is the mock recorder for MockCharmhubClient.
type MockCharmhubClientMockRecorder struct {
	mock *MockCharmhubClient
}

// NewMockCharmhubClient creates a new mock instance.
func NewMockCharmhubClient(ctrl *gomock.Controller) *MockCharmhubClient {
	mock := &MockCharmhubClient{ctrl: ctrl}
	mock.recorder = &MockCharmhubClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCharmhubClient) EXPECT() *MockCharmhubClientMockRecorder {
	return m.recorder
}

// Refresh mocks base method.
func (m *MockCharmhubClient) Refresh(arg0 context.Context, arg1 charmhub.RefreshConfig) ([]transport.RefreshResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", arg0, arg1)
	ret0, _ := ret[0].([]transport.RefreshResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh.
func (mr *MockCharmhubClientMockRecorder) Refresh(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockCharmhubClient)(nil).Refresh), arg0, arg1)
}
