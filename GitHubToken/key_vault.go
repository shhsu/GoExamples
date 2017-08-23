package azure

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"strings"

	"azure.com/acr/acr-build-runner/runner/domain"
	"github.com/Azure/azure-sdk-for-go/dataplane/keyvault"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
)

type KeyValutClient struct {
	config *domain.AzureConfig
}

func NewKeyValutClient(config *domain.AzureConfig) *KeyValutClient {
	return &KeyValutClient{config: config}
}

func (c *KeyValutClient) GetPrivateKey(kvConfig *domain.KeyVaultConfig, keyRef domain.KeyReference) (*rsa.PrivateKey, error) {
	cloud := c.config.Cloud
	oauthConfig, err := adal.NewOAuthConfig(cloud.ActiveDirectoryEndpoint, c.config.Tenant)
	if err != nil {
		return nil, err
	}

	resource := cloud.KeyVaultEndpoint
	if strings.HasSuffix(resource, "/") {
		resource = resource[0 : len(resource)-1]
	}

	sp := c.config.ServicePrincipal
	spToken, err := adal.NewServicePrincipalToken(*oauthConfig, sp.ID, sp.Password, resource)
	if err != nil {
		return nil, err
	}

	client := keyvault.New()
	client.Authorizer = autorest.NewBearerAuthorizer(spToken)
	keyBundle, err := client.GetKey(kvConfig.BaseURL, keyRef.Name, keyRef.Version)
	if err != nil {
		return nil, err
	}

	webKey := keyBundle.Key
	nValue, err := toBigInt(webKey.N)
	if err != nil {
		return nil, err
	}
	eBigInt, err := toBigInt(webKey.E)
	if err != nil {
		return nil, err
	}
	eValue := int(eBigInt.Int64())

	// NPE because webkey.D is nil
	var dValue *big.Int
	if webKey.D != nil {
		dValue, err = toBigInt(webKey.D)
		if err != nil {
			return nil, err
		}
	}

	return &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: nValue,
			E: eValue,
		},
		D:      dValue,
		Primes: []*big.Int{},
	}, nil
}

func toBigInt(webKeyValue *string) (*big.Int, error) {
	bytesValue, err := decodeWebKeyValue(*webKeyValue)
	if err != nil {
		return &big.Int{}, err
	}
	return (&big.Int{}).SetBytes(bytesValue), nil
}

func decodeWebKeyValue(str string) ([]byte, error) {
	remainder := len(str) % 4
	if remainder > 0 {
		str = str + strings.Repeat("=", 4-remainder)
	}
	return base64.URLEncoding.DecodeString(str)
}
