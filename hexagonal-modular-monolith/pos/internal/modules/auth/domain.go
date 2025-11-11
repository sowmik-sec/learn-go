package auth

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	RoleIDs   []string  `json:"role_ids"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Role struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	PermissionIDs []string  `json:"permission_ids"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Permission struct {
	ID       string `json:"id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Name     string `json:"name"`
}

type RoleName string

const (
	RoleAdmin          RoleName = "admin"
	RoleManager        RoleName = "manager"
	RoleCashier        RoleName = "cashier"
	RoleInventoryStaff RoleName = "inventory_staff"
)

type PermissionAction string

const (
	ActionCreate PermissionAction = "create"
	ActionRead   PermissionAction = "read"
	ActionUpdate PermissionAction = "update"
	ActionDelete PermissionAction = "delete"
	ActionList   PermissionAction = "list"
)
