package cmd

import (
	"errors"
	"flag"
	"kube-svc-ctl/vault"
	"os"
)

func ValidateVaultFlags(url *string, token *string) (*vault.VaultCredentials, error) {
	urlEnv, vaultUrlOk := os.LookupEnv("VAULT_URL")
	tokenEnv, vaultTokenOk := os.LookupEnv("VAULT_TOKEN")
	if (*url == "" || *token == "") && (!vaultUrlOk || !vaultTokenOk) {
		return nil, errors.New("No vault token and/or url supplied ")
	}
	if !vaultUrlOk || !vaultTokenOk {
		creds := vault.VaultCredentials{
			Url:   *url,
			Token: *token,
		}
		return &creds, nil
	}
	creds := vault.VaultCredentials{
		Url:   urlEnv,
		Token: tokenEnv,
	}
	return &creds, nil
}

func ServiceNameFlag() *CommandFlag {
	sName := NewFlag(
		"service", func(set *flag.FlagSet, cf *CommandFlag) *interface{} {
			var data interface{}
			data = set.String(cf.Name, "", "Service name (Required)")
			return &data
		})
	sName.AddValidator(func(cf *CommandFlag) error {
		if cf.GetValuePtr() == nil || (*cf.GetValuePtr()) == "" {
			return errors.New("service name (required) is not specified")
		} else {
			return nil
		}
	})
	return sName
}

func VaultUrlFlag() *CommandFlag {
	vaultUrl := NewFlag("vault-url", func(set *flag.FlagSet, cf *CommandFlag) *interface{} {
		var data interface{}
		data = set.String(cf.Name, "http://localhost:8200", "Hashicorp Vault Url. Can be overridden by VAULT_URL environment variable. Default: http://localhost:8200")
		return &data
	})
	return vaultUrl
}

func VaultTokenFlag() *CommandFlag {
	vaultToken := NewFlag("vault-token", func(set *flag.FlagSet, cf *CommandFlag) *interface{} {
		var data interface{}
		data = set.String(cf.Name, "", "Hashicorp Vault Access Token. Can be overridden by VAULT_TOKEN environment variable (RECOMMENDED)")
		return &data
	})
	return vaultToken
}
