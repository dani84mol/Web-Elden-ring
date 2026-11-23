// handlers/api.go
package handlers

import (
	"elden-guide/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

func (app *Application) NotasAPI(w http.ResponseWriter, r *http.Request) {
	usuario := app.getUsuario(r)
	if usuario == "" {
		http.Error(w, `{"error": "No autorizado"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// READ
		notas, _ := models.GetNotas(app.DB, usuario)
		if notas == nil {
			notas = []models.Nota{} // Para no devolver null en JSON
		}
		json.NewEncoder(w).Encode(notas)

	case http.MethodPost:
		// CREATE
		var nota models.Nota
		json.NewDecoder(r.Body).Decode(&nota)
		id, _ := models.CreateNota(app.DB, usuario, nota.Titulo, nota.Contenido)
		nota.ID = id
		w.WriteHeader(http.StatusCreated) // 201 Created
		json.NewEncoder(w).Encode(nota)

	case http.MethodPut:
		// UPDATE
		var nota models.Nota
		json.NewDecoder(r.Body).Decode(&nota)
		models.UpdateNota(app.DB, nota.ID, usuario, nota.Titulo, nota.Contenido)
		w.WriteHeader(http.StatusOK) // 200 OK
		json.NewEncoder(w).Encode(map[string]string{"status": "actualizado"})

	case http.MethodDelete:
		// DELETE
		idStr := r.URL.Query().Get("id")
		id, _ := strconv.Atoi(idStr)
		models.DeleteNota(app.DB, id, usuario)
		w.WriteHeader(http.StatusOK) // 200 OK
		json.NewEncoder(w).Encode(map[string]string{"status": "eliminado"})
	}
}
