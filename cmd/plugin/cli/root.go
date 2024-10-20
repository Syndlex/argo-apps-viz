package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"strings"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
	rootCmd               = &cobra.Command{
		Use:           "argo-apps-viz",
		Short:         "Hi, you need to execute one of the sub commands",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "kubectl argo-apps-viz",
		},
	}
)

func RootCmd() *cobra.Command {

	cobra.OnInitialize(initConfig)

	KubernetesConfigFlags = genericclioptions.NewConfigFlags(false)
	KubernetesConfigFlags.AddFlags(rootCmd.Flags())

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return rootCmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
