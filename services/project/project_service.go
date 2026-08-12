package project

import (
	"database/sql"
	"errors"
	"meu-provedor/config"
	"meu-provedor/models"
)

// Create insere um novo projeto
func Create(req models.ProjectRequest) error {
	if req.Code == "" {
		return errors.New("code is required")
	}

	exists, err := CodeExists(req.Code)
	if err != nil {
		return err
	}
	if exists {
		return models.ErrProjectCodeExists
	}

	_, err = config.MasterDB.Exec(`
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

// List retorna todos os projetos
func List() ([]models.Project, error) {
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
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Code,
			&p.ApiKey,
			&p.Type,
			&p.Version,
			&p.Status,
			&p.CreatedAt,
		)
		if err != nil {
			continue
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// Update atualiza um projeto existente (SEM alterar o code)
func Update(id int64, req models.ProjectUpdateRequest) error {
	_, err := config.MasterDB.Exec(`
		UPDATE projects
		SET name=?, api_key=?, type=?, version=?, status=?
		WHERE id=?`,
		req.Name,
		req.ApiKey,
		req.Type,
		req.Version,
		req.Status,
		id,
	)
	return err
}

// Delete remove um projeto
func Delete(id int64) error {
	_, err := config.MasterDB.Exec(`DELETE FROM projects WHERE id=?`, id)
	return err
}

// CodeExists verifica se um código de projeto já existe
func CodeExists(code string) (bool, error) {
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
