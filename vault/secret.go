package vault

type VaultSecret struct {
	Secret map[string]interface{}
}

func mapIndexExists(m *map[string]interface{}, i string) bool {
	_, exists := (*m)[i]
	return exists
}

func NewVaultSecret(plainSecret *map[string]interface{}) (*VaultSecret, error) {
	var s map[string]interface{}
	// check that we are working with a newer Vault version
	// retrieve values from "data" submap in this case
	if mapIndexExists(plainSecret, "data") && mapIndexExists(plainSecret, "metadata") {
		s = (*plainSecret)["data"].(map[string]interface{})
	} else { // use returned values directly otherwise
		s = *plainSecret
	}

	return &VaultSecret{Secret: s}, nil
}
