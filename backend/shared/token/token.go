package token

import (
	"database/sql"
	"mlock/shared"
	"mlock/shared/datastore"

	googleAuthIDTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
)

// GetUserFromToken will first verify that the token is valid, then return the corresponding user object.
func GetUserFromToken(db *sql.DB, token string) (datastore.User, error) {
	// Verify the token.
	v := googleAuthIDTokenVerifier.Verifier{}
	if err := v.VerifyIDToken(token, []string{shared.GetConfig("GOOGLE_SIGNIN_CLIENT_ID")}); err != nil {
		return datastore.User{}, err
	}

	// Grab the claim set (to get the email).
	claimSet, err := googleAuthIDTokenVerifier.Decode(token)
	if err != nil {
		return datastore.User{}, err
	}

	// Grab the user.
	return datastore.GetUserByEmail(db, claimSet.Email)
}
