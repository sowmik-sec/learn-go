package auth

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, email, username, password string, roleIDs []string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
}

type RoleRepository interface {
	CreateRole(ctx context.Context, name, description string, permissionIDs []string) (*Role, error)
	GetRoleByID(ctx context.Context, id string) (*Role, error)
	GetRoleByName(ctx context.Context, name string) (*Role, error)
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, id string) error
	ListRoles(ctx context.Context, limit, offset int) ([]*Role, error)

	CreatePermission(ctx context.Context, resource, action, name string) (*Permission, error)
	GetPermissionByID(ctx context.Context, id string) (*Permission, error)
	GetPermissionsByResource(ctx context.Context, resource string) ([]*Permission, error)
	ListPermissions(ctx context.Context, limit, offset int) ([]*Permission, error)
}

type AuthService interface {
	// User management
	CreateUser(ctx context.Context, email, username, password string) (*User, error)

	AuthenticateUser(ctx context.Context, emailOrUsername, password string) (*User, error)
	GetUserByID(ctx context.Context, userId string) (*User, error)
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) bool

	HasPermission(ctx context.Context, userID, resource, action string) (bool, error)
	HasAnyPermission(ctx context.Context, userID string, permissions []PermissionCheck) (bool, error)
	GetUserPermissions(ctx context.Context, userID string) ([]*Permission, error)

	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]*Role, error)

	CreateRoleWithPermissions(ctx context.Context, name, description string, permissions []PermissionCheck) (*Role, error)
}

type PermissionCheck struct {
	Resource string
	Action   string
}
