package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"azure.com/acr/acr-build-runner/runner/azure"
	"azure.com/acr/acr-build-runner/runner/domain"
)

func main() {
	config, err := domain.LoadConfig("conf.dev.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	key, err := readPKFromFile()
	if err != nil {
		panic(fmt.Errorf("Failed to retrieve private key: %s", err))
	}
	kvClient := azure.NewKeyValutClient(&config.Azure)
	cloudKey, err := kvClient.GetPrivateKey(&config.KeyVault, config.GitHubIntegration.IntegrationKey)
	fmt.Println(cloudKey.N.Bytes())
	fmt.Println(key.N.Bytes())
	fmt.Println(key.D.Bytes())
	// githubClient := github.NewIntegrationAppClient(&config.GitHubIntegration)
	// githubToken, err := githubClient.GetAccessToken(key, 47558)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(githubToken.Token)
	// token, err := github.GetGithubToken(azure.PublicCloud, &config.ServicePrincipal,
	// 	config.AzureTenantID, config.AzureSubscriptionID,
	// 	"https://shhsukey.vault.azure.net",
	// 	"githubintkeyportal", "95d0e5d681744d0a9587461a32886686",
	// 	// "githubintsecret", "91bbe99df9e341c796a9aa6b1efa3f20",
	// 	// "https://acrbuilderdemokv.vault.azure.net", "acrbuildergithub", "",
	// 	// 1838, 47772,
	// 	4515, 47558,
	// )
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// DELETEME
func readPKFromFile() (*rsa.PrivateKey, error) {
	keyBytes, err := ioutil.ReadFile("C:\\Users\\shhsu\\Downloads\\acr-build-demo-github-app.2017-08-21.private-key.pem")
	if err != nil {
		return nil, err
	}
	//fmt.Println(text)
	keyObj, err := loadPrivateKey(keyBytes)
	key, ok := keyObj.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("no cast")
	}
	return key, nil
}

func loadPrivateKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	var priv interface{}
	priv, err0 := x509.ParsePKCS1PrivateKey(input)
	if err0 == nil {
		return priv, nil
	}

	priv, err1 := x509.ParsePKCS8PrivateKey(input)
	if err1 == nil {
		return priv, nil
	}

	priv, err2 := x509.ParseECPrivateKey(input)
	if err2 == nil {
		return priv, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s', '%s' and '%s'", err0, err1, err2)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
