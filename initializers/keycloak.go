package initializers

import (
	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"log"
	"time"
)

type KeycloakClient struct {
	clientID          string
	clientSecret      string
	Realm             string
	Gocloak           gocloak.GoCloak
	Token             *gocloak.JWT
	tokenObtainedTime time.Time
}

func ConnectToKeycloak(env Env) (*KeycloakClient, error) {
	client := gocloak.NewClient(env.KeycloakUrl)
	token, err := client.LoginClient(
		context.Background(),
		env.ClientId,
		env.ClientSecret,
		env.Realm,
	)
	if err != nil {
		log.Println("Error connecting to keycloak", err)
		return &KeycloakClient{}, err
	}
	log.Println("Connected to keycloak")
	return &KeycloakClient{
		clientID:          env.ClientId,
		clientSecret:      env.ClientSecret,
		Realm:             env.Realm,
		Gocloak:           *client,
		Token:             token,
		tokenObtainedTime: time.Now(),
	}, nil
}

func (kc *KeycloakClient) EnsureToken(ctx context.Context) error {
	if kc.Token == nil {
		return fmt.Errorf("no token available")
	}

	if time.Since(kc.tokenObtainedTime).Seconds() < float64(kc.Token.ExpiresIn)*time.Second.Seconds() {
		return nil
	}

	newToken, err := kc.Gocloak.RefreshToken(ctx, kc.Token.RefreshToken, kc.clientID, kc.clientSecret, kc.Realm)
	if err != nil {
		newToken, err = kc.Gocloak.LoginClient(ctx, kc.clientID, kc.clientSecret, kc.Realm)
		if err != nil {
			return fmt.Errorf("failed to login: %w", err)
		}
	}
	kc.Token = newToken
	kc.tokenObtainedTime = time.Now()
	return nil
}
