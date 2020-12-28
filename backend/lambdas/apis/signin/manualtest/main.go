package main

import (
	"fmt"

	googleAuthIDTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
)

const (
	// Don't commit actual values.
	GOOGLE_SIGNIN_CLIENT_ID = ""
	TOKEN                   = ""
)

func main() {
	fmt.Println("starting...")

	email, err := getEmail2(TOKEN)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	fmt.Printf("-- %s --\n", email)
}

func getEmail2(token string) (string, error) {

	v := googleAuthIDTokenVerifier.Verifier{}
	if err := v.VerifyIDToken(token, []string{GOOGLE_SIGNIN_CLIENT_ID}); err != nil {
		return "", err
	}

	claimSet, err := googleAuthIDTokenVerifier.Decode(token)
	if err != nil {
		return "", err
	}

	return claimSet.Email, nil
}
