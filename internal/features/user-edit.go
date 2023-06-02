package features

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/bingemate/keycloak-service/initializers"
	"strings"
)

type UserEditService struct {
	keycloakClient  *initializers.KeycloakClient
	userInfoService *UserInfoService
}

func NewUserEditService(keycloakClient *initializers.KeycloakClient, userInfoService *UserInfoService) *UserEditService {
	return &UserEditService{keycloakClient: keycloakClient, userInfoService: userInfoService}
}

func (s *UserEditService) UpdateUser(userID string, username, firstname, lastname, email string) (*gocloak.User, error) {
	user, err := s.userInfoService.GetUser(userID)
	if err != nil {
		return nil, err
	}

	username = strings.TrimSpace(username)
	firstname = strings.TrimSpace(firstname)
	lastname = strings.TrimSpace(lastname)
	email = strings.TrimSpace(email)

	if username != "" {
		user.Username = &username
	}
	if firstname != "" {
		user.FirstName = &firstname
	}
	if lastname != "" {
		user.LastName = &lastname
	}
	if email != "" {
		user.Email = &email
	}
	err = s.keycloakClient.EnsureToken(context.Background())
	if err != nil {
		return nil, err
	}

	err = s.keycloakClient.Gocloak.UpdateUser(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		*user,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserEditService) UpdateUserPassword(userID string, newPassword string) error {
	_, err := s.userInfoService.GetUser(userID)
	if err != nil {
		return err
	}
	err = s.keycloakClient.EnsureToken(context.Background())
	if err != nil {
		return err
	}
	return s.keycloakClient.Gocloak.SetPassword(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		userID,
		s.keycloakClient.Realm,
		newPassword,
		false,
	)
}

func (s *UserEditService) AddUserRole(userID string, role string) error {
	_, err := s.userInfoService.GetUser(userID)
	if err != nil {
		return err
	}

	realmRole, err := s.keycloakClient.Gocloak.GetRealmRole(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		role,
	)
	if err != nil {
		return err
	}
	err = s.keycloakClient.EnsureToken(context.Background())
	if err != nil {
		return err
	}
	return s.keycloakClient.Gocloak.AddRealmRoleToUser(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		userID,
		[]gocloak.Role{*realmRole},
	)
}

func (s *UserEditService) RemoveUserRole(userID string, role string) error {
	_, err := s.userInfoService.GetUser(userID)
	if err != nil {
		return err
	}

	realmRole, err := s.keycloakClient.Gocloak.GetRealmRole(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		role,
	)
	if err != nil {
		return err
	}
	err = s.keycloakClient.EnsureToken(context.Background())
	if err != nil {
		return err
	}
	return s.keycloakClient.Gocloak.DeleteRealmRoleFromUser(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		userID,
		[]gocloak.Role{*realmRole},
	)
}

func (s *UserEditService) DeleteUser(userID string) error {
	_, err := s.userInfoService.GetUser(userID)
	if err != nil {
		return err
	}
	err = s.keycloakClient.EnsureToken(context.Background())
	if err != nil {
		return err
	}
	return s.keycloakClient.Gocloak.DeleteUser(
		context.Background(),
		s.keycloakClient.Token.AccessToken,
		s.keycloakClient.Realm,
		userID,
	)
}
