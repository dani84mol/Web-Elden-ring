package models

import "database/sql"

// ========================
// CAPA DE ACCESO A DATOS (MVC)
// ========================

func MarcarJefeDerrotado(db *sql.DB, usuario string, jefeID string) error {
	_, err := db.Exec("INSERT OR IGNORE INTO jefes_derrotados (usuario, jefe_id) VALUES (?, ?)", usuario, jefeID)
	return err
}

func RevivirJefe(db *sql.DB, usuario string, jefeID string) error {
	_, err := db.Exec("DELETE FROM jefes_derrotados WHERE usuario = ? AND jefe_id = ?", usuario, jefeID)
	return err
}

func ToggleFavorito(db *sql.DB, usuario string, claseID string) error {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM favoritos WHERE usuario = ? AND clase_id = ?", usuario, claseID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists > 0 {
		_, err = db.Exec("DELETE FROM favoritos WHERE usuario = ? AND clase_id = ?", usuario, claseID)
	} else {
		_, err = db.Exec("INSERT INTO favoritos (usuario, clase_id) VALUES (?, ?)", usuario, claseID)
	}
	return err
}

func ToggleObjeto(db *sql.DB, usuario string, objetoID string) error {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM objetos_obtenidos WHERE usuario = ? AND objeto_id = ?", usuario, objetoID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists > 0 {
		_, err = db.Exec("DELETE FROM objetos_obtenidos WHERE usuario = ? AND objeto_id = ?", usuario, objetoID)
	} else {
		_, err = db.Exec("INSERT INTO objetos_obtenidos (usuario, objeto_id) VALUES (?, ?)", usuario, objetoID)
	}
	return err
}

// Consultas de obtención de datos
func GetFavoritos(db *sql.DB, usuario string) (map[string]bool, error) {
	favsMap := make(map[string]bool)
	rows, err := db.Query("SELECT clase_id FROM favoritos WHERE usuario = ?", usuario)
	if err != nil {
		return favsMap, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			favsMap[id] = true
		}
	}
	return favsMap, nil
}

func GetDerrotados(db *sql.DB, usuario string) (map[string]bool, error) {
	derrotadosMap := make(map[string]bool)
	rows, err := db.Query("SELECT jefe_id FROM jefes_derrotados WHERE usuario = ?", usuario)
	if err != nil {
		return derrotadosMap, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			derrotadosMap[id] = true
		}
	}
	return derrotadosMap, nil
}

func GetObtenidos(db *sql.DB, usuario string) (map[string]bool, error) {
	obtenidosMap := make(map[string]bool)
	rows, err := db.Query("SELECT objeto_id FROM objetos_obtenidos WHERE usuario = ?", usuario)
	if err != nil {
		return obtenidosMap, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			obtenidosMap[id] = true
		}
	}
	return obtenidosMap, nil
}

// En models/queries.go (añadir al final)
type Nota struct {
	ID        int    `json:"id"`
	Titulo    string `json:"titulo"`
	Contenido string `json:"contenido"`
}

func GetNotas(db *sql.DB, usuario string) ([]Nota, error) {
	rows, err := db.Query("SELECT id, titulo, contenido FROM notas WHERE usuario = ?", usuario)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notas []Nota
	for rows.Next() {
		var n Nota
		if err := rows.Scan(&n.ID, &n.Titulo, &n.Contenido); err == nil {
			notas = append(notas, n)
		}
	}
	return notas, nil
}

func CreateNota(db *sql.DB, usuario, titulo, contenido string) (int, error) {
	res, err := db.Exec("INSERT INTO notas (usuario, titulo, contenido) VALUES (?, ?, ?)", usuario, titulo, contenido)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func UpdateNota(db *sql.DB, id int, usuario, titulo, contenido string) error {
	_, err := db.Exec("UPDATE notas SET titulo = ?, contenido = ? WHERE id = ? AND usuario = ?", titulo, contenido, id, usuario)
	return err
}

func DeleteNota(db *sql.DB, id int, usuario string) error {
	_, err := db.Exec("DELETE FROM notas WHERE id = ? AND usuario = ?", id, usuario)
	return err
}
