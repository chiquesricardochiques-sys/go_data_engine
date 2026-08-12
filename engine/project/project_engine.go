package project

import (
	"database/sql"
	"meu-provedor/config"
	"meu-provedor/models"
)

func InsertProject(req models.ProjectRequest) error {
	_, err := config.MasterDB.Exec(`
		INSERT INTO projects (name, code, api_key, type, version, status)
		VALUES (?, ?, ?, ?, ?, ?)`,
		req.Name,
		req.Code,
		req.ApiKey,
		req.Type,
		req.Version,
		req.Status,
	)
	return err
}

func ListProjects() ([]models.Project, error) {
	rows, err := config.MasterDB.Query(`
		SELECT id, name, code, api_key, type, version, status, created_at
		FROM projects`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var p models.Project
		rows.Scan(
			&p.ID,
			&p.Name,
			&p.Code,
			&p.ApiKey,
			&p.Type,
			&p.Version,
			&p.Status,
			&p.CreatedAt,
		)
		projects = append(projects, p)
	}

	return projects, nil
}

func UpdateProject(id int64, req models.ProjectRequest) error {
	_, err := config.MasterDB.Exec(`
		UPDATE projects
		SET name=?, code=?, type=?, version=?, status=?
		WHERE id=?`,
		req.Name,
		req.Code,
		req.Type,
		req.Version,
		req.Status,
		id,
	)
	return err
}

func DeleteProject(id int64) error {
	_, err := config.MasterDB.Exec(`DELETE FROM projects WHERE id=?`, id)
	return err
}

func ProjectCodeExists(code string) (bool, error) {
	var exists int
	err := config.MasterDB.QueryRow(
		`SELECT 1 FROM projects WHERE code=? LIMIT 1`,
		code,
	).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	return err == nil, err
}

