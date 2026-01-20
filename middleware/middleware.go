package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/Be2Bag/erp-demo/config"
	"github.com/Be2Bag/erp-demo/dto"
	"github.com/Be2Bag/erp-demo/models"
	"github.com/Be2Bag/erp-demo/pkg/util"
	"github.com/Be2Bag/erp-demo/ports"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Middleware struct {
	JWT config.JWTConfig
}

func NewMiddleware(jwtCfg config.JWTConfig) *Middleware {
	return &Middleware{JWT: jwtCfg}
}

// Replace previous standalone AuthCookieMiddleware with method version.
func (m *Middleware) AuthCookieMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("auth_token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusUnauthorized,
				MessageEN:  "No auth token cookie found",
				MessageTH:  "ไม่พบคุกกี้ auth token",
				Status:     "error",
				Data:       nil,
			})
		}
		if m.JWT.SecretKey == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusInternalServerError,
				MessageEN:  "JWT secret not configured",
				MessageTH:  "ไม่ได้ตั้งค่า JWT secret",
				Status:     "error",
				Data:       nil,
			})
		}

		claims, err := util.VerifyAndParseJWTClaims(token, m.JWT.SecretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusUnauthorized,
				MessageEN:  "Invalid or expired token",
				MessageTH:  "โทเคนไม่ถูกต้องหรือหมดอายุ",
				Status:     "error",
				Data:       nil,
			})
		}

		c.Locals("auth_claims", claims)
		return c.Next()
	}
}

func GetClaims(c *fiber.Ctx) (*dto.JWTClaims, error) {
	claims, ok := c.Locals("auth_claims").(*dto.JWTClaims)
	if !ok || claims == nil || claims.UserID == "" {
		return nil, fiber.ErrUnauthorized
	}
	return claims, nil
}

func (m *Middleware) TimeoutMiddleware(d time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// BUG (เดิม): ใช้ c.Context() ซึ่งเป็น *fasthttp.RequestCtx ไม่ใช่ context.Context -> ทำให้ compile ผิด/ทำงานไม่ถูกต้อง
		parentCtx := c.UserContext()
		ctx, cancel := context.WithTimeout(parentCtx, d)
		defer cancel()

		c.SetUserContext(ctx)

		err := c.Next()

		// หาก handler ใช้ ctx จาก c.UserContext() จะโดน timeout ได้ถูกต้อง
		if ctx.Err() == context.DeadlineExceeded {
			return c.Status(fiber.StatusRequestTimeout).JSON(dto.BaseResponse{
				StatusCode: fiber.StatusRequestTimeout,
				MessageEN:  "Request timed out",
				MessageTH:  "คำขอใช้เวลานานเกินไป",
				Status:     "error",
				Data:       nil,
			})
		}
		return err
	}
}

// AuditLogMiddleware logs all authenticated requests (POST, PUT, DELETE) to the audit log
func AuditLogMiddleware(auditSvc ports.AuditLogService, jwtSecret string) fiber.Handler {
	// Paths to exclude from logging
	excludePaths := []string{
		"/swagger",
		"/health",
		"/v1/audit-log", // Don't log audit log queries themselves
	}

	return func(c *fiber.Ctx) error {
		method := c.Method()
		path := c.Path()

		// Only log POST, PUT, DELETE (skip GET for better performance)
		if method == "GET" || method == "OPTIONS" || method == "HEAD" {
			return c.Next()
		}

		// Skip excluded paths
		for _, excludePath := range excludePaths {
			if strings.HasPrefix(path, excludePath) {
				return c.Next()
			}
		}

		// Record start time
		startTime := time.Now()

		// Execute the handler
		err := c.Next()

		// Calculate response time
		responseTime := time.Since(startTime).Milliseconds()

		// Get status code
		statusCode := c.Response().StatusCode()

		// Try to get user claims (may not exist for unauthenticated requests)
		var userID, email, employeeCode, role, fullName string
		if claims, claimErr := GetClaims(c); claimErr == nil && claims != nil {
			userID = claims.UserID
			email = claims.Email
			employeeCode = claims.EmployeeCode
			role = claims.Role
			fullName = claims.TitleTH + claims.FirstNameTH + " " + claims.LastNameTH
		}

		// Skip if no user (unauthenticated request)
		if userID == "" {
			return err
		}

		// Extract resource and action from path
		resource, action, resourceID := extractResourceInfo(method, path, c.Params("id"))

		// Get request body (limited size for security)
		requestBody := ""
		if len(c.Body()) > 0 && len(c.Body()) < 10000 {
			// Sanitize sensitive fields
			requestBody = sanitizeRequestBody(string(c.Body()))
		}

		// Create audit log entry
		auditLog := &models.AuditLog{
			LogID:        uuid.New().String(),
			UserID:       userID,
			Email:        email,
			EmployeeCode: employeeCode,
			Role:         role,
			FullName:     fullName,
			Method:       method,
			Path:         path,
			QueryParams:  string(c.Request().URI().QueryString()),
			RequestBody:  requestBody,
			IPAddress:    c.IP(),
			UserAgent:    c.Get("User-Agent"),
			StatusCode:   statusCode,
			ResponseTime: responseTime,
			Action:       action,
			Resource:     resource,
			ResourceID:   resourceID,
			CreatedAt:    time.Now(),
		}

		// Log asynchronously (non-blocking)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = auditSvc.Log(ctx, auditLog)
		}()

		return err
	}
}

// extractResourceInfo extracts resource type, action, and ID from the request
func extractResourceInfo(method, path, paramID string) (resource, action, resourceID string) {
	// Extract resource from path (e.g., /service/api/v1/user/create -> user)
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "v1" && i+1 < len(parts) {
			resource = parts[i+1]
			break
		}
	}

	// Determine action based on HTTP method
	switch method {
	case "POST":
		action = "CREATE"
	case "PUT", "PATCH":
		action = "UPDATE"
	case "DELETE":
		action = "DELETE"
	default:
		action = "READ"
	}

	// Check for specific action keywords in path
	pathLower := strings.ToLower(path)
	if strings.Contains(pathLower, "/create") {
		action = "CREATE"
	} else if strings.Contains(pathLower, "/update") {
		action = "UPDATE"
	} else if strings.Contains(pathLower, "/delete") {
		action = "DELETE"
	} else if strings.Contains(pathLower, "/confirm") {
		action = "UPDATE"
	} else if strings.Contains(pathLower, "/copy") {
		action = "CREATE"
	}

	resourceID = paramID

	return resource, action, resourceID
}

// sanitizeRequestBody removes sensitive fields from request body
func sanitizeRequestBody(body string) string {
	// Simple sanitization - replace common sensitive field patterns
	sensitivePatterns := []string{
		`"password"`, `"Password"`,
		`"secret"`, `"Secret"`,
		`"token"`, `"Token"`,
		`"credit_card"`, `"creditCard"`,
	}

	result := body
	for _, pattern := range sensitivePatterns {
		if strings.Contains(result, pattern) {
			// Replace the value after the sensitive key with [REDACTED]
			// This is a simple approach; for production, use proper JSON parsing
			result = strings.ReplaceAll(result, pattern, pattern[:len(pattern)-1]+`":"[REDACTED]"`)
		}
	}

	return result
}
