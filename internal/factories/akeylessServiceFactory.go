package factories

import (
	"akeyless-gateway-migrator/migrator/internal/services"

	"github.com/akeylesslabs/akeyless-go/v2"
)

func BuildAkeylessService(url ...string) *services.AkeylessService {

	urlString := "https://api.akeyless.io"
	if url != nil {
		urlString = url[0] + "/api/v2"
	}
	client := akeyless.NewAPIClient(&akeyless.Configuration{
		Servers: []akeyless.ServerConfiguration{
			{
				URL: urlString,
			},
		},
	}).V2Api

	return services.NewAkeylessService(client)
}
