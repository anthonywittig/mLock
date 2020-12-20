package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mlock/shared"
	"mlock/shared/datastore"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	googleAuthIDTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
)

type CreateBody struct {
	GoogleToken string
}

type CreateResponse struct {
	Token string
}

func main() {
	if err := shared.LoadConfig(); err != nil {
		fmt.Printf("ERROR loading config: %s\n", err.Error())
		return
	}
	lambda.Start(HandleRequest)
}

/*
Need a GCP session for this to work.
func getEmail(token string) (string, error) {
	json := []byte{}
	tokenValidator, err := idtoken.NewValidator(context.Background(), option.WithCredentialsJSON(json))
	if err != nil {
		return "", err
	}

	googleClientID := shared.GetConfig("GOOGLE_SIGNIN_CLIENT_ID")

	payload, err := tokenValidator.Validate(context.Background(), token, googleClientID)
	if err != nil {
		return "", err
	}

	email, ok := (payload.Claims["email"]).(string)
	if !ok {
		return "", fmt.Errorf("couldn't get email from token")
	}

	return email, nil
}
*/

/*
Never tried this one out.
func getEmail(idToken string) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.New(httpClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
}
*/

func getEmail(token string) (string, error) {
	v := googleAuthIDTokenVerifier.Verifier{}
	if err := v.VerifyIDToken(token, []string{shared.GetConfig("GOOGLE_SIGNIN_CLIENT_ID")}); err != nil {
		return "", err
	}

	claimSet, err := googleAuthIDTokenVerifier.Decode(token)
	if err != nil {
		return "", err
	}

	return claimSet.Email, nil
}

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "POST":
		return create(ctx, req)
	default:
		return shared.APIResponse(http.StatusNotImplemented, fmt.Errorf("not implemented - %s", req.HTTPMethod))
	}
}

func create(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var body CreateBody
	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	// Validate and get email.
	email, err := getEmail(body.GoogleToken)
	if err != nil {
		return nil, fmt.Errorf("error getting email: %s", err.Error())
	}

	// Verify that we have the user in our DB.
	user, err := datastore.GetUserByEmail(nil, email)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %s", err.Error())
	}

	return shared.APIResponse(http.StatusOK, CreateResponse{Token: user.Email})
}
