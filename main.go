package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	akeylessSourceToken      string
	akeylessDestinationToken string
	sourceGatewayURL         string
	destinationGatewayURL    string
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
		if sourceGatewayURL == "" {
			sourceGatewayURL = getEnvVar("SOURCE_GATEWAY_URL")
		}
		if destinationGatewayURL == "" {
			destinationGatewayURL = getEnvVar("DESTINATION_GATEWAY_URL")
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
	kubernetesCmd.Flags().StringVar(&sourceGatewayURL, "source-gateway-url", "", "Source gateway URL")
	kubernetesCmd.Flags().StringVar(&destinationGatewayURL, "destination-gateway-url", "", "Destination gateway URL")
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
