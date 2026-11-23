package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var validarUsuario = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
var validarCorreo = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func generateSessionID() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (app *Application) ProcesarRegistro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/registro.html", http.StatusSeeOther)
		return
	}

	r.ParseForm()
	usuario := strings.TrimSpace(r.FormValue("usuario"))
	correo := strings.TrimSpace(r.FormValue("correo"))
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	// Validaciones básicas
	if !validarUsuario.MatchString(usuario) || !validarCorreo.MatchString(correo) || len(password) < 8 || password != confirm {
		log.Printf("WARN: Intento de registro denegado por validación para el usuario '%s'", usuario)
		http.Redirect(w, r, "/registro.html?error=usuario", http.StatusSeeOther)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := app.DB.Exec("INSERT INTO usuarios (usuario, correo, password_hash) VALUES (?, ?, ?)", usuario, correo, string(hash))

	if err != nil {
		log.Printf("WARN: Colisión en registro (Usuario o correo ya existe) para '%s' o '%s': %v", usuario, correo, err)
		http.Redirect(w, r, "/registro.html?error=existe", http.StatusSeeOther)
		return
	}

	log.Printf("INFO: Nuevo usuario registrado exitosamente: '%s'", usuario)
	http.Redirect(w, r, "/login.html", http.StatusSeeOther)
}

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	usuario := strings.TrimSpace(r.FormValue("usuario"))
	password := r.FormValue("password")

	var hash string
	err := app.DB.QueryRow("SELECT password_hash FROM usuarios WHERE usuario = ?", usuario).Scan(&hash)

	if err != nil {
		log.Printf("WARN: Intento de login con usuario inexistente: '%s'", usuario)
		http.Redirect(w, r, "/login.html?error=usuario", http.StatusSeeOther)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("WARN: Contraseña incorrecta para el usuario: '%s'", usuario)
		http.Redirect(w, r, "/login.html?error=password", http.StatusSeeOther)
		return
	}

	// Crear sesión segura
	sessionID := generateSessionID()
	app.mu.Lock()
	app.Sessions[sessionID] = usuario
	app.mu.Unlock()

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,                 // Pág 83: Protección XSS
		Secure:   true,                 // Pág 83: Solo HTTPS
		SameSite: http.SameSiteLaxMode, // Pág 83: Protección CSRF
	}
	http.SetCookie(w, cookie)

	log.Printf("INFO: Usuario '%s' ha iniciado sesión correctamente", usuario)

	dest := "/"
	if c, err := r.Cookie("redirect_after_login"); err == nil {
		dest = c.Value
		http.SetCookie(w, &http.Cookie{Name: "redirect_after_login", Value: "", Path: "/", MaxAge: -1})
	}
	http.Redirect(w, r, dest, http.StatusSeeOther)
}

func (app *Application) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session_id"); err == nil {
		app.mu.Lock()
		usuario := app.Sessions[cookie.Value]
		delete(app.Sessions, cookie.Value)
		app.mu.Unlock()
		log.Printf("INFO: Usuario '%s' ha cerrado sesión", usuario)
	}

	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: "", Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
