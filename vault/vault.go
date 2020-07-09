package vault

import (
	"crypto/tls"
	"github.com/hashicorp/vault/api"
	"net/http"
	"time"
)

type VaultCredentials struct {
	Url   string
	Token string
}

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var httpClient = &http.Client{
	Timeout:   10 * time.Second,
	Transport: tr,
}

func NewVaultAuthenticatedClient(credentials *VaultCredentials) *api.Client {
	client, err := api.NewClient(&api.Config{Address: credentials.Url, HttpClient: httpClient})
	if err != nil {
		panic(err)
	}
	client.SetToken(credentials.Token)
	return client
}

func ReadVaultSecret(client *api.Client, path string) (map[string]interface{}, error) {
	data, err := client.Logical().Read(path)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return data.Data["data"].(map[string]interface{}), nil
}