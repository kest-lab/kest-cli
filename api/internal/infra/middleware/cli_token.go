package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kest-labs/kest/api/pkg/response"
)

type ProjectCLITokenValidator interface {
	ValidateCLIToken(ctx context.Context, projectID uint, rawToken string, requiredScopes []string) (uint, uint, error)
}

// RequireProjectCLIToken validates a project-scoped CLI token and enforces required scopes.
func RequireProjectCLIToken(validator ProjectCLITokenValidator, requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		projectIDStr := c.Param("id")
		if projectIDStr == "" {
			projectIDStr = c.Param("pid")
		}
		projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid Project ID")
			c.Abort()
			return
		}

		tokenID, createdBy, err := validator.ValidateCLIToken(c.Request.Context(), uint(projectID), parts[1], requiredScopes)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		c.Set("cliTokenID", tokenID)
		c.Set("cliTokenCreatedBy", createdBy)
		c.Next()
	}
}
