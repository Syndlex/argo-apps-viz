package main

import (
	"argo-apps-viz/cmd/plugin/cli" // Updated import path

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	cli.InitAndExecute()
}
