package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDCProvider struct {
	provider  *oidc.Provider
	oauth2Cfg oauth2.Config
	verifier  *oidc.IDTokenVerifier
}

func NewOIDCProvider(ctx context.Context, tenantID, clientID, clientSecret, redirectURI string) (*OIDCProvider, error) {
	issuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tenantID)

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("could not discover OIDC provider at %s: %w", issuer, err)
	}

	return &OIDCProvider{
		provider: provider,
		oauth2Cfg: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  redirectURI,
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		verifier: provider.Verifier(&oidc.Config{ClientID: clientID}),
	}, nil
}

func (p *OIDCProvider) AuthCodeURL(state, nonce string) string {
	return p.oauth2Cfg.AuthCodeURL(state, oauth2.SetAuthURLParam("nonce", nonce))
}

func (p *OIDCProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return p.oauth2Cfg.Exchange(ctx, code)
}

func (p *OIDCProvider) VerifyIDToken(ctx context.Context, rawIDToken string) (*oidc.IDToken, error) {
	return p.verifier.Verify(ctx, rawIDToken)
}
