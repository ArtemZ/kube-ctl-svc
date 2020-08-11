package secret

import "log"
import "gopkg.in/yaml.v2"

type Secret struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Type       string
	Data       SecretData
	Metadata   SecretMetadata
}
type SecretData struct {
	Dockerconfigjson string `yaml:".dockerconfigjson"`
}
type SecretMetadata struct {
	Name string
}

func (s *Secret) YamlManifest() string {
	y, err := yaml.Marshal(s)
	if err != nil {
		log.Fatalf("Yaml Marshalling error: %v", err)
	}
	return string(y)
}

func NewSecret(name string, Type string, dockerconfigjson string) *Secret {
	meta := SecretMetadata{Name: name}
	data := SecretData{Dockerconfigjson: dockerconfigjson}
	s := Secret{
		ApiVersion: "v1",
		Kind:       "Secret",
		Type:       Type,
		Data:       data,
		Metadata:   meta,
	}
	return &s
}
