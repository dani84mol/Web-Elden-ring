package handlers

import (
	"elden-guide/internal/models"
	"log"
	"net/http"
)

type ClasesPageData struct {
	Usuario          string
	Favoritos        map[string]bool
	ResultadoOraculo string // Almacena el nombre (Ej: "Vagante")
	ResultadoId      string // Almacena el ID (Ej: "vagante")
}

var ordenJefes = []string{
	"margit", "godrick", "lobo", "rennala", "centinela",
	"morgott", "gigante", "sacrodermo", "maliketh",
	"gideon", "godfrey", "radagon", "bestia",
}

func (app *Application) ToggleFavorito(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	claseID := r.FormValue("clase_id")
	usuario := app.getUsuario(r)

	if usuario == "" {
		redirectPath := "/Clases.html"
		if claseID != "" {
			redirectPath += "#" + claseID
		}
		http.SetCookie(w, &http.Cookie{Name: "redirect_after_login", Value: redirectPath, Path: "/"})
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	err := models.ToggleFavorito(app.DB, usuario, claseID)
	if err != nil {
		log.Printf("ERROR: Fallo al procesar favorito de clase '%s' para usuario '%s': %v", claseID, usuario, err)
	} else {
		log.Printf("INFO: Usuario '%s' ha alternado favorito de la clase '%s'", usuario, claseID)
	}
	http.Redirect(w, r, "/Clases.html#"+claseID, http.StatusSeeOther)
}

func (app *Application) ToggleObjeto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	objetoID := r.FormValue("objeto_id")
	zonaID := r.FormValue("zona_id")
	usuario := app.getUsuario(r)

	if usuario == "" {
		redirectPath := "/Zonas.html"
		if zonaID != "" {
			redirectPath += "#" + zonaID
		}
		http.SetCookie(w, &http.Cookie{Name: "redirect_after_login", Value: redirectPath, Path: "/"})
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	err := models.ToggleObjeto(app.DB, usuario, objetoID)
	if err != nil {
		log.Printf("ERROR: Fallo al procesar objeto '%s' para usuario '%s': %v", objetoID, usuario, err)
	} else {
		log.Printf("INFO: Usuario '%s' ha alternado estado del objeto '%s'", usuario, objetoID)
	}

	http.Redirect(w, r, "/Zonas.html#"+zonaID, http.StatusSeeOther)
}

func (app *Application) MarcarMuerto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	jefeID := r.FormValue("jefe_id")
	usuario := app.getUsuario(r)

	if usuario == "" {
		redirectPath := "/Jefes.html"
		if jefeID != "" {
			redirectPath += "#" + jefeID
		}
		http.SetCookie(w, &http.Cookie{Name: "redirect_after_login", Value: redirectPath, Path: "/"})
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	err := models.MarcarJefeDerrotado(app.DB, usuario, jefeID)
	if err != nil {
		log.Printf("ERROR: Fallo al marcar jefe '%s' derrotado por usuario '%s': %v", jefeID, usuario, err)
	} else {
		log.Printf("INFO: Usuario '%s' ha derrotado al jefe '%s'", usuario, jefeID)
	}

	proximoJefe := ordenJefes[0]
	for i, id := range ordenJefes {
		if id == jefeID {
			if i+1 < len(ordenJefes) {
				proximoJefe = ordenJefes[i+1]
			}
			break
		}
	}
	http.Redirect(w, r, "/Jefes.html#"+proximoJefe, http.StatusSeeOther)
}

func (app *Application) RevivirJefe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	jefeID := r.FormValue("jefe_id")
	usuario := app.getUsuario(r)

	if usuario == "" {
		redirectPath := "/Jefes.html"
		if jefeID != "" {
			redirectPath += "#" + jefeID
		}
		http.SetCookie(w, &http.Cookie{Name: "redirect_after_login", Value: redirectPath, Path: "/"})
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
		return
	}

	err := models.RevivirJefe(app.DB, usuario, jefeID)
	if err != nil {
		log.Printf("ERROR: Fallo al revivir jefe '%s' para usuario '%s': %v", jefeID, usuario, err)
	} else {
		log.Printf("INFO: Usuario '%s' ha revivido al jefe '%s'", usuario, jefeID)
	}

	http.Redirect(w, r, "/Jefes.html#"+jefeID, http.StatusSeeOther)
}

