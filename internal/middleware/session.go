package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// SessionName must be unique among SIATO applications running on localhost.
// Browser cookies are scoped by domain and path, not by port, so using the
// generic "siato_session" would collide with other local SIATO services.
const SessionName = "siato_lms_session"

type SessionManager struct{ store *sessions.CookieStore }

func NewSessionManager(secret string, secure bool, maxAge int) *SessionManager {
	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{Path: "/", MaxAge: maxAge, HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode}
	return &SessionManager{store: store}
}
func (m *SessionManager) Get(c *gin.Context) (*sessions.Session, error) {
	return m.store.Get(c.Request, SessionName)
}
func (m *SessionManager) UserID(c *gin.Context) (uint64, bool) {
	s, err := m.Get(c)
	if err != nil {
		return 0, false
	}
	v, ok := s.Values["user_id"].(uint64)
	return v, ok && v > 0
}
func (m *SessionManager) CSRF(c *gin.Context) (string, error) {
	s, err := m.Get(c)
	if err != nil {
		return "", err
	}
	if token, ok := s.Values["csrf"].(string); ok && token != "" {
		return token, nil
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(b)
	s.Values["csrf"] = token
	return token, s.Save(c.Request, c.Writer)
}
func (m *SessionManager) Destroy(c *gin.Context) error {
	s, e := m.Get(c)
	if e != nil {
		return e
	}
	s.Options.MaxAge = -1
	return s.Save(c.Request, c.Writer)
}
