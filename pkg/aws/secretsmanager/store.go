package secretsmanager

import (
	"context"
	"errors"
	"strings"

	"get.porter.sh/plugin/aws/pkg/aws/awsconfig"
	"get.porter.sh/porter/pkg/secrets/plugins"
	"get.porter.sh/porter/pkg/secrets/plugins/host"
	"get.porter.sh/porter/pkg/tracing"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/smithy-go"
	"github.com/hashicorp/go-hclog"
	"go.opentelemetry.io/otel/attribute"
)

var _ plugins.SecretsProtocol = &Store{}

const (
	SecretKeyName = "secret"
)

// Store implements the backing store for secrets in AWS Secrets Manager.
type Store struct {
	logger    hclog.Logger
	config    awsconfig.Config
	client    *secretsmanager.Client
	hostStore host.Store
}

func NewStore(cfg awsconfig.Config, l hclog.Logger) *Store {
	return &Store{
		config:    cfg,
		logger:    l,
		hostStore: host.NewStore(),
	}
}

func (s *Store) Connect(ctx context.Context) error {
	_, log := tracing.StartSpan(ctx)
	defer log.EndSpan()
	if s.client != nil {
		return nil
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(s.config.Region))
	if err != nil {
		return err
	}

	client := secretsmanager.NewFromConfig(awsCfg)
	s.client = client

	return nil
}

func (s *Store) Resolve(ctx context.Context, keyName string, keyValue string) (string, error) {
	ctx, log := tracing.StartSpan(ctx)
	defer log.EndSpan()
	if strings.ToLower(keyName) != SecretKeyName {
		return s.hostStore.Resolve(ctx, keyName, keyValue)
	}

	log.SetAttributes(attribute.String("requested-secret", keyValue))

	if err := s.Connect(ctx); err != nil {
		return "", err
	}

	log.Debugf("getting secret %s from AWS secrets manager", keyValue)
	params := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(keyValue),
	}
	result, err := s.client.GetSecretValue(ctx, params)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}

// Create saves the secret to AWS Secrets Manager using the keyValue as the
// secret key.
// It implements the Create method on the secret plugins' interface.
func (s *Store) Create(ctx context.Context, keyName string, keyValue string, value string) error {
	ctx, log := tracing.StartSpan(ctx)
	defer log.EndSpan()
	if strings.ToLower(keyName) != SecretKeyName {
		return log.Errorf("unsupported secret type: %s. Only %s is supported", keyName, SecretKeyName)
	}

	if err := s.Connect(ctx); err != nil {
		return err
	}

	// Secret doesn't exist, proceed with creation
	log.Debugf("creating secret %s in AWS secrets manager", keyValue)
	_, err := s.client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(keyValue),
		SecretString: aws.String(value),
	})

	// If the error is not because the secret doesn't exist, return the error
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && apiErr.ErrorCode() != "ResourceExistsException" {
		return log.Errorf("failed to create secret %s in AWS secrets manager: %w", keyValue, err)
	}

	// Secret exists, update it
	log.Debugf("updating existing secret %s in AWS secrets manager", keyValue)
	_, err = s.client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(keyValue),
		SecretString: aws.String(value),
	})
	if err != nil {
		return log.Errorf("failed to update secret %s in AWS secrets manager: %w", keyValue, err)
	}
	return nil
}
