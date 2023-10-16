package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	akeylessSourceToken      string
	akeylessDestinationToken string
	sourceGatewayConfigURL         string
	destinationGatewayConfigURL    string
	filterConfigFilePath     string
)

var rootCmd = &cobra.Command{
	Use:   "agmigrator",
	Short: "Akeyless Gateway Configuration Migrator",
	Long:  `agmigrator is a CLI tool for migrating configuration from one Akeyless Gateway cluster to another.`,
}

var kubernetesCmd = &cobra.Command{
	Use:   "kubernetes",
	Short: "Migrate kubernetes auth configs from a source Akeyless gateway cluster to a destination Akeyless gateway cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if akeylessSourceToken == "" {
			akeylessSourceToken = getEnvVar("AKEYLESS_SOURCE_TOKEN")
		}
		if akeylessDestinationToken == "" {
			akeylessDestinationToken = getEnvVar("AKEYLESS_DESTINATION_TOKEN")
		}
		if sourceGatewayConfigURL == "" {
			sourceGatewayConfigURL = getEnvVar("SOURCE_GATEWAY_URL")
		}
		if destinationGatewayConfigURL == "" {
			destinationGatewayConfigURL = getEnvVar("DESTINATION_GATEWAY_URL")
		}
		if filterConfigFilePath == "" {
			filterConfigFilePath = getEnvVar("FILTER_CONFIG_FILE_PATH")
		}
		// Do Stuff Here
	},
}

func getEnvVar(name string) string {
	if !strings.HasPrefix(name, "AKEYLESS") {
		name = "AKEYLESS_" + name
	}
	return os.Getenv(name)
}

func init() {
	rootCmd.AddCommand(kubernetesCmd)

	kubernetesCmd.Flags().StringVar(&akeylessSourceToken, "akeyless-source-token", "", "Akeyless source token")
	kubernetesCmd.Flags().StringVar(&akeylessDestinationToken, "akeyless-destination-token", "", "Akeyless destination token")
	kubernetesCmd.Flags().StringVar(&sourceGatewayConfigURL, "source-gateway-config-url", "", "Source gateway Config URL")
	kubernetesCmd.Flags().StringVar(&destinationGatewayConfigURL, "destination-gateway-config-url", "", "Destination gateway config URL")
	kubernetesCmd.Flags().StringVar(&filterConfigFilePath, "filter-config-file-path", "", "Filter config file path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
