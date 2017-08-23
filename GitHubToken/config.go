package domain

type Config struct {
	KeyVault          KeyVaultConfig          `json:"keyvault"`
	GitHubIntegration GitHubIntegrationConfig `json:"github-integration"`
}


type AzureConfig struct {
	CloudName        string            `json:"cloud-name,omitempty"`
	Cloud            azure.Environment `json:"cloud,omitempty"`
	ServicePrincipal ServicePrincipal  `json:"service-principal"`
	Tenant           string            `json:"tenant"`
	Subscription     string            `json:"subscription"`
}

type ServicePrincipal struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type KeyVaultConfig struct {
	BaseURL string `json:"baseURL"`
}

type KeyReference struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type GitHubIntegrationConfig struct {
	AppID          int          `json:"app-id"`
	IntegrationKey KeyReference `json:"integration-key"`
}
