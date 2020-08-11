package main

import (
	"kube-svc-ctl/docker"
	"kube-svc-ctl/k8s/secret"
	"kube-svc-ctl/svc"
	"kube-svc-ctl/vault"
)
import "flag"
import "os"
import "fmt"

func validateKubernetesResource(name *string, namespace *string, resourceType *string) bool {
	return false
}
func generateSvcConfig(name *string, tag *string, credentials *vault.VaultCredentials) string {
	c := vault.NewVaultAuthenticatedClient(credentials)
	r, err := vault.ReadVaultSecret(c, "secret/data/service/"+*name)
	if err != nil {
		panic(err)
	}

	s := svc.NewHelmConfig(tag, &r)
	return s.ToYaml()
}

func generateRegistrySecretManifest(credentials *vault.VaultCredentials) {
	c := vault.NewVaultAuthenticatedClient(credentials)
	d, err := vault.ReadVaultSecret(c, "secret/data/common/registry")
	var r map[string]interface{}
	if subdata, ok := d["data"]; ok {
		r = subdata.(map[string]interface{})
	} else {
		r = d
	}
	if err != nil {
		panic(err)
	}
	var registryPassword string
	if p, ok := r["registry_password"].(string); ok {
		registryPassword = p
	} else {
		panic("Registry password retrieved from Vault is not valid")
	}
	creds := docker.RegistryCredentials{
		Addr:     "registry.gitlab.com",
		Username: "k8s-registry-token",
		Password: registryPassword,
	}
	m := secret.NewSecret(
		"gitlab-registry",
		"kubernetes.io/dockerconfigjson",
		creds.EncryptedDockerconfigjson())
	fmt.Println(m.YamlManifest())
}

func addVaultFlags(set *flag.FlagSet) {
	set.String(
		"vault-url",
		"http://localhost:8200",
		"Hashicorp Vault Url. Can be overridden by VAULT_URL environment variable. Default: http://localhost:8200",
	)
	set.String(
		"vault-token",
		"",
		"Hashicorp Vault Access Token. Can be overridden by VAULT_TOKEN environment variable (RECOMMENDED)",
	)
}

func validateVaultFlags(url *string, token *string) *vault.VaultCredentials {
	url_env, vaultUrlOk := os.LookupEnv("VAULT_URL")
	token_env, vaultTokenOk := os.LookupEnv("VAULT_TOKEN")
	if (*url == "" || *token == "") && (!vaultUrlOk || !vaultTokenOk) {
		return nil
	}
	if !vaultUrlOk || !vaultTokenOk {
		creds := vault.VaultCredentials{
			Url:   *url,
			Token: *token,
		}
		return &creds
	}
	creds := vault.VaultCredentials{
		Url:   url_env,
		Token: token_env,
	}
	return &creds
}

func main() {
	validateCommand := flag.NewFlagSet("validate", flag.ExitOnError)
	validateCommandService := validateCommand.String("service", "", "Service name (Required)")
	validateCommandResource := validateCommand.String("resource", "", "Resource type <secret|> (Required)")
	validateCommandResourceName := validateCommand.String("resource-name", "", "Resource name to validate (Required)")

	generateConfigCommand := flag.NewFlagSet("generate-svc-config", flag.ExitOnError)
	generateConfigService := generateConfigCommand.
		String("service", "", "Service name")
	generateConfigCommand.Bool("add-docker-image-tag", true, "Add information about service's docker tag to generated configuration. Default: true")
	generateConfigTag := generateConfigCommand.String("tag", "latest", "Service's docker image tag. Default: latest")
	addVaultFlags(generateConfigCommand)

	generateSecretManifest := flag.NewFlagSet("generate-secret-manifest", flag.ExitOnError)
	addVaultFlags(generateSecretManifest)

	if len(os.Args) < 2 {
		fmt.Println("One of the following commands is required: validate|generate-svc-config|generate-secret-manifest")
	}
	switch os.Args[1] {
	case "validate":
		_ = validateCommand.Parse(os.Args[2:])
	case "generate-svc-config":
		_ = generateConfigCommand.Parse(os.Args[2:])
	case "generate-secret-manifest":
		_ = generateSecretManifest.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
	if validateCommand.Parsed() {
		// Check required flags
		if *validateCommandService == "" || *validateCommandResource == "" {
			validateCommand.PrintDefaults()
			os.Exit(1)
		}
		validateKubernetesResource(validateCommandResourceName, validateCommandService, validateCommandResource)
	}
	if generateConfigCommand.Parsed() {
		if *generateConfigService == "" {
			generateConfigCommand.PrintDefaults()
			os.Exit(1)
		}
		u := generateConfigCommand.Lookup("vault-url").Value.String()
		t := generateConfigCommand.Lookup("vault-token").Value.String()
		addTag := generateConfigCommand.Lookup("add-docker-image-tag").Value.String()
		creds := validateVaultFlags(&u, &t)
		if creds == nil {
			fmt.Println("Vault flags are invalid")
			os.Exit(1)
		}
		if addTag == "false" {
			generateConfigTag = nil
		}
		fmt.Println(generateSvcConfig(generateConfigService, generateConfigTag, creds))
	}
	if generateSecretManifest.Parsed() {
		u := generateSecretManifest.Lookup("vault-url").Value.String()
		t := generateSecretManifest.Lookup("vault-token").Value.String()
		creds := validateVaultFlags(&u, &t)
		if creds == nil {
			fmt.Println("Vault flags are invalid")
			os.Exit(1)
		}
		generateRegistrySecretManifest(creds)
	}
}
