package oai

import (
	"github.com/go-logr/zapr"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/jlewi/bsctl/pkg/config"
	"github.com/jlewi/monogo/files"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

// NewClient helper function to create a new OpenAI client from  a config
func NewClient(cfg config.Config) (*openai.Client, error) {
	log := zapr.NewLogger(zap.L())
	// ************************************************************************
	// Setup middleware
	// ************************************************************************

	// Handle retryable errors
	// To handle retryable errors we use hashi corp's retryable client. This client will automatically retry on
	// retryable errors like 429; rate limiting
	retryClient := retryablehttp.NewClient()
	httpClient := retryClient.StandardClient()

	log.Info("Configuring OpenAI client")
	if cfg.OpenAI == nil {
		return nil, errors.New("OpenAI configuration is required")
	}

	apiKey := ""
	if cfg.OpenAI.APIKeyFile != "" {
		var err error
		raw, err := files.Read(cfg.OpenAI.APIKeyFile)
		if err != nil {
			return nil, err
		}
		apiKey = string(raw)
	}
	// If baseURL is customized then we could be using a custom endpoint that may not require an API key
	if apiKey == "" && cfg.OpenAI.BaseURL == "" {
		return nil, errors.New("OpenAI APIKeyFile is required when using OpenAI")
	}
	clientConfig := openai.DefaultConfig(apiKey)
	if cfg.OpenAI.BaseURL != "" {
		log.Info("Using custom OpenAI BaseURL", "baseURL", cfg.OpenAI.BaseURL)
		clientConfig.BaseURL = cfg.OpenAI.BaseURL
	}

	clientConfig.HTTPClient = httpClient
	client := openai.NewClientWithConfig(clientConfig)

	return client, nil
}
