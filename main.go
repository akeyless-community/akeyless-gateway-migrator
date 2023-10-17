package main

import (
	"akeyless-gateway-migrator/migrator/internal/factories"
	"akeyless-gateway-migrator/migrator/internal/services"
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gojek/heimdall/httpclient"
	"github.com/spf13/cobra"
)

type KubeAuthConfig struct {
	Name                 string `json:"name,omitempty"`
	ID                   string `json:"id,omitempty"`
	ProtectionKey        string `json:"protection_key,omitempty"`
	AuthMethodAccessID   string `json:"auth_method_access_id,omitempty"`
	AuthMethodPrvKeyPem  string `json:"auth_method_prv_key_pem,omitempty"`
	AmTokenExpiration    int    `json:"am_token_expiration,omitempty"`
	K8SHost              string `json:"k8s_host,omitempty"`
	K8SCaCert            string `json:"k8s_ca_cert,omitempty"`
	K8STokenReviewerJwt  string `json:"k8s_token_reviewer_jwt,omitempty"`
	K8SIssuer            string `json:"k8s_issuer,omitempty"`
	K8SPubKeysPem        string `json:"k8s_pub_keys_pem,omitempty"`
	DisableIssValidation bool   `json:"disable_iss_validation,omitempty"`
	UseLocalCaJwt        bool   `json:"use_local_ca_jwt,omitempty"`
	ClusterAPIType       string `json:"cluster_api_type,omitempty"`
}

type KubeAuthConfigs struct {
	K8SAuths []KubeAuthConfig `json:"k8s_auths,omitempty"`
}

var (
	akeylessSourceToken         string
	akeylessDestinationToken    string
	sourceGatewayConfigURL      string
	destinationGatewayConfigURL string
	filterConfigFilePath        string
	debugFlag                   bool
)

var timeout = 30000 * time.Millisecond

var rootCmd = &cobra.Command{
	Use:   "agmigrator",
	Short: "Akeyless Gateway Migrator",
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
		run(akeylessSourceToken, akeylessDestinationToken, sourceGatewayConfigURL, destinationGatewayConfigURL, filterConfigFilePath)
	},
}

func run(akeylessSourceToken string, akeylessDestinationToken string, sourceGatewayConfigURL string, destinationGatewayConfigURL string, filterConfigFilePath string) {
	if debugFlag {
		// output all the argument values
		fmt.Println("akeylessSourceToken:", akeylessSourceToken)
		fmt.Println("sourceGatewayConfigURL:", sourceGatewayConfigURL)
		fmt.Println("akeylessDestinationToken:", akeylessDestinationToken)
		fmt.Println("destinationGatewayConfigURL:", destinationGatewayConfigURL)
		fmt.Println("filterConfigFilePath:", filterConfigFilePath)
	}

	if akeylessSourceToken != "" {
		fmt.Println("Validating source token")
		runValidateToken(akeylessSourceToken, sourceGatewayConfigURL)
	}

	if akeylessDestinationToken != "" {
		fmt.Println("Validating destination token")
		runValidateToken(akeylessDestinationToken, destinationGatewayConfigURL)
	}

	k8sAuthConfigs := lookupK8sAuthConfigs(akeylessSourceToken, sourceGatewayConfigURL)

	fmt.Println("Found", len(k8sAuthConfigs.K8SAuths), "k8s auth configs")

	if filterConfigFilePath != "" {
		file, err := os.Open(filterConfigFilePath)
		if err != nil {
			fmt.Println("Unable to open filter config file:", filterConfigFilePath, err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		filterConfigNames := make(map[string]bool)
		for scanner.Scan() {
			filterConfigNames[scanner.Text()] = true
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading filter config file:", filterConfigFilePath, err)
			return
		}

		for _, k8sAuthConfig := range k8sAuthConfigs.K8SAuths {
			if _, ok := filterConfigNames[k8sAuthConfig.Name]; ok {
				// Migrate this config
				fmt.Println("Migrating config:", k8sAuthConfig.Name)
				runMigrate(k8sAuthConfig)
			}
		}
	} else {
		for _, k8sAuthConfig := range k8sAuthConfigs.K8SAuths {
			// Migrate this config
			fmt.Println("Migrating config:", k8sAuthConfig.Name)
			runMigrate(k8sAuthConfig)
		}
	}
}

// Check that each of the tokens are valid
func runValidateToken(akeylessToken string, gatewayConfigURL string) error {
	akeylessService := factories.BuildAkeylessService(sourceGatewayConfigURL)
	validateToken, err := akeylessService.ValidateToken(context.Background(), akeylessSourceToken)
	if err != nil {
		fmt.Println("Unable to validate token at URL:", gatewayConfigURL, err)
		return err
	} else {
		if validateToken.IsValid != nil && !*validateToken.IsValid {
			fmt.Println("Token is not valid")
			return errors.New("token is not valid")
		}
		// print line to indicate if token is valid and the token expiration time
		expirationTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", *validateToken.Expiration)
		if err != nil {
			fmt.Println("Unable to parse expiration time:", err)
			return err
		}
		// print the difference between the current time and the token expiration time, rounded to the nearest second
		fmt.Println("Token expires in:", time.Until(expirationTime).Round(time.Second))
	}
	return err
}

func lookupK8sAuthConfigs(sourceToken string, sourceGatewayConfigURL string) KubeAuthConfigs {
	url := sourceGatewayConfigURL + "/config/k8s-auths"

	var k8sAuthConfigs KubeAuthConfigs

	httpRequestClient := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

	// Create an http.Request instance
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	bearerToken := "Bearer " + sourceToken
	req.Header.Add("Authorization", bearerToken)
	// Call the `Do` method, which has a similar interface to the `http.Do` method
	res, err := httpRequestClient.Do(req)
	if err != nil {
		fmt.Println("Unable to get k8s auth configs:", sourceGatewayConfigURL, err)
		return generateEmptyK8sAuthConfigs()
	}

	body, err := ioutil.ReadAll(res.Body)

	err2 := json.Unmarshal(body, &k8sAuthConfigs)
	if err2 != nil {
		fmt.Println(err)
	}

	return k8sAuthConfigs
}

func generateEmptyK8sAuthConfigs() KubeAuthConfigs {
	k8sAuthConfigs := KubeAuthConfigs{
		K8SAuths: []KubeAuthConfig{},
	}
	return k8sAuthConfigs
}

func runMigrate(k8sAuthConfig KubeAuthConfig) {
	akeylessService := factories.BuildAkeylessService(destinationGatewayConfigURL)

	k8sAuthMethod := &services.K8SAuthMethod{
		AccessId:   k8sAuthConfig.AuthMethodAccessID,
		PrivateKey: k8sAuthConfig.AuthMethodPrvKeyPem,
	}

	k8sDetails := &services.K8SDetails{
		K8SHost:             k8sAuthConfig.K8SHost,
		K8SIssuer:           k8sAuthConfig.K8SIssuer,
		KubeCACert:          k8sAuthConfig.K8SCaCert,
		K8SServiceAccountToken: k8sAuthConfig.K8STokenReviewerJwt,
	}

	_, err := akeylessService.CreateAuthConfigK8S(context.Background(), k8sAuthConfig.Name, k8sAuthMethod, k8sDetails, akeylessDestinationToken)
	if err != nil {
		fmt.Println("Failed to migrate config:", k8sAuthConfig.Name, err)
	} else {
		fmt.Println("Successfully migrated config:", k8sAuthConfig.Name)
	}
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
	kubernetesCmd.Flags().BoolVar(&debugFlag, "debug", false, "Enable debug mode")
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
