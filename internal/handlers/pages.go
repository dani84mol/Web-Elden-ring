package handlers

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strings"

	"elden-guide/internal/models"
)

func (app *Application) getUsuario(r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.Sessions[cookie.Value]
}

// 🔥 NUEVA FUNCIÓN: Uso de Búfer para HTML (Exigido en Pág. 72)
func (app *Application) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	var buf bytes.Buffer
	err := app.Templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		log.Printf("ERROR: Fallo al renderizar la plantilla '%s': %v", name, err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}

func (app *Application) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		path := "./public" + r.URL.Path
		if _, err := os.Stat(path); os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			data := struct{ Usuario string }{Usuario: app.getUsuario(r)}
			app.renderTemplate(w, "404.html", data)
			return
		}
		http.ServeFile(w, r, path)
		return
	}
	data := struct{ Usuario string }{Usuario: app.getUsuario(r)}
	app.renderTemplate(w, "index.html", data)
}

func (app *Application) Clases(w http.ResponseWriter, r *http.Request) {
	usuario := app.getUsuario(r)
	favoritos := make(map[string]bool)

	// Obtener favoritos desde la base de datos si el usuario está logueado
	if usuario != "" {
		var err error
		// Se asume que tu función en models se llama GetFavoritos
		favoritos, err = models.GetFavoritos(app.DB, usuario)
		if err != nil {
			log.Printf("ERROR: Fallo al obtener favoritos del usuario '%s': %v", usuario, err)
		}
	}

	data := ClasesPageData{
		Usuario:   usuario,
		Favoritos: favoritos,
	}

	app.renderTemplate(w, "Clases.html", data)
}
func (app *Application) Jefes(w http.ResponseWriter, r *http.Request) {
	nombreUsuario := app.getUsuario(r)
	derrotadosMap := make(map[string]bool)

	if nombreUsuario != "" {
		var err error
		derrotadosMap, err = models.GetDerrotados(app.DB, nombreUsuario)
		if err != nil {
			log.Printf("ERROR: Fallo al obtener jefes derrotados del usuario '%s': %v", nombreUsuario, err)
		}
	}

	data := struct {
		Usuario    string
		Derrotados map[string]bool
	}{Usuario: nombreUsuario, Derrotados: derrotadosMap}
	app.renderTemplate(w, "Jefes.html", data)
}

func (app *Application) Zonas(w http.ResponseWriter, r *http.Request) {
	nombreUsuario := app.getUsuario(r)
	obtenidosMap := make(map[string]bool)

	if nombreUsuario != "" {
		var err error
		obtenidosMap, err = models.GetObtenidos(app.DB, nombreUsuario)
		if err != nil {
			log.Printf("ERROR: Fallo al obtener objetos de zona del usuario '%s': %v", nombreUsuario, err)
		}
	}

	data := struct {
		Usuario   string
		Obtenidos map[string]bool
	}{Usuario: nombreUsuario, Obtenidos: obtenidosMap}
	app.renderTemplate(w, "Zonas.html", data)
}

func (app *Application) Consejos(w http.ResponseWriter, r *http.Request) {
	data := struct{ Usuario string }{Usuario: app.getUsuario(r)}
	app.renderTemplate(w, "Consejos.html", data)
}

func (app *Application) render500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	data := struct{ Usuario string }{Usuario: app.getUsuario(r)}
	app.renderTemplate(w, "500.html", data)
}

func (app *Application) LoginPage(w http.ResponseWriter, r *http.Request) {
	referer := r.Referer()
	if referer != "" && !strings.Contains(referer, "/login.html") && !strings.Contains(referer, "/registro.html") {
		http.SetCookie(w, &http.Cookie{Name: "redirect_after_login", Value: referer, Path: "/"})
	}

	errorParam := r.URL.Query().Get("error")
	var mensaje string
	if errorParam == "usuario" || errorParam == "password" {
		mensaje = "Las credenciales no coinciden con ningún Sinluz registrado."
	}
	app.renderTemplate(w, "login.html", struct{ Error string }{Error: mensaje})
}

func (app *Application) RegistroPage(w http.ResponseWriter, r *http.Request) {
	referer := r.Referer()
	if referer != "" && !strings.Contains(referer, "/login.html") && !strings.Contains(referer, "/registro.html") {
		http.SetCookie(w, &http.Cookie{Name: "redirect_after_login", Value: referer, Path: "/"})
	}

	errorParam := r.URL.Query().Get("error")
	var mensaje string
	switch errorParam {
	case "usuario":
		mensaje = "El nombre de usuario debe ser de 3 a 20 caracteres."
	case "correo":
		mensaje = "El correo electrónico no es válido."
	case "password":
		mensaje = "La contraseña es demasiado débil."
	case "confirm":
		mensaje = "Las contraseñas no coinciden."
	case "existe":
		mensaje = "Este nombre o correo ya pertenece a otra alma."
	}
	app.renderTemplate(w, "registro.html", struct{ Error string }{Error: mensaje})
}
