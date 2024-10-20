package cli

import (
	"context"
	_ "embed"
	"github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned/typed/application/v1alpha1"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/syndlex/argo-apps-viz/pkg/logger"
	"github.com/syndlex/argo-apps-viz/pkg/model/appsofapps"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	appsOfAppsFile = "apps-of-apps.html"
)

var argoAppsOffApps = &cobra.Command{
	Use:   "apps-of-apps",
	Short: "Generate documentation from your ArgoCD applications and applicationsSets within your cluster",
	RunE: func(c *cobra.Command, args []string) error {
		logger := logger.NewLogger()
		logger.Info("Generating Dependency Chart :)")
		graph, err := runAoa(c.Flags())
		if err != nil {
			return err
		}
		err = CreateFile(appsOfAppsFile, graph)
		if err != nil {
			return err
		}
		logger.Info("Finished look in: " + appsOfAppsFile)
		return err
	},

	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(argoAppsOffApps)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// graphAoaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	argoAppsOffApps.Flags().StringArrayP("root", "r", []string{}, "When using a Recursive Apps of Apps pattern please add this here")
	argoAppsOffApps.Flags().BoolP("tree", "t", false, "Set this if you wand a tree instead of a graph")
}

func runAoa(flags *pflag.FlagSet) (components.Charter, error) {
	log := logger.NewLogger()
	config, err := KubernetesConfigFlags.ToRESTConfig()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	argoclient, err := v1alpha1.NewForConfig(config)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var ctx = context.Background()
	applicationList, err := argoclient.Applications("argocd").List(ctx, v1.ListOptions{})
	if err != nil {
		log.Info("Problem while getting ArgoCd domains")
		log.Error(err)
		return nil, err
	}

	applicationSetList, err := argoclient.ApplicationSets("argocd").List(ctx, v1.ListOptions{})
	if err != nil {
		log.Info("Problem while getting ArgoCd domains")
		log.Error(err)
		return nil, err
	}

	isTree, err := flags.GetBool("tree")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if isTree {
		roots, err := flags.GetStringArray("root")
		if err != nil {
			log.Error(err)
			return nil, err
		}
		return appsofapps.RenderTree(applicationSetList, applicationList, roots), nil
	} else {
		return appsofapps.AppsOfAppsRenderGraph(applicationSetList, applicationList), nil
	}
}
