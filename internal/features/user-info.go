package features

import (
	"context"
	"github.com/bingemate/keycloak-service/initializers"
	"log"
)

type UserInfoService struct {
	keycloakClient *initializers.KeycloakClient
}

func NewUserInfoService(keycloakClient *initializers.KeycloakClient) *UserInfoService {
	return &UserInfoService{keycloakClient: keycloakClient}
}

func (s *UserInfoService) GetUsername(userID string) (string, error) {
	err := s.keycloakClient.EnsureToken(context.Background())
	if err != nil {
		log.Println("Error getting username", err)
		return "", err
	}
	userInfo, err := s.keycloakClient.Gocloak.GetUserByID(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		userID,
	)
	if err != nil {
		log.Println("Error getting username", err)
		return "", err
	}
	return *userInfo.Username, nil
}
