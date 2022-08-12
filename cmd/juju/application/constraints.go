// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"fmt"
	"io"

	"github.com/juju/cmd/v3"
	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"github.com/juju/names/v4"

	"github.com/juju/juju/api/client/application"
	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/juju/block"
	"github.com/juju/juju/cmd/juju/config"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/core/constraints"
)

const (
	constraintsSummary = `Gets or sets machine constraints for an application.`
	constraintsDetails = `
To view all machine constraints for an application, run
    juju constraints <app>
By default, constraints will be printed in a key=value format. You can instead
print it in json or yaml format using the --format flag:
    juju constraints <app> --format yaml

To view the value of a single constraint, run
    juju constraints <app> key
To set constraints, run
    juju constraints <app> key1=val1 key2=val2 ...
This sets "key1" to "val1", etc. Note that setting constraints will clear all
previous constraint values. For example, if you run
    juju constraints <app> a=1 b=2
    juju constraints <app> a=3
the "b=2" constraint will no longer apply. To clear all constraints, run
    juju constraints --clear

Constraints can be imported from a yaml file using the --file flag:
    juju constraints <app> --file=path/to/cfg.yaml
This allows you to e.g. save an app's constraints to a file:
    juju constraints app1 > cfg.yaml
and then import the constraints later. You can also read from stdin using "-",
which allows you to pipe constraints from one app to another:
    juju constraints app1 | juju constraints app2 --file -
You can simultaneously read constraints from a yaml file and set constraints
as above. The command-line args will override any values specified in the file.

This command sets constraints on an application level. To set constraints on a
model level, use the model-constraints command. Application-level constraints
take precedence over model-level constraints.

See also:
    config
    model-constraints
`
)

// NewConstraintsCommand returns a command which gets/sets application
// constraints.
func NewConstraintsCommand() modelcmd.ModelCommand {
	return modelcmd.Wrap(
		&constraintsCommand{
			configBase: config.ConfigCommandBase{
				Resettable: false,
			},
		},
	)
}

type constraintsAPI interface {
	Close() error
	GetConstraints(...string) ([]constraints.Value, error)
	SetConstraints(string, constraints.Value) error
}

type constraintsCommand struct {
	modelcmd.ModelCommandBase
	configBase config.ConfigCommandBase
	api        constraintsAPI
	out        cmd.Output

	applicationName string
	clear           bool
}

func (c *constraintsCommand) getAPI() (constraintsAPI, error) {
	if c.api != nil {
		return c.api, nil
	}
	root, err := c.NewAPIRoot()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return application.NewClient(root), nil
}

func (c *constraintsCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    "constraints",
		Args:    "<application> [<constraint>[=<value>] ... | --file <path>]",
		Purpose: constraintsSummary,
		Doc:     constraintsDetails,
	})
}

func formatConstraints(writer io.Writer, value interface{}) error {
	fmt.Fprint(writer, value.(constraints.Value).String())
	return nil
}

func (c *constraintsCommand) SetFlags(f *gnuflag.FlagSet) {
	// Set the -B / --no-browser-login flag, and model/controller specific flags
	c.ModelCommandBase.SetFlags(f)
	// Set ConfigCommandBase flags
	c.configBase.SetFlags(f)

	// Set the --format and -o flags
	c.out.AddFlags(f, "constraints", map[string]cmd.Formatter{
		"constraints": formatConstraints,
		"yaml":        cmd.FormatYaml,
		"json":        cmd.FormatJson,
	})

	// Set --clear flag to reset all constraints
	f.BoolVar(&c.clear, "clear", false, "Reset all application constraints")
}

func (c *constraintsCommand) Init(args []string) error {
	if len(args) == 0 {
		return errors.Errorf("no application name specified")
	}
	if !names.IsValidApplication(args[0]) {
		return errors.Errorf("invalid application name %q", args[0])
	}

	c.applicationName = args[0]
	return c.configBase.Init(args[1:])
}

func (c *constraintsCommand) Run(ctx *cmd.Context) error {
	client, err := c.getAPI()
	if err != nil {
		return err
	}
	defer client.Close()

	if c.clear {
		return c.clearConstraints(client)
	}

	for _, action := range c.configBase.Actions {
		var err error
		switch action {
		case config.GetOne:
			err = c.getConstraint(client, ctx)
		case config.SetArgs:
			err = c.setConstraints(client, c.configBase.ValsToSet)
		case config.SetFile:
			var attrs config.Attrs
			attrs, err = c.configBase.ReadFile(ctx)
			if err != nil {
				return errors.Trace(err)
			}
			err = c.setConstraints(client, attrs)
		default:
			err = c.getAllConstraints(client, ctx)
		}
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (c *constraintsCommand) clearConstraints(client constraintsAPI) error {
	err := client.SetConstraints(c.applicationName, constraints.Value{})
	return block.ProcessBlockedError(err, block.BlockChange)
}

func (c *constraintsCommand) getAllConstraints(client constraintsAPI, ctx *cmd.Context) error {
	cons, err := client.GetConstraints(c.applicationName)
	if err != nil {
		return err
	}
	return c.out.Write(ctx, cons[0])
}

func (c *constraintsCommand) getConstraint(client constraintsAPI, ctx *cmd.Context) error {
	cons, err := client.GetConstraints(c.applicationName)
	if err != nil {
		return err
	}
	// find key to get within cons[0]
	val, err := cons[0].Lookup(c.configBase.KeysToGet[0])
	if err != nil {
		return errors.Trace(err)
	}
	return c.out.Write(ctx, val)
}

func (c *constraintsCommand) setConstraints(client constraintsAPI, attrs config.Attrs) error {
	// Turn attrs into a constraints.Value
	kvpairs := []string{}
	for k, v := range attrs {
		kvpairs = append(kvpairs, fmt.Sprintf("%s=%v", k, v))
	}
	cons, err := constraints.Parse(kvpairs...)
	if err != nil {
		return errors.Trace(err)
	}

	err = client.SetConstraints(c.applicationName, cons)
	return block.ProcessBlockedError(err, block.BlockChange)
}
