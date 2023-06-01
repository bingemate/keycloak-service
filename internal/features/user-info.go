package features

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
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

func (s *UserInfoService) SearchUsers(username string, includeRoles bool) ([]*gocloak.User, error) {
	err := s.keycloakClient.EnsureToken(context.Background())
	if err != nil {
		log.Println("Error getting username", err)
		return nil, err
	}
	users, err := s.keycloakClient.Gocloak.GetUsers(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		gocloak.GetUsersParams{
			Search: &username,
			Exact:  gocloak.BoolP(false),
		},
	)
	if err != nil {
		log.Println("Error getting username", err)
		return nil, err
	}
	if users == nil {
		return []*gocloak.User{}, nil
	}

	if includeRoles {

		for _, user := range users {
			roles, err := s.keycloakClient.Gocloak.GetRealmRolesByUserID(
				context.Background(),
				s.keycloakClient.Token.AccessToken,
				s.keycloakClient.Realm,
				*user.ID,
			)
			if err != nil {
				log.Println("Error getting username", err)
				return nil, err
			}
			rolesArray := make([]string, len(roles))
			for i, role := range roles {
				rolesArray[i] = *role.Name
			}
			user.RealmRoles = &rolesArray
		}
	}
	if err != nil {
		log.Println("Error getting username", err)
		return nil, err
	}

	return users, nil
}
