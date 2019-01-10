package auth

import (
	"context"
	"os"
	"strings"

	"google.golang.org/api/option"

	"github.com/aws/aws-lambda-go/events"

	firebase "firebase.google.com/go"
	auth "firebase.google.com/go/auth"
	"github.com/pkg/errors"
)

func AuthFirebaseForLambda(request *events.APIGatewayProxyRequest) (*auth.Token, error) {
	var token *auth.Token

	authHeader, ok := request.Headers["Authorization"]
	if !ok {
		return token, errors.New("Authorization header not found")
	}

	return authFirebase(authHeader)
}

func authFirebase(authHeader string) (*auth.Token, error) {
	var token *auth.Token
	ctx := context.Background()
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)

	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_JSON_PATH"))
	var config *firebase.Config

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return token, errors.Wrap(err, "Failed to make firebase app")
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return token, errors.Wrap(err, "Failed to make firebase auth client")
	}

	token, err = client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return token, errors.Wrap(err, "Failed to verify firebase auth id token")
	}

	return token, nil
}
