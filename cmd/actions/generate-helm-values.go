package actions

import (
	"flag"
	"fmt"
	"kube-svc-ctl/cmd"
	"kube-svc-ctl/helm"
	"kube-svc-ctl/vault"
	"net"
	"os"
)

func GenerateHelmValues() error {
	addDockerImage := cmd.NewFlag("add-docker-image-tag", func(set *flag.FlagSet, cf *cmd.CommandFlag) *interface{} {
		var data interface{}
		data = set.Bool(cf.Name, true, "Add information about service's docker tag to generated configuration. Default: true")
		return &data
	})

	dockerImageTag := cmd.NewFlag("tag", func(set *flag.FlagSet, cf *cmd.CommandFlag) *interface{} {
		var data interface{}
		data = set.String(cf.Name, "latest", "Service's docker image tag. Default: latest")
		return &data
	})

	targetYamlTree := cmd.NewFlag("target-yaml-tree", func(set *flag.FlagSet, cf *cmd.CommandFlag) *interface{} {
		var data interface{}
		data = set.String(cf.Name, "secrets.datas", "Target YAML tree path where to put secrets. Default: .secrets")
		return &data
	})

	targetYamlFormat := cmd.NewFlag("target-yaml-format", func(set *flag.FlagSet, cf *cmd.CommandFlag) *interface{} {
		var data interface{}
		data = set.String(cf.Name, "map", "Target YAML formatting. Acceptable values: map, list. Default: map")
		return &data
	})

	c := cmd.NewCommand("generate-helm-values")
	action := func(c cmd.Command) {
		var tokenPtr *string
		var urlPtr *string
		var serviceNamePtr *string
		var tagPtr *string
		tokenValue := *c.GetFlag("vault-token").GetValuePtr()
		urlValue := *c.GetFlag("vault-url").GetValuePtr()
		serviceNameValue := *c.GetFlag("service").GetValuePtr()
		tagValue := *c.GetFlag("tag").GetValuePtr()

		if tokenValue != nil {
			tokenPtr = tokenValue.(*string)
		}
		if urlValue != nil {
			urlPtr = tokenValue.(*string)
		}
		if serviceNameValue != nil {
			serviceNamePtr = serviceNameValue.(*string)
		}
		if tagValue != nil {
			tagPtr = tagValue.(*string)
		}

		targetTreeValue := *c.GetFlag("target-yaml-tree").GetValuePtr()
		var targetTreePtr *string
		if targetTreeValue != nil {
			targetTreePtr = targetTreeValue.(*string)
		}

		targetFormatValue := *c.GetFlag("target-yaml-format").GetValuePtr()
		var targetFormatPtr *string
		if targetFormatValue != nil {
			targetFormatPtr = targetFormatValue.(*string)
		}

		vaultCredentials, err := cmd.ValidateVaultFlags(tokenPtr, urlPtr)
		if err != nil {
			panic(err)
		}
		vaultClient := vault.NewVaultAuthenticatedClient(vaultCredentials)
		var retries = 10
		r, err := vault.ReadVaultSecret(vaultClient, fmt.Sprintf("secret/data/service/%s", *serviceNamePtr))
		for retries != 0 && err != nil {
			if test, ok := err.(net.Error); ok && test.Timeout() {
				r, err = vault.ReadVaultSecret(vaultClient, fmt.Sprintf("secret/data/service/%s", *serviceNamePtr))
				retries = retries - 1
			} else {
				panic(err)
			}
		}
		if retries == 0 && err != nil {
			panic(err)
		}

		s := helm.NewHelmConfig(tagPtr, &r)
		fmt.Println(s.ToYaml(targetTreePtr, targetFormatPtr))
	}
	c.AddAction(cmd.CommandAction(action))
	c.AddFlag(cmd.ServiceNameFlag())
	c.AddFlag(addDockerImage)
	c.AddFlag(dockerImageTag)
	c.AddFlag(targetYamlTree)
	c.AddFlag(targetYamlFormat)
	c.AddFlag(cmd.VaultUrlFlag())
	c.AddFlag(cmd.VaultTokenFlag())

	fs, err := c.Parse(os.Args[2:])
	if err != nil {
		return err
	}
	if fs.Parsed() {
		c.Execute()
	}
	return nil
}
