package db

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (a *InMemoryAdapter) SeedInitialData(ctx context.Context) error {
	log.Println("🔄 Starting initial auth data seeding...")

	permissions := []struct {
		resource string
		action   string
		name     string
	}{
		// Order permissions
		{"order", "create", "Create Order"},
		{"order", "read", "Read Order"},
		{"order", "update", "Update Order"},

		// Product permissions
		{"product", "create", "Create Product"},
		{"product", "read", "Read Product"},
		{"product", "update", "Update Product"},
		{"product", "delete", "Delete Product"},

		// User permissions
		{"user", "create", "Create User"},
		{"user", "read", "Read User"},
		{"user", "update", "Update User"},
	}
	permissionIDs := make(map[string]string)
	for _, p := range permissions {
		perm, err := a.CreatePermission(ctx, p.resource, p.action, p.name)
		if err != nil {
			log.Printf("Warning: Could not create permission %s:%s: %v", p.resource, p.action, err)
			continue
		}
		permissionIDs[p.resource+"_"+p.action] = perm.ID
	}
	// Create roles
	roles := []struct {
		name        string
		description string
		permissions []string
	}{
		{
			name:        "admin",
			description: "Full system access",
			permissions: []string{
				permissionIDs["order_create"], permissionIDs["order_read"], permissionIDs["order_update"],
				permissionIDs["product_create"], permissionIDs["product_read"], permissionIDs["product_update"], permissionIDs["product_delete"],
				permissionIDs["user_create"], permissionIDs["user_read"], permissionIDs["user_update"],
			},
		},
		{
			name:        "manager",
			description: "Management access",
			permissions: []string{
				permissionIDs["order_read"], permissionIDs["order_update"],
				permissionIDs["product_create"], permissionIDs["product_read"], permissionIDs["product_update"],
				permissionIDs["user_read"],
			},
		},
		{
			name:        "cashier",
			description: "Order processing access",
			permissions: []string{
				permissionIDs["order_create"], permissionIDs["order_read"],
				permissionIDs["product_read"],
			},
		},
		{
			name:        "inventory_staff",
			description: "Inventory management access",
			permissions: []string{
				permissionIDs["product_create"], permissionIDs["product_read"], permissionIDs["product_update"],
			},
		},
	}
	for _, r := range roles {
		_, err := a.CreateRole(ctx, r.name, r.description, r.permissions)
		if err != nil {
			log.Printf("Warning: Could not create role %s: %v", r.name, err)
		}
	}
	// Create default admin user
	adminRole, err := a.GetRoleByName(ctx, "admin")
	if err != nil {
		log.Printf("Warning: Could not find admin role: %v", err)
		return nil
	}

	// Create default admin user (update password if exists)
	adminPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Warning: Could not hash admin password: %v", err)
		return nil
	}
	log.Printf("Seeding: Generated password hash: %s", string(adminPassword)[:20]+"...")

	existingUser, err := a.GetUserByUsername(ctx, "admin")
	if err == nil {
		// Update existing user's password
		existingUser.Password = string(adminPassword)
		log.Println("Seeding: Updated existing admin user password")
	} else {
		_, err = a.CreateUser(ctx, "admin@pos.com", "admin", string(adminPassword), []string{adminRole.ID})
		if err != nil {
			log.Printf("Warning: Could not create admin user: %v", err)
		} else {
			log.Println("DB Adapter: Created admin user")
		}
	}

	log.Println("✅ Initial data seeding completed successfully")
	return nil
}
