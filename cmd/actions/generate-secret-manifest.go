package actions

import (
	"fmt"
	"kube-svc-ctl/cmd"
	"kube-svc-ctl/docker"
	"kube-svc-ctl/k8s/secret"
	"kube-svc-ctl/vault"
	"os"
)

func GenerateSecretManifest() error {
	c := cmd.NewCommand("generate-secret-manifest")
	c.AddFlag(cmd.VaultUrlFlag())
	c.AddFlag(cmd.VaultTokenFlag())

	c.AddAction(func(c cmd.Command) {
		var tokenPtr *string
		var urlPtr *string
		tokenValue := *c.GetFlag("vault-token").GetValuePtr()
		urlValue := *c.GetFlag("vault-url").GetValuePtr()
		if tokenValue != nil {
			tokenPtr = tokenValue.(*string)
		}
		if urlValue != nil {
			urlPtr = urlValue.(*string)
		}
		println(urlPtr)
		vaultCredentials, err := cmd.ValidateVaultFlags(urlPtr, tokenPtr)
		if err != nil {
			panic(err)
		}

		vaultClient := vault.NewVaultAuthenticatedClient(vaultCredentials)
		d, err := vault.ReadVaultSecret(vaultClient, "secret/data/common/registry")
		var r map[string]interface{}
		if subdata, ok := d["data"]; ok {
			r = subdata.(map[string]interface{})
		} else {
			r = d
		}
		if err != nil {
			panic(err)
		}
		var registryUsername string
		var registryPassword string
		var registryAddr string

		if p, ok := r["registry_password"].(string); ok {
			registryPassword = p
		} else {
			panic("Registry password retrieved from Vault is not valid")
		}
		if p, ok := r["registry_username"].(string); ok {
			registryUsername = p
		} else {
			registryUsername = "k8s-registry-token"
		}
		if p, ok := r["registry_addr"].(string); ok {
			registryAddr = p
		} else {
			registryAddr = "registry.gitlab.com"
		}
		creds := docker.RegistryCredentials{
			Addr:     registryAddr,
			Username: registryUsername,
			Password: registryPassword,
		}
		m := secret.NewSecret(
			"gitlab-registry",
			"kubernetes.io/dockerconfigjson",
			creds.EncryptedDockerconfigjson())
		fmt.Println(m.YamlManifest())
	})
	fs, err := c.Parse(os.Args[2:])
	if err != nil {
		return err
	}
	if fs.Parsed() {
		c.Execute()
	}
	return nil
}
