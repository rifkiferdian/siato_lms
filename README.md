# SIATO LMS

LMS perguruan tinggi berbasis Go, Gin, GORM, MySQL, dan server-side rendered templates. Repository berisi fondasi Phase 1, modul jadwal kuliah, dan modul kelas dengan filter akses berdasarkan kepesertaan serta pengampu.

## Menjalankan lokal

Prasyarat: Go 1.24+, Node.js 20+, MySQL 8 atau MariaDB 10.4+, dan GNU Make (opsional).

1. Salin `.env.example` menjadi `.env`, lalu ganti `SESSION_SECRET` dengan nilai acak minimal 32 karakter dan sesuaikan database.
2. Buat database: `CREATE DATABASE siato_lms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`
3. Jalankan `go mod download`.
4. Jalankan `npm install` lalu `npm run build`. Di PowerShell dengan execution policy ketat gunakan `npm.cmd`.
5. Jalankan `go run ./cmd/migrate`, lalu `go run ./cmd/seed`.
6. Jalankan `go run ./cmd/server`, lalu buka `http://localhost:8080`.

Untuk live reload dengan Air, jalankan `air`. Konfigurasi `.air.toml` akan membangun entrypoint `./cmd/server` menjadi `tmp/main.exe` pada Windows.

Akun development: `admin@siato.test`, `akademik@siato.test`, `dosen@siato.test`, atau `mahasiswa@siato.test`; kata sandi seluruhnya `Password123!`. Wajib diganti atau dihapus sebelum production.

## Pemeriksaan

`go test ./...`, `go vet ./...`, dan `go build ./cmd/server`.

## Catatan keamanan

Set `COOKIE_SECURE=true` di balik HTTPS. Private storage tidak dipasang sebagai static route. Phase berikutnya akan menambahkan object-level authorization untuk entitas akademik, upload private tervalidasi, audit event bisnis, serta Origin validation untuk WebSocket.

## Status roadmap

- Selesai: fondasi Phase 1, login/logout, session, CSRF, RBAC dasar, migration/seed, dashboard, modul jadwal kuliah/mengajar, daftar dan detail kelas.
- Belum: CRUD master akademik (Phase 2), detail kelas/pertemuan, QR, materi/tugas, messaging, notification, laporan. Menu untuk fitur tersebut belum berupa route aktif sampai implementasinya tersedia.
