package svc

import (
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

type HelmConfig struct {
	tag        *string
	secrets    *map[string]interface{}
	targetTree *string
}

func (s *HelmConfig) ToYaml() string {
	tree := make(map[string]interface{})

	if s.targetTree != nil {
		// makes map[value][value][value] out of
		// "value.value.value" string
		var lastBranch *map[string]interface{} = nil
		splitTargetTree := strings.Split(*s.targetTree, ".")
		for index, elem := range splitTargetTree {
			subtree := make(map[string]interface{})
			if lastBranch == nil {
				tree[elem] = subtree
			} else {
				(*lastBranch)[elem] = subtree
			}
			if len(splitTargetTree) == index+1 {
				(*lastBranch)[elem] = s.secrets
			} else {
				lastBranch = &subtree
			}
		}

	} else {
		tree = *s.secrets
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

func mapIndexExists(m *map[string]interface{}, i string) bool {
	_, exists := (*m)[i]
	return exists
}

func NewHelmConfig(dockerImageTag *string, secrets *map[string]interface{}, targetTree *string) HelmConfig {
	var s map[string]interface{}
	// check that we are working with a newer Vault version
	// retrieve values from "data" submap in this case
	if mapIndexExists(secrets, "data") && mapIndexExists(secrets, "metadata") {
		s = (*secrets)["data"].(map[string]interface{})
	} else { // use returned values directly otherwise
		s = *secrets
	}
	sHelmConfig := HelmConfig{
		secrets:    &s,
		tag:        dockerImageTag,
		targetTree: targetTree,
	}
	return sHelmConfig

}
