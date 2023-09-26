package helm

import (
	"gopkg.in/yaml.v2"
	"kube-svc-ctl/utils"
	"log"
)

type Config struct {
	tag     *string
	secrets *SecretTree
}

func (s *Config) ToYaml(targetTree *string, format *string) string {
	var tree map[string]interface{}
	secrets := *s.secrets
	if *format == "list" {
		tree = secrets.MakeList(targetTree)
	} else {
		tree = secrets.MakeMap(targetTree)
	}

	if s.tag != nil {
		tree["image"] = map[string]map[string]string{}
		t := make(map[string]interface{})
		t["tag"] = s.tag
		tree["image"] = t
	}
	y, err := yaml.Marshal(tree)
	if err != nil {
		log.Fatalf("Yaml Marshalling error: %v", err)
	}
	return string(y)
}

func NewHelmConfig(dockerImageTag *string, secrets *map[string]interface{}) Config {
	var s map[string]interface{}
	// check that we are working with a newer Vault version
	// retrieve values from "data" submap in this case
	if utils.MapIndexExists(secrets, "data") && utils.MapIndexExists(secrets, "metadata") {
		s = (*secrets)["data"].(map[string]interface{})
	} else { // use returned values directly otherwise
		s = *secrets
	}
	var secretTree = SecretTree{secrets: s}
	sHelmConfig := Config{
		secrets: &secretTree,
		tag:     dockerImageTag,
	}
	return sHelmConfig

}
