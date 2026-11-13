package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"siato_lms/internal/config"
	"siato_lms/internal/database"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/auth"
	"siato_lms/internal/modules/classes"
	"siato_lms/internal/modules/dashboard"
	"siato_lms/internal/modules/people"
	"siato_lms/internal/modules/rbac"
	"siato_lms/internal/modules/schedules"
	"siato_lms/internal/modules/usermanagement"
	"siato_lms/internal/shared/view"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		log.Error("configuration error", "error", err)
		os.Exit(1)
	}
	db, err := database.Open(cfg, log)
	if err != nil {
		log.Error("database error", "error", err)
		os.Exit(1)
	}
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery(), requestLogger(log), middleware.SecurityHeaders())
	if err := r.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		log.Error("trusted proxies", "error", err)
		os.Exit(1)
	}
	r.Static("/static", "./web/static")
	sm := middleware.NewSessionManager(cfg.SessionSecret, cfg.CookieSecure, cfg.SessionMaxAge)
	r.Use(middleware.CSRF(sm))
	views := view.New(sm, "./web/templates", cfg.AppName)
	rbacService := rbac.NewService(db)
	authHandler := auth.NewHandler(auth.NewService(auth.NewRepository(db)), rbacService, sm, views)
	auth.RegisterRoutes(r, authHandler, sm)
	dashboard.RegisterRoutes(r, dashboard.NewHandler(rbacService, sm, views, db), sm, rbacService)
	scheduleHandler := schedules.NewHandler(schedules.NewService(schedules.NewRepository(db)), rbacService, sm, views, db)
	schedules.RegisterRoutes(r, scheduleHandler, sm, rbacService)
	classHandler := classes.NewHandler(classes.NewService(classes.NewRepository(db)), rbacService, sm, views, db)
	classes.RegisterRoutes(r, classHandler, sm, rbacService)
	peopleHandler := people.NewHandler(people.NewService(people.NewRepository(db)), rbacService, sm, views, db)
	people.RegisterRoutes(r, peopleHandler, sm, rbacService)
	userManagementHandler := usermanagement.NewHandler(usermanagement.NewService(usermanagement.NewRepository(db)), rbacService, sm, views, db)
	usermanagement.RegisterRoutes(r, userManagementHandler, sm, rbacService)
	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusSeeOther, "/dashboard") })
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	server := &http.Server{Addr: cfg.Address, Handler: r, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	log.Info("server starting", "address", cfg.Address, "environment", cfg.Environment)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func requestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info("http request", "method", c.Request.Method, "path", c.Request.URL.Path, "status", c.Writer.Status(), "duration_ms", time.Since(start).Milliseconds(), "client_ip", c.ClientIP())
	}
}
