package db

import (
	"context"
	"fmt"
	"log"
	"pos/internal/modules/auth"
	"pos/internal/modules/inventory"
	"pos/internal/modules/orders"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type InMemoryAdapter struct {
	products    map[string]*inventory.Product
	orders      map[string]*orders.Order
	users       map[string]*auth.User
	roles       map[string]*auth.Role
	permissions map[string]*auth.Permission
	mu          sync.RWMutex
}

func NewInMemoryAdapter() *InMemoryAdapter {
	return &InMemoryAdapter{
		products:    make(map[string]*inventory.Product),
		orders:      make(map[string]*orders.Order),
		users:       make(map[string]*auth.User),
		roles:       make(map[string]*auth.Role),
		permissions: make(map[string]*auth.Permission),
	}
}

func (a *InMemoryAdapter) SaveProduct(product *inventory.Product) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.products[product.ID] = product
	log.Printf("DB Adapter: Saved product %s", product.Name)
	return nil
}

func (a *InMemoryAdapter) GetProduct(id string) (*inventory.Product, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	product, ok := a.products[id]
	if !ok {
		return nil, fmt.Errorf("product %s not found in DB", id)
	}
	return product, nil
}

func (a *InMemoryAdapter) GetAllProducts() ([]*inventory.Product, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	products := make([]*inventory.Product, 0, len(a.products))
	for _, product := range a.products {
		products = append(products, product)
	}
	return products, nil
}

func (a *InMemoryAdapter) GetProductByName(name string) (*inventory.Product, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, product := range a.products {
		if product.Name == name {
			return product, nil
		}
	}
	return nil, fmt.Errorf("product with name %s not found in DB", name)
}

func (a *InMemoryAdapter) UpdateProductStock(id string, newStock int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	product, ok := a.products[id]
	if !ok {
		return fmt.Errorf("product %s not found in DB", id)
	}
	product.Stock = newStock
	log.Printf("DB Adapter: Updated stock for %s to %d", product.Name, newStock)
	return nil
}

func (a *InMemoryAdapter) SaveOrder(order *orders.Order) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.orders[order.ID] = order
	log.Printf("DB Adapter: Saved order %s (Status: %s)", order.ID, order.Status)
	return nil
}

func (a *InMemoryAdapter) GetOrderByID(id string) (*orders.Order, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	order, ok := a.orders[id]
	if !ok {
		return nil, fmt.Errorf("order %s not found in DB", id)
	}
	return order, nil
}

func (a *InMemoryAdapter) CreateUser(ctx context.Context, email, username, passwordHash string, roleIDs []string) (*auth.User, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, user := range a.users {
		if user.Email == email {
			return nil, fmt.Errorf("user with email %s already exists", email)
		}
		if user.Username == username {
			return nil, fmt.Errorf("user with username %s already exists", username)
		}
	}
	user := &auth.User{
		ID:        uuid.New().String(),
		Email:     email,
		Username:  username,
		Password:  passwordHash,
		RoleIDs:   roleIDs,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	a.users[user.ID] = user
	log.Printf("DB Adapter: Created user %s (%s) with password hash starting: %s", user.Username, user.Email, passwordHash[:20]+"...")
	return user, nil
}

func (a *InMemoryAdapter) GetUserByID(ctx context.Context, id string) (*auth.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	user, ok := a.users[id]
	if !ok {
		return nil, fmt.Errorf("user %s not found", id)
	}
	return user, nil
}

func (a *InMemoryAdapter) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	for _, user := range a.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user with email %s not found", email)
}

func (a *InMemoryAdapter) GetUserByUsername(ctx context.Context, username string) (*auth.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	log.Printf("DB: Looking for user with username: %s", username)
	log.Printf("DB: Total users in DB: %d", len(a.users))
	for _, user := range a.users {
		log.Printf("DB: Checking user: %s", user.Username)
		if user.Username == username {
			log.Printf("DB: Found user %s", user.Username)
			return user, nil
		}
	}
	log.Printf("DB: User %s not found", username)
	return nil, fmt.Errorf("user with username %s not found", username)
}

func (a *InMemoryAdapter) UpdateUser(ctx context.Context, user *auth.User) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.users[user.ID]; !ok {
		return fmt.Errorf("user %s not found", user.ID)
	}
	user.UpdatedAt = time.Now()
	a.users[user.ID] = user
	log.Printf("DB Adapter: Updated user %s", user.Username)
	return nil
}

