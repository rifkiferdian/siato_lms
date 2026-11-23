package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

// SessionName must be unique among SIATO applications running on localhost.
// Browser cookies are scoped by domain and path, not by port, so using the
// generic "siato_session" would collide with other local SIATO services.
const SessionName = "siato_lms_session"

type SessionManager struct {
	store  *sessions.CookieStore
	db     *gorm.DB
	maxAge int
}

func NewSessionManager(secret string, secure bool, maxAge int, db *gorm.DB) *SessionManager {
	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{Path: "/", MaxAge: maxAge, HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode}
	return &SessionManager{store: store, db: db, maxAge: maxAge}
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
	token, _ := s.Values["session_token"].(string)
	if !ok || v == 0 || token == "" || m.db == nil {
		return 0, false
	}
	var count int64
	m.db.Raw("SELECT COUNT(*) FROM user_sessions WHERE user_id=? AND token_hash=? AND revoked_at IS NULL AND expires_at>UTC_TIMESTAMP()", v, hashToken(token)).Scan(&count)
	if count == 0 {
		return 0, false
	}
	return v, true
}
func hashToken(v string) string { sum := sha256.Sum256([]byte(v)); return hex.EncodeToString(sum[:]) }
func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, e := rand.Read(b); e != nil {
		return "", e
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func (m *SessionManager) Start(c *gin.Context, userID uint64) error {
	token, e := randomToken()
	if e != nil {
		return e
	}
	expires := time.Now().UTC().Add(time.Duration(m.maxAge) * time.Second)
	if e = m.db.Exec("INSERT INTO user_sessions(user_id,token_hash,ip_address,user_agent,expires_at) VALUES(?,?,?,?,?)", userID, hashToken(token), c.ClientIP(), c.Request.UserAgent(), expires).Error; e != nil {
		return e
	}
	s, _ := m.Get(c)
	s.Values["user_id"] = userID
	s.Values["session_token"] = token
	return s.Save(c.Request, c.Writer)
}
func (m *SessionManager) CurrentTokenHash(c *gin.Context) string {
	s, e := m.Get(c)
	if e != nil {
		return ""
	}
	token, _ := s.Values["session_token"].(string)
	if token == "" {
		return ""
	}
	return hashToken(token)
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
	if token, ok := s.Values["session_token"].(string); ok && token != "" && m.db != nil {
		_ = m.db.Exec("UPDATE user_sessions SET revoked_at=UTC_TIMESTAMP() WHERE token_hash=?", hashToken(token)).Error
	}
	s.Options.MaxAge = -1
	return s.Save(c.Request, c.Writer)
}
func (m *SessionManager) DestroyAll(c *gin.Context, userID uint64) error {
	if e := m.db.Exec("UPDATE user_sessions SET revoked_at=UTC_TIMESTAMP() WHERE user_id=? AND revoked_at IS NULL", userID).Error; e != nil {
		return e
	}
	return m.Destroy(c)
}
