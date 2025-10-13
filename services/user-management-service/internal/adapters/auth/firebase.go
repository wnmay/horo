package app

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/wnmay/horo/services/user-management-service/internal/config"
	"google.golang.org/api/option"
)

func InitFirebase(ctx context.Context, cfg *config.Config) *auth.Client {
	var firebaseKeyfile = cfg.FirebaseAccountKeyFile
	opt := option.WithCredentialsFile(firebaseKeyfile)
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
