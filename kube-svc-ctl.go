package main

import (
	"kube-svc-ctl/cmd/actions"
)
import "flag"
import "os"
import "fmt"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("One of the following commands is required: validate|generate-helm-values|generate-secret-manifest")
	}
	switch os.Args[1] {
	case "validate":
		actions.ValidateCommand()
	case "generate-svc-config":
		err := actions.GenerateHelmValues()
		if err != nil {
			fmt.Println(err)
		}
	case "generate-helm-values":
		err := actions.GenerateHelmValues()
		if err != nil {
			fmt.Println(err)
		}
	case "generate-secret-manifest":
		err := actions.GenerateSecretManifest()
		if err != nil {
			fmt.Println(err)
		}
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

}
