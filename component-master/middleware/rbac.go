package middleware

import (
	"component-master/config"
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gofiber/fiber/v2"
)

// CasbinMiddleware represents the Casbin auth middleware
type CasbinMiddleware struct {
	enforcer *casbin.Enforcer
}

// NewCasbinMiddleware creates a new Casbin middleware
func NewCasbinMiddleware(conf *config.CasbinConfig) (*CasbinMiddleware, error) {
	// Load model from file
	m, err := model.NewModelFromFile(conf.ModelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	// Load policy from file
	adapter := fileadapter.NewAdapter(conf.PolicyPath)

	// Create enforcer
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	// Load policies
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("failed to load policy: %w", err)
	}

	return &CasbinMiddleware{
		enforcer: enforcer,
	}, nil
}

// Authorize middleware checks if the user has permission
func (cm *CasbinMiddleware) Authorize() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user from context (set by previous auth middleware)
		user := c.Locals("user").(Claims)
		role := user.Role

		// Get request path and method
		path := c.Path()
		method := c.Method()

		// Check permission
		allowed, err := cm.enforcer.Enforce(role, path, method)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check permission",
			})
		}

		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Permission denied",
			})
		}

		return c.Next()
	}
}

// AddPolicy adds a new policy rule
func (cm *CasbinMiddleware) AddPolicy(role, path, method string) bool {
	success, _ := cm.enforcer.AddPolicy(role, path, method)
	return success
}

// RemovePolicy removes a policy rule
func (cm *CasbinMiddleware) RemovePolicy(role, path, method string) bool {
	success, _ := cm.enforcer.RemovePolicy(role, path, method)
	return success
}

// HasPolicy checks if a policy exists
func (cm *CasbinMiddleware) HasPolicy(role, path, method string) bool {
	hasPolicy, _ := cm.enforcer.HasPolicy(role, path, method)
	return hasPolicy
}

// GetRoles returns all roles assigned to a user
func (cm *CasbinMiddleware) GetRoles(username string) []string {
	roles, _ := cm.enforcer.GetRolesForUser(username)
	return roles
}

// AddRoleForUser assigns a role to a user
func (cm *CasbinMiddleware) AddRoleForUser(username, role string) bool {
	success, _ := cm.enforcer.AddRoleForUser(username, role)
	return success
}
