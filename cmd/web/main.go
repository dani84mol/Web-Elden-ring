package main

import (
	"html/template"
	"log"
	"net/http"

	"elden-guide/internal/database"
	"elden-guide/internal/handlers"
)

func main() {

	// ========================
	// INICIALIZAR BASE DE DATOS
	// ========================
	db, err := database.InitDB("./auth.db")
	if err != nil {
		log.Fatal("Error al inicializar la base de datos:", err)
	}
	defer db.Close()

	// 🔥 CORRECCIÓN RENDIMIENTO: Compilar plantillas UNA SOLA VEZ
	tmpl, err := template.ParseGlob("./public/*.html")
	if err != nil {
		log.Fatal("Error compilando las plantillas HTML:", err)
	}

	// Crear app con dependencias
	app := handlers.New(db, tmpl)

	// Router
	mux := http.NewServeMux()

	// ========================
	// RUTAS AUTH
	// ========================
	mux.HandleFunc("/procesar-registro", app.ProcesarRegistro)
	mux.HandleFunc("/procesar-login", app.Login)
	mux.HandleFunc("/logout", app.Logout)

	mux.HandleFunc("/login.html", app.LoginPage)
	mux.HandleFunc("/registro.html", app.RegistroPage)

	// ========================
	// RUTAS PÁGINAS
	// ========================
	mux.HandleFunc("/", app.Index)
	mux.HandleFunc("/index.html", app.Index)
	mux.HandleFunc("/Clases.html", app.Clases)
	mux.HandleFunc("/Jefes.html", app.Jefes)
	mux.HandleFunc("/Zonas.html", app.Zonas)
	mux.HandleFunc("/Consejos.html", app.Consejos)

	// ========================
	// RUTAS DE ACCIÓN (Juego) - 🔥 ESTO FALTABA
	// ========================
	mux.HandleFunc("/toggle-favorito", app.ToggleFavorito)
	mux.HandleFunc("/toggle-objeto", app.ToggleObjeto)
	mux.HandleFunc("/marcar-muerto", app.MarcarMuerto)
	mux.HandleFunc("/revivir-jefe", app.RevivirJefe)
	mux.HandleFunc("/recomendar-clase", app.RecomendarClaseHandler)
	mux.HandleFunc("/api/notas", app.NotasAPI)
	// ========================
	// ARCHIVOS ESTÁTICOS
	// ========================
	fs := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
