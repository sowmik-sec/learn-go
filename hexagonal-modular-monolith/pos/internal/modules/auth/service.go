package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo UserRepository
	roleRepo RoleRepository
}

func NewAuthService(userRepo UserRepository, roleRepo RoleRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *authService) CreateUser(ctx context.Context, email, username, password string) (*User, error) {
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return s.userRepo.CreateUser(ctx, email, username, hashedPassword, []string{})
}

func (s *authService) AuthenticateUser(ctx context.Context, emailOrUsername, password string) (*User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, emailOrUsername)
	if err != nil {
		user, err = s.userRepo.GetUserByUsername(ctx, emailOrUsername)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}
	}
	log.Printf("Auth: Found user %s, stored hash: %s", user.Username, user.Password)
	log.Printf("Auth: Input password: %s", password)
	if !s.VerifyPassword(password, user.Password) {
		log.Printf("Auth: Password verification failed")
		return nil, errors.New("invalid credentials")
	}
	log.Printf("Auth: Password verification succeeded")
	return user, nil
}

func (s *authService) HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password must be at least 8 characters long")
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

func (s *authService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *authService) HasPermission(ctx context.Context, userID, resource, action string) (bool, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	for _, roleID := range user.RoleIDs {
		role, err := s.roleRepo.GetRoleByID(ctx, roleID)
		if err != nil {
			continue
		}
		for _, permID := range role.PermissionIDs {
			permission, err := s.roleRepo.GetPermissionByID(ctx, permID)
			if err != nil {
				continue
			}
			if permission.Resource == resource && permission.Action == action {
				return true, nil
			}
		}
	}
	return false, nil
}

func (s *authService) HasAnyPermission(ctx context.Context, userID string, permissions []PermissionCheck) (bool, error) {
	for _, perm := range permissions {
		hasPerm, err := s.HasPermission(ctx, userID, perm.Resource, perm.Action)
		if err != nil {
			return false, err
		}
		if hasPerm {
			return true, nil
		}
	}
	return false, nil
}

func (s *authService) GetUserPermissions(ctx context.Context, userID string) ([]*Permission, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	var allPermissions []*Permission
	permissionSet := make(map[string]bool)

	for _, roleID := range user.RoleIDs {
		role, err := s.roleRepo.GetRoleByID(ctx, roleID)
		if err != nil {
			continue
		}
		for _, permID := range role.PermissionIDs {
			if permissionSet[permID] {
				continue
			}
			permission, err := s.roleRepo.GetPermissionByID(ctx, permID)
			if err != nil {
				continue
			}
			allPermissions = append(allPermissions, permission)
			permissionSet[permID] = true
		}
	}
	return allPermissions, nil
}

func (s *authService) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	_, err = s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("role does not exists: %w", err)
	}
	for _, existingRoleID := range user.RoleIDs {
		if existingRoleID == roleID {
			return errors.New("user already has this role")
		}
	}
	user.RoleIDs = append(user.RoleIDs, roleID)
	user.UpdatedAt = time.Now()
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *authService) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	for i, existingRoleID := range user.RoleIDs {
		if existingRoleID == roleID {
			user.RoleIDs = append(user.RoleIDs[:i], user.RoleIDs[i+1:]...)
			user.UpdatedAt = time.Now()
			return s.userRepo.UpdateUser(ctx, user)
		}
	}
	return errors.New("user does ot have this role")
}

func (s *authService) GetUserRoles(ctx context.Context, userID string) ([]*Role, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	var roles []*Role
	for _, roleID := range user.RoleIDs {
		role, err := s.roleRepo.GetRoleByID(ctx, roleID)
		if err != nil {
			continue
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (s *authService) CreateRoleWithPermissions(ctx context.Context, name, description string, permissions []PermissionCheck) (*Role, error) {
	var permissionIDs []string
	for _, perm := range permissions {
		existingPerms, err := s.roleRepo.GetPermissionsByResource(ctx, perm.Resource)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing permissions: %w", err)
		}

		var permID string
		for _, existingPerm := range existingPerms {
			if existingPerm.Action == perm.Action {
				permID = existingPerm.ID
				break
			}
		}
		if permID == "" {
			newPerm, err := s.roleRepo.CreatePermission(ctx, perm.Resource, perm.Action, fmt.Sprintf("%s %s", perm.Resource, perm.Action))
			if err != nil {
				return nil, fmt.Errorf("failed to create permission: %w", err)
			}
			permID = newPerm.ID
		}
		permissionIDs = append(permissionIDs, permID)
	}
	return s.roleRepo.CreateRole(ctx, name, description, permissionIDs)
}

func (s *authService) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}
