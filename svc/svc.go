package svc

import (
	"gopkg.in/yaml.v2"
	"log"
)

type HelmConfig struct {
	Image   DockerImageHelmConfig
	Secrets SecretsHelmConfig
}

type DockerImageHelmConfig struct {
	Tag string
}

type SecretsHelmConfig struct {
	Datas map[string]interface{}
}

func mapIndexExists(m *map[string]interface{}, i string) bool {
	_, exists := (*m)[i]
	return exists
}

func (s *HelmConfig) ToYaml() string {
	y, err := yaml.Marshal(s)
	if err != nil {
		log.Fatalf("Yaml Marshalling error: %v", err)
	}
	return string(y)
}

func NewHelmConfig(dockerImageTag *string, secrets *map[string]interface{}) HelmConfig {
	var s map[string]interface{}
	// check that we are working with a newer Vault version
	// retrieve values from "data" submap in this case
	if mapIndexExists(secrets, "data") && mapIndexExists(secrets, "metadata") {
		s = (*secrets)["data"].(map[string]interface{})
	} else { // use returned values directly otherwise
		s = *secrets
	}
	sHelmConfig := SecretsHelmConfig{
		Datas: s,
	}
	var h HelmConfig
	if dockerImageTag != nil {
		i := DockerImageHelmConfig{
			Tag: *dockerImageTag,
		}
		h = HelmConfig{
			Image:   i,
			Secrets: sHelmConfig,
		}
		return h
	}
	h = HelmConfig{
		Secrets: sHelmConfig,
	}
	return h

}
