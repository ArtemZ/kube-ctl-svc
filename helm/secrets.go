package helm

import (
	"fmt"
	"strings"
)

type SecretTree struct {
	secrets map[string]interface{}
}

func (s *SecretTree) MakeMap(targetTree *string) map[string]interface{} {
	splitTargetTree := strings.Split(*targetTree, ".")
	var lastChild map[string]interface{}
	for index := range splitTargetTree {
		index = len(splitTargetTree) - 1 - index // reverse order

		if index == len(splitTargetTree)-1 { // last element
			lastChild = make(map[string]interface{})
			lastChild[splitTargetTree[index]] = s.secrets
		} else {
			lastChild[splitTargetTree[index]] = lastChild
		}
	}
	return lastChild
}

func (s *SecretTree) MakeList(targetTree *string) map[string]interface{} {
	var secretList = make([]map[string]string, 0)
	for k, v := range s.secrets {
		value := map[string]string{
			"key":   k,
			"value": fmt.Sprint(v),
		}
		secretList = append(secretList, value)
	}
	splitTargetTree := strings.Split(*targetTree, ".")
	var lastChild map[string]interface{}
	for index := range splitTargetTree {
		index = len(splitTargetTree) - 1 - index // reverse order

		if index == len(splitTargetTree)-1 { // last element
			lastChild = make(map[string]interface{})
			lastChild[splitTargetTree[index]] = secretList
		} else {
			lastChild[splitTargetTree[index]] = lastChild
		}
	}
	return lastChild
}
