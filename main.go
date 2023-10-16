package main

import (
	"encoding/json"
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
	akeylessSourceToken      string
	akeylessDestinationToken string
	sourceGatewayConfigURL         string
	destinationGatewayConfigURL    string
	filterConfigFilePath     string
)

var timeout = 30000 * time.Millisecond

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
		run(akeylessSourceToken, akeylessDestinationToken, sourceGatewayConfigURL, destinationGatewayConfigURL, filterConfigFilePath)
	},
}

func run(akeylessSourceToken string, akeylessDestinationToken string, sourceGatewayConfigURL string, destinationGatewayConfigURL string, filterConfigFilePath string) {
	k8sAuthConfigs := lookupK8sAuthConfigs(akeylessSourceToken, sourceGatewayConfigURL)
	fmt.Println("Found", len(k8sAuthConfigs.K8SAuths), "k8s auth configs")
	
}

func lookupK8sAuthConfigs(sourceToken string, sourceGatewayConfigURL string) KubeAuthConfigs {
	url := sourceGatewayConfigURL + "/configs/k8s-auths"

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