func (a *InMemoryAdapter) DeleteUser(ctx context.Context, id string) error {
	a.mu.Lock()
	if _, ok := a.users[id]; !ok {
		return fmt.Errorf("user %s not found", id)
	}
	delete(a.users, id)
	log.Printf("DB Adapter: Deleted user %s", id)
	return nil
}

func (a *InMemoryAdapter) ListUsers(ctx context.Context, limit, offset int) ([]*auth.User, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	users := make([]*auth.User, 0, len(a.users))
	for _, user := range a.users {
		users = append(users, user)
	}
	start := offset
	if start > len(users) {
		start = len(users)
	}
	end := start + limit
	if end > len(users) {
		end = len(users)
	}
	return users[start:end], nil
}

func (a *InMemoryAdapter) CreateRole(ctx context.Context, name, description string, permissionIDs []string) (*auth.Role, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if role name already exists
	for _, role := range a.roles {
		if strings.EqualFold(role.Name, name) {
			return nil, fmt.Errorf("role with name %s already exists", name)
		}
	}

	role := &auth.Role{
		ID:            uuid.New().String(),
		Name:          name,
		Description:   description,
		PermissionIDs: permissionIDs,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	a.roles[role.ID] = role
	log.Printf("DB Adapter: Created role %s", role.Name)
	return role, nil
}

func (a *InMemoryAdapter) GetRoleByID(ctx context.Context, id string) (*auth.Role, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	role, ok := a.roles[id]
	if !ok {
		return nil, fmt.Errorf("role %s not found", id)
	}
	return role, nil
}

func (a *InMemoryAdapter) GetRoleByName(ctx context.Context, name string) (*auth.Role, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, role := range a.roles {
		if strings.EqualFold(role.Name, name) {
			return role, nil
		}
	}
	return nil, fmt.Errorf("role with name %s not found", name)
}

func (a *InMemoryAdapter) UpdateRole(ctx context.Context, role *auth.Role) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.roles[role.ID]; !ok {
		return fmt.Errorf("role %s not found", role.ID)
	}

	role.UpdatedAt = time.Now()
	a.roles[role.ID] = role
	log.Printf("DB Adapter: Updated role %s", role.Name)
	return nil
}

func (a *InMemoryAdapter) DeleteRole(ctx context.Context, id string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.roles[id]; !ok {
		return fmt.Errorf("role %s not found", id)
	}

	delete(a.roles, id)
	log.Printf("DB Adapter: Deleted role %s", id)
	return nil
}

func (a *InMemoryAdapter) ListRoles(ctx context.Context, limit, offset int) ([]*auth.Role, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	roles := make([]*auth.Role, 0, len(a.roles))
	for _, role := range a.roles {
		roles = append(roles, role)
	}

	// Apply offset and limit
	start := offset
	if start > len(roles) {
		start = len(roles)
	}
	end := start + limit
	if end > len(roles) {
		end = len(roles)
	}

	return roles[start:end], nil
}

func (a *InMemoryAdapter) CreatePermission(ctx context.Context, resource, action, name string) (*auth.Permission, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, perm := range a.permissions {
		if perm.Resource == resource && perm.Action == action {
			return nil, fmt.Errorf("permission for %s:%s already exists", resource, action)
		}
	}
	permission := &auth.Permission{
		ID:       uuid.New().String(),
		Resource: resource,
		Action:   action,
		Name:     name,
	}
	a.permissions[permission.ID] = permission
	log.Printf("DB Adapter: Created permission %s", permission.Name)
	return permission, nil
}

func (a *InMemoryAdapter) GetPermissionByID(ctx context.Context, id string) (*auth.Permission, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	permission, ok := a.permissions[id]
	if !ok {
		return nil, fmt.Errorf("permission %s not found", id)
	}
	return permission, nil
}

func (a *InMemoryAdapter) GetPermissionsByResource(ctx context.Context, resource string) ([]*auth.Permission, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var permissions []*auth.Permission
	for _, perm := range a.permissions {
		if perm.Resource == resource {
			permissions = append(permissions, perm)
		}
	}
	return permissions, nil
}

func (a *InMemoryAdapter) ListPermissions(ctx context.Context, limit, offset int) ([]*auth.Permission, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	permissions := make([]*auth.Permission, 0, len(a.permissions))
	for _, perm := range a.permissions {
		permissions = append(permissions, perm)
	}

	// Apply offset and limit
	start := offset
	if start > len(permissions) {
		start = len(permissions)
	}
	end := start + limit
	if end > len(permissions) {
		end = len(permissions)
	}

	return permissions[start:end], nil
}
