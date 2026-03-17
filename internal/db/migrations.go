package db

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed migration/*.sql
var migrationsFS embed.FS

func MigrationsFS() http.FileSystem {
	sub, _ := fs.Sub(migrationsFS, "migration")
	return http.FS(sub)
}
