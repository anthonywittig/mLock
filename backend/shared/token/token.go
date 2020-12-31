package token

import (
	"context"
	"mlock/shared"
	"mlock/shared/datastore"
	"net/http"

	googleAuthIDTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
)

type TokenData struct {
	Token      string
	TokenValid bool
	User       datastore.User
	UserValid  bool
	Error      error
}

// GetUserFromToken will first verify that the token is valid, then return the corresponding user object.
func GetUserFromToken(ctx context.Context, token string) (TokenData, error) {
	// Verify the token.
	v := googleAuthIDTokenVerifier.Verifier{}
	if err := v.VerifyIDToken(token, []string{shared.GetConfig("GOOGLE_SIGNIN_CLIENT_ID")}); err != nil {
		// For now we'll just assume the token is bad (could be network error etc.).
		return TokenData{Error: err}, nil
	}

	// Grab the claim set (to get the email).
	claimSet, err := googleAuthIDTokenVerifier.Decode(token)
	if err != nil {
		return TokenData{Error: err}, nil
	}

	// Grab the user.
	user, ok, err := datastore.GetUserByEmail(ctx, claimSet.Email)
	if err != nil {
		return TokenData{}, err
	}

	// In the future we could create our own token, but for now we'll just piggy back on Google's.
	tokenData := TokenData{Token: token, TokenValid: true}

	if !ok {
		tokenData.Error = shared.NewAPIError("user not authorized", http.StatusForbidden)
		return tokenData, nil
	}

	tokenData.User = user
	tokenData.UserValid = true
	return tokenData, nil
}
