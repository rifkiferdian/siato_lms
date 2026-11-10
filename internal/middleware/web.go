package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(self), microphone=(), geolocation=()")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; font-src 'self'; img-src 'self' data:; connect-src 'self' ws: wss:")
		c.Next()
	}
}
func CSRF(sm *SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		s, e := sm.Get(c)
		expected, _ := s.Values["csrf"].(string)
		actual := c.GetHeader("X-CSRF-Token")
		if actual == "" {
			actual = c.PostForm("_csrf")
		}
		if e != nil || expected == "" || len(expected) != len(actual) || subtle.ConstantTimeCompare([]byte(expected), []byte(actual)) != 1 {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
func RequireAuth(sm *SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := sm.UserID(c)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
		c.Set("user_id", id)
		c.Next()
	}
}
func RequirePermission(check func(uint64, string) bool, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, ok := c.Get("user_id")
		if !ok || !check(id.(uint64), permission) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
func LoginLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(0.2), 5)
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost && !limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) { c.Header("Cache-Control", "no-store"); c.Next() }
}
func WantsJSON(c *gin.Context) bool {
	return strings.Contains(c.GetHeader("Accept"), "application/json")
}
