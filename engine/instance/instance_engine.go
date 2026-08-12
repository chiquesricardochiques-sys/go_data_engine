package instance

import (
	"database/sql"
	"encoding/json"

	"meu-provedor/config"
	"meu-provedor/models"
)

func InsertInstance(req models.InstanceRequest) error {
	settingsJSON, _ := json.Marshal(req.Settings)

	_, err := config.MasterDB.Exec(`
		INSERT INTO instancias_projetion
		(project_id, name, code, description, status, settings)
		VALUES (?, ?, ?, ?, ?, ?)`,
		req.ProjectID,
		req.Name,
		req.Code,
		req.Description,
		req.Status,
		settingsJSON,
	)

	return err
}

func ListInstances(projectID *int64) ([]models.Instance, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if projectID != nil {
		rows, err = config.MasterDB.Query(`
			SELECT id, project_id, name, code, description, status, settings, created_at
			FROM instancias_projetion
			WHERE project_id = ?`,
			*projectID,
		)
	} else {
		rows, err = config.MasterDB.Query(`
			SELECT id, project_id, name, code, description, status, settings, created_at
			FROM instancias_projetion`,
		)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []models.Instance

	for rows.Next() {
		var (
			i        models.Instance
			settings []byte
		)

		rows.Scan(
			&i.ID,
			&i.ProjectID,
			&i.Name,
			&i.Code,
			&i.Description,
			&i.Status,
			&settings,
			&i.CreatedAt,
		)

		_ = json.Unmarshal(settings, &i.Settings)
		instances = append(instances, i)
	}

	return instances, nil
}

func UpdateInstance(id int64, req models.InstanceRequest) error {
	settingsJSON, _ := json.Marshal(req.Settings)

	_, err := config.MasterDB.Exec(`
		UPDATE instancias_projetion
		SET name=?, code=?, description=?, status=?, settings=?
		WHERE id=?`,
		req.Name,
		req.Code,
		req.Description,
		req.Status,
		settingsJSON,
		id,
	)

	return err
}

func DeleteInstance(id int64) error {
	_, err := config.MasterDB.Exec(
		`DELETE FROM instancias_projetion WHERE id=?`,
		id,
	)
	return err
}


