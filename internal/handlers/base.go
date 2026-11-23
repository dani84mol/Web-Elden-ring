package handlers

import (
	"database/sql"
	"html/template"
	"sync"
)

// Application guarda las dependencias a nivel de aplicación
type Application struct {
	DB        *sql.DB
	Templates *template.Template

	// 🔥 CORRECCIÓN: Almacén de sesiones seguras en memoria (Evita Session Hijacking)
	Sessions map[string]string
	mu       sync.RWMutex
}

// New crea una nueva instancia de la aplicación con sus dependencias
func New(db *sql.DB, tmpl *template.Template) *Application {
	return &Application{
		DB:        db,
		Templates: tmpl,
		Sessions:  make(map[string]string),
	}
}