func (app *Application) RecomendarClaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/clases", http.StatusSeeOther)
		return
	}

	// Recoger las respuestas enviadas desde HTML
	r.ParseForm()
	p1 := r.FormValue("p1")
	p2 := r.FormValue("p2")
	p3 := r.FormValue("p3")
	p4 := r.FormValue("p4")
	p5 := r.FormValue("p5")

	// 1. Inicializar el mapa de puntuaciones a 0
	puntuaciones := map[string]int{
		"vagante": 0, "guerrero": 0, "heroe": 0, "bandido": 0, "astrologo": 0,
		"profeta": 0, "samurai": 0, "prisionero": 0, "confesor": 0, "miserable": 0,
	}

	// 2. Sumar puntos en base a las respuestas (Pregunta 1)
	switch p1 {
	case "A":
		puntuaciones["vagante"] += 2
		puntuaciones["confesor"] += 1
	case "B":
		puntuaciones["guerrero"] += 2
		puntuaciones["samurai"] += 1
	case "C":
		puntuaciones["heroe"] += 2
	case "D":
		puntuaciones["astrologo"] += 2
		puntuaciones["profeta"] += 2
	case "E":
		puntuaciones["bandido"] += 2
	case "F":
		puntuaciones["miserable"] += 2
	}

	// Pregunta 2
	switch p2 {
	case "A":
		puntuaciones["vagante"] += 2
	case "B":
		puntuaciones["heroe"] += 2
	case "C":
		puntuaciones["samurai"] += 2
		puntuaciones["guerrero"] += 1
		puntuaciones["bandido"] += 1
	case "D":
		puntuaciones["astrologo"] += 2
		puntuaciones["prisionero"] += 1
	case "E":
		puntuaciones["profeta"] += 2
		puntuaciones["confesor"] += 1
	}

	// Pregunta 3
	switch p3 {
	case "A":
		puntuaciones["vagante"] += 1
		puntuaciones["guerrero"] += 1
		puntuaciones["heroe"] += 1
		puntuaciones["samurai"] += 1
		puntuaciones["bandido"] += 1
	case "B":
		puntuaciones["astrologo"] += 1
		puntuaciones["profeta"] += 1
	case "C":
		puntuaciones["prisionero"] += 3
	case "D":
		puntuaciones["confesor"] += 3
	case "E":
		puntuaciones["miserable"] += 1
	}

	// Pregunta 4
	switch p4 {
	case "A":
		puntuaciones["bandido"] += 2
		puntuaciones["samurai"] += 2
	case "B":
		puntuaciones["guerrero"] += 2
	case "C":
		puntuaciones["vagante"] += 2
	case "D":
		puntuaciones["astrologo"] += 1
		puntuaciones["profeta"] += 1
		puntuaciones["prisionero"] += 1
		puntuaciones["confesor"] += 1
	case "E":
		puntuaciones["heroe"] += 2
	}

	// Pregunta 5
	switch p5 {
	case "A":
		puntuaciones["vagante"] += 1
	case "B":
		puntuaciones["guerrero"] += 1
		puntuaciones["samurai"] += 1
		puntuaciones["bandido"] += 1
	case "C":
		puntuaciones["astrologo"] += 1
		puntuaciones["profeta"] += 1
	case "D":
		puntuaciones["prisionero"] += 1
		puntuaciones["confesor"] += 1
	case "E":
		puntuaciones["miserable"] += 3
	}

	// 1. Identificar al ganador y su nombre legible
	maxPuntos := -1
	claseGanadora := ""
	for clase, puntos := range puntuaciones {
		if puntos > maxPuntos {
			maxPuntos = puntos
			claseGanadora = clase
		}
	}

	nombresLegibles := map[string]string{
		"vagante": "Vagante", "guerrero": "Guerrero", "heroe": "Héroe",
		"bandido": "Bandido", "astrologo": "Astrólogo", "profeta": "Profeta",
		"samurai": "Samurái", "prisionero": "Prisionero", "confesor": "Confesor",
		"miserable": "Miserable",
	}

	// 2. OBTENER DATOS DE SESIÓN Y FAVORITOS
	usuario := app.getUsuario(r)
	favoritos := make(map[string]bool)

	if usuario != "" {
		// Igual que en pages.go, obtenemos el mapa de la base de datos
		favoritos, _ = models.GetFavoritos(app.DB, usuario)
	}

	data := ClasesPageData{
		Usuario:          usuario,
		Favoritos:        favoritos,
		ResultadoId:      claseGanadora,
		ResultadoOraculo: nombresLegibles[claseGanadora],
	}

	// 3. RENDERIZAR USANDO EL MÉTODO DE TU APP
	app.renderTemplate(w, "Clases.html", data)
}
