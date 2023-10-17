package services

import (
	"context"
	"encoding/base64"

	"github.com/akeylesslabs/akeyless-go/v2"
)

type AkeylessService struct {
	client *akeyless.V2ApiService
}

func NewAkeylessService(client *akeyless.V2ApiService) *AkeylessService {
	return &AkeylessService{
		client: client,
	}
}

type K8SAuthMethod struct {
	AccessId   string
	PrivateKey string
}

type K8SDetails struct {
	K8SHost                string
	KubeCACert             string
	K8SIssuer              string
	K8SServiceAccountToken string
}

func (service *AkeylessService) ValidateToken(ctx context.Context, token string) (*akeyless.ValidateTokenOutput, error) {
	validateToken := akeyless.ValidateToken{
		Token: &token,
	}
	output, _, err := service.client.ValidateToken(ctx).Body(validateToken).Execute()
	return &output, err
}

func (service *AkeylessService) CreateAuthConfigK8S(ctx context.Context, authConfigName string, k8sAuthMethod *K8SAuthMethod, k8sDetails *K8SDetails, token string) (string, error) {

	createK8sAuthConfigBody := akeyless.NewGatewayCreateK8SAuthConfig(k8sAuthMethod.AccessId, k8sDetails.K8SHost, authConfigName, k8sAuthMethod.PrivateKey)

	base64CACert := base64.StdEncoding.EncodeToString([]byte(k8sDetails.KubeCACert))
	createK8sAuthConfigBody.SetK8sCaCert(base64CACert)
	createK8sAuthConfigBody.SetToken(token)
	createK8sAuthConfigBody.SetK8sIssuer(k8sDetails.K8SIssuer)
	createK8sAuthConfigBody.SetTokenReviewerJwt(k8sDetails.K8SServiceAccountToken)

	_, _, err := service.client.GatewayCreateK8SAuthConfig(ctx).Body(*createK8sAuthConfigBody).Execute()
	if err != nil {
		return "", err
	}

	return authConfigName, nil

}
