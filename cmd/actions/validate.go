package actions

import (
	"errors"
	"flag"
	"kube-svc-ctl/cmd"
)

func ValidateCommand() {
	rType := cmd.NewFlag(
		"resource", func(set *flag.FlagSet, cf *cmd.CommandFlag) *interface{} {
			var data interface{}
			data = set.String(cf.Name, "", "Resource type <secret|> (Required)")
			return &data
		})
	rType.AddValidator(func(cf *cmd.CommandFlag) error {
		if (*cf.GetValuePtr()).(string) != "secret" {
			return errors.New("resource type (required) must be one of those: secret|")
		}
		return nil
	})

	rName := cmd.NewFlag(
		"resource-name", func(set *flag.FlagSet, cf *cmd.CommandFlag) *interface{} {
			var data interface{}
			data = set.String(cf.Name, "", "Resource name to validate (Required)")
			return &data
		})
	rName.AddValidator(func(cf *cmd.CommandFlag) error {
		if (*cf.GetValuePtr()).(string) == "" {
			return errors.New("resource name (required) is not specified")
		}
		return nil
	})
	c := cmd.NewCommand("validate")
	c.AddFlag(cmd.ServiceNameFlag())
	c.AddFlag(rName)
	c.AddFlag(rType)
}
