package cli

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned/typed/application/v1alpha1"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/spf13/cobra"
	"github.com/syndlex/argo-apps-viz/pkg/logger"
	"github.com/syndlex/argo-apps-viz/pkg/model"
	"github.com/syndlex/argo-apps-viz/pkg/model/dependencies"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	dependenciesFile = "dependencies.html"
)

// graphAoaCmd represents the graphAoa command
var graphdepsCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "Generate dependency documentation from your ArgoCD applications and applicationsSets within your cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := logger.NewLogger()
		logger.Info("Generating Apps of Apps Chart :)")
		tree, err := runPlugin()
		if err != nil {
			return err
		}
		err = CreateFile(dependenciesFile, tree)
		if err != nil {
			return err
		}
		logger.Info("Finished look in: " + dependenciesFile)
		return err
	},
}

func init() {
	rootCmd.AddCommand(graphdepsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// graphAoaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// graphAoaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runPlugin() (*charts.Tree, error) {
	config, err := KubernetesConfigFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	argoclient, err := v1alpha1.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	var ctx = context.Background()
	applicationList, err := argoclient.Applications("argocd").List(ctx, v1.ListOptions{})
	if err != nil {
		fmt.Println("Problem while getting ArgoCD domains")
		return nil, err
	}

	applicationSetList, err := argoclient.ApplicationSets("argocd").List(ctx, v1.ListOptions{})
	if err != nil {
		fmt.Println("Problem while getting ArgoCd domains")
		return nil, err
	}

	c := model.Model(applicationList, applicationSetList)
	return dependencies.RenderTree(c), nil
}
