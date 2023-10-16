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
			akeylessSourceToken = os.Getenv("AKEYLESS_SOURCE_TOKEN")
		}
		if akeylessDestinationToken == "" {
			akeylessDestinationToken = os.Getenv("AKEYLESS_DESTINATION_TOKEN")
		}
		if sourceGatewayURL == "" {
			sourceGatewayURL = os.Getenv("SOURCE_GATEWAY_URL")
		}
		if destinationGatewayURL == "" {
			destinationGatewayURL = os.Getenv("DESTINATION_GATEWAY_URL")
		}
		if filterConfigFilePath == "" {
			filterConfigFilePath = os.Getenv("FILTER_CONFIG_FILE_PATH")
		}
		// Do Stuff Here
	},
}

func init() {
	rootCmd.AddCommand(kubernetesCmd)

	kubernetesCmd.Flags().StringVar(&akeylessSourceToken, "akeyless-source-token", "", "Akeyless source token")
	kubernetesCmd.Flags().StringVar(&akeylessDestinationToken, "akeyless-destination-token", "", "Akeyless destination token")
	kubernetesCmd.Flags().StringVar(&sourceGatewayURL, "source-gateway-url", "", "Source gateway URL")
	kubernetesCmd.Flags().StringVar(&destinationGatewayURL, "destination-gateway-url", "", "Destination gateway URL")
	kubernetesCmd.Flags().StringVar(&filterConfigFilePath, "filter-config-file-path", "", "Filter config file path")

	viper.BindPFlag("akeyless-source-token", kubernetesCmd.Flags().Lookup("akeyless-source-token"))
	viper.BindPFlag("akeyless-destination-token", kubernetesCmd.Flags().Lookup("akeyless-destination-token"))
	viper.BindPFlag("source-gateway-url", kubernetesCmd.Flags().Lookup("source-gateway-url"))
	viper.BindPFlag("destination-gateway-url", kubernetesCmd.Flags().Lookup("destination-gateway-url"))
	viper.BindPFlag("filter-config-file-path", kubernetesCmd.Flags().Lookup("filter-config-file-path"))

	viper.AutomaticEnv()
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
