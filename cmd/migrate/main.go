package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"path/filepath"
	"siato_lms/internal/config"
	"sort"
	"strings"
)

func main() {
	c, e := config.Load()
	if e != nil {
		log.Fatal(e)
	}
	db, e := sql.Open("mysql", c.DB.DSN())
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	if _, e = db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version VARCHAR(100) PRIMARY KEY, applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=InnoDB`); e != nil {
		log.Fatal(e)
	}
	files, _ := filepath.Glob("migrations/*.up.sql")
	sort.Strings(files)
	for _, f := range files {
		v := strings.TrimSuffix(filepath.Base(f), ".up.sql")
		var n int
		_ = db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version=?", v).Scan(&n)
		if n > 0 {
			continue
		}
		b, e := os.ReadFile(f)
		if e != nil {
			log.Fatal(e)
		}
		tx, e := db.Begin()
		if e != nil {
			log.Fatal(e)
		}
		if _, e = tx.Exec(string(b)); e != nil {
			_ = tx.Rollback()
			log.Fatalf("migration %s: %v", v, e)
		}
		if _, e = tx.Exec("INSERT INTO schema_migrations(version) VALUES(?)", v); e != nil {
			_ = tx.Rollback()
			log.Fatal(e)
		}
		if e = tx.Commit(); e != nil {
			log.Fatal(e)
		}
		fmt.Println("applied", v)
	}
}
