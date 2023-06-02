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
			roles, err := s.keycloakClient.Gocloak.GetCompositeRealmRolesByUserID(
				context.Background(),
				s.keycloakClient.Token.AccessToken,
				s.keycloakClient.Realm,
				*user.ID,
			)
			if err != nil {
				log.Println("Error getting user's roles", err)
				return nil, err
			}
			rolesArray := make([]string, len(roles))
			for i, role := range roles {
				rolesArray[i] = *role.Name
			}
			user.RealmRoles = &rolesArray
		}
	}

	return users, nil
}

func (s *UserInfoService) GetUser(userID string) (*gocloak.User, error) {
	err := s.keycloakClient.EnsureToken(context.Background())
	if err != nil {
		log.Println("Error getting username", err)
		return nil, err
	}
	user, err := s.keycloakClient.Gocloak.GetUserByID(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		userID,
	)
	if err != nil {
		log.Println("Error getting user", err)
		return nil, err
	}
	roles, err := s.keycloakClient.Gocloak.GetCompositeRealmRolesByUserID(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		*user.ID,
	)
	if err != nil {
		log.Println("Error getting user's roles", err)
		return nil, err
	}
	rolesArray := make([]string, len(roles))
	for i, role := range roles {
		rolesArray[i] = *role.Name
	}
	user.RealmRoles = &rolesArray
	return user, nil
}

func (s *UserInfoService) GetUsers(query string, page, limit int) ([]*gocloak.User, int, error) {
	err := s.keycloakClient.EnsureToken(context.Background())
	if err != nil {
		log.Println("Error refreshing access token", err)
		return nil, 0, err
	}
	params := gocloak.GetUsersParams{
		First: gocloak.IntP((page - 1) * limit),
		Max:   gocloak.IntP(limit),
	}
	if query != "" {
		params.Search = &query
		params.Exact = gocloak.BoolP(false)
	}

	users, err := s.keycloakClient.Gocloak.GetUsers(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		params,
	)
	if err != nil {
		log.Println("Error getting users", err)
		return nil, 0, err
	}
	for _, user := range users {
		roles, err := s.keycloakClient.Gocloak.GetCompositeRealmRolesByUserID(
			context.Background(),
			s.keycloakClient.Token.AccessToken,
			s.keycloakClient.Realm,
			*user.ID,
		)
		if err != nil {
			log.Println("Error getting user's roles", err)
			return nil, 0, err
		}
		rolesArray := make([]string, len(roles))
		for i, role := range roles {
			rolesArray[i] = *role.Name
		}
		user.RealmRoles = &rolesArray
	}
	count, err := s.keycloakClient.Gocloak.GetUserCount(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		params,
	)
	if err != nil {
		log.Println("Error getting user count", err)
		return nil, 0, err
	}

	return users, count, nil
}

func (s *UserInfoService) GetAvailableRoles() (*[]string, error) {
	err := s.keycloakClient.EnsureToken(context.Background())
	if err != nil {
		return nil, err
	}
	roles, err := s.keycloakClient.Gocloak.GetRealmRoles(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		gocloak.GetRoleParams{},
	)
	if err != nil {
		return nil, err
	}
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = *role.Name
	}
	return &roleNames, nil
}
