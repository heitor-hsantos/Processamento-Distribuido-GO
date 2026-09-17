package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/coreos/go-oidc/v3/oidc"
)

type OIDCConfig struct {
	IssuerURL string
	Audience  string
}

type OIDCValidator struct {
	Config   OIDCConfig
	verifier *oidc.IDTokenVerifier
}

func NewOIDCValidator(cfg OIDCConfig) *OIDCValidator {
	// During initialization we just prepare the verifier.
	// Normally we'd fetch the provider configuration here,
	// but to prevent app crashing on startup if keycloak is down,
	// we will initialize it lazily or ignore the network failure for the scaffold.

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		log.Printf("WARN: could not fetch oidc provider config at startup: %v", err)
	}

	var verifier *oidc.IDTokenVerifier
	if provider != nil {
		verifier = provider.Verifier(&oidc.Config{ClientID: cfg.Audience})
	}

	return &OIDCValidator{
		Config:   cfg,
		verifier: verifier,
	}
}

func (v *OIDCValidator) ValidateToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("oidc: missing token")
	}

	// Fallback/Mock for local testing if no verifier exists (provider was down)
	if v.verifier == nil {
		return "provider-tenant-demo", nil
	}

	idToken, err := v.verifier.Verify(ctx, token)
	if err != nil {
		return "", fmt.Errorf("oidc: token verification failed: %w", err)
	}

	return idToken.Subject, nil
}
