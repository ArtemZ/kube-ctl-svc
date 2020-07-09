package docker

import (
	"encoding/base64"
	"fmt"
)

type RegistryCredentials struct {
	Addr, Username, Password string
}

func (cr *RegistryCredentials) Encrypted() string {
	return base64.StdEncoding.
		EncodeToString([]byte(fmt.Sprintf("%s:%s", cr.Username, cr.Password)))
}

func (cr *RegistryCredentials) EncryptedDockerconfigjson() string {
	plain := fmt.
		Sprintf("{ \"auths\": { \"%s\": {\"auth\": \"%s\"} } }", cr.Addr, cr.Encrypted())
	return base64.StdEncoding.EncodeToString([]byte(plain))
}
