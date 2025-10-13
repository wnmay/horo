package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

func InitFirebase(ctx context.Context, firebaseAccountKeyFile string) *auth.Client {
	opt := option.WithCredentialsFile(firebaseAccountKeyFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	firebaseAuthClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Failed to get firebaseAuthClient: %v", err)
	}
	return firebaseAuthClient
}
