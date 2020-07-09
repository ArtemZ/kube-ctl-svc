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
	Values map[string]interface{}
}

func (s *HelmConfig) ToYaml() string {
	y, err := yaml.Marshal(s)
	if err != nil {
		log.Fatalf("Yaml Marshalling error: %v", err)
	}
	return string(y)
}

func NewHelmConfig(dockerImageTag *string, secrets *map[string]interface{}) HelmConfig {
	s := SecretsHelmConfig{
		Values: *secrets,
	}
	i := DockerImageHelmConfig{
		Tag: *dockerImageTag,
	}
	return HelmConfig{
		Image:   i,
		Secrets: s,
	}
}
