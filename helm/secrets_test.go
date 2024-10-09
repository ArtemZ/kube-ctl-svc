package helm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecrets_ToMap(t *testing.T) {
	secrets := make(map[string]interface{})
	secrets["secret1"] = 1
	secrets["secret2"] = "b"
	tree := "some.weird.tree"
	s := SecretTree{secrets: secrets}
	secretTree := s.MakeMap(&tree)
	assert.Contains(t, secretTree, "some")
	assert.Contains(t, secretTree["some"], "weird")
	var sMap = secretTree["some"].(map[string]interface{})
	assert.Contains(t, sMap, "weird")
	var tMap = sMap["weird"].(map[string]interface{})
	assert.Contains(t, tMap, "tree")
	var vMap = tMap["tree"].(map[string]interface{})
	assert.Equal(t, vMap["secret2"], "b")
}

func TestSecrets_ToList(t *testing.T) {
	secrets := make(map[string]interface{})
	secrets["secret1"] = 1
	secrets["secret2"] = "b"
	tree := "tree"

	s := SecretTree{secrets: secrets}
	print(fmt.Sprint(s.MakeList(&tree)))
}
