package helm

import "testing"

func TestHelmConfig_ToYaml(t *testing.T) {
	dockerTag := "latest"
	secrets := make(map[string]interface{})
	secrets["secret1"] = 1
	secrets["secret2"] = "b"
	tree := "some.weird.tree"

	c := NewHelmConfig(&dockerTag, &secrets, &tree)

	got := c.ToYaml()
	want := `image:
  tag: latest
some:
  weird:
    tree:
      secret1: 1
      secret2: b
`
	if got != want {
		t.Fatalf("Yaml not matching, got %q, want %q", got, want)
	}
	println(got)
}
