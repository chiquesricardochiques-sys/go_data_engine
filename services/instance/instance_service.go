package instance

import (
	"database/sql"
	"encoding/json"
	"errors"

	"meu-provedor/config"
	"meu-provedor/models"
)

// =======================
// Validação centralizada
// =======================
func validate(req models.InstanceRequest) error {
	if req.ProjectID <= 0 {
		return errors.New("project_id is required")
	}
	if req.ClientName == "" {
		return errors.New("client_name is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.PaymentDay < 1 || req.PaymentDay > 28 {
		return errors.New("payment_day must be between 1 and 28")
	}
	if req.Price < 0 {
		return errors.New("price must be >= 0")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Code == "" {
		return errors.New("code is required")
	}
	return nil
}

// =======================
// CREATE
// =======================
func Create(req models.InstanceRequest) error {
	if err := validate(req); err != nil {
		return err
	}

	settingsJSON, _ := json.Marshal(req.Settings)

	_, err := config.MasterDB.Exec(`
		INSERT INTO instancias_projetion
		(project_id, client_name, email, phone, price, payment_day,
		 name, code, description, status, settings)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		req.ProjectID,
		req.ClientName,
		req.Email,
		req.Phone,
		req.Price,
		req.PaymentDay,
		req.Name,
		req.Code,
		req.Description,
		req.Status,
		settingsJSON,
	)

	return err
}

// =======================
// LIST
// =======================
func List(projectID *int64) ([]models.Instance, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if projectID != nil {
		rows, err = config.MasterDB.Query(`
			SELECT id, project_id, client_name, email, phone, price, payment_day,
			       name, code, description, status, settings, created_at, updated_at
			FROM instancias_projetion
			WHERE project_id = ?`,
			*projectID,
		)
	} else {
		rows, err = config.MasterDB.Query(`
			SELECT id, project_id, client_name, email, phone, price, payment_day,
			       name, code, description, status, settings, created_at, updated_at
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

		err := rows.Scan(
			&i.ID,
			&i.ProjectID,
			&i.ClientName,
			&i.Email,
			&i.Phone,
			&i.Price,
			&i.PaymentDay,
			&i.Name,
			&i.Code,
			&i.Description,
			&i.Status,
			&settings,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(settings, &i.Settings)
		instances = append(instances, i)
	}

	return instances, nil
}

// =======================
// UPDATE
// =======================
func Update(id int64, req models.InstanceRequest) error {
	if id <= 0 {
		return errors.New("invalid instance id")
	}
	if err := validate(req); err != nil {
		return err
	}

	settingsJSON, _ := json.Marshal(req.Settings)

	_, err := config.MasterDB.Exec(`
		UPDATE instancias_projetion
		SET client_name=?, email=?, phone=?, price=?, payment_day=?,
		    name=?, code=?, description=?, status=?, settings=?
		WHERE id=?`,
		req.ClientName,
		req.Email,
		req.Phone,
		req.Price,
		req.PaymentDay,
		req.Name,
		req.Code,
		req.Description,
		req.Status,
		settingsJSON,
		id,
	)

	return err
}

// =======================
// DELETE
// =======================
func Delete(id int64) error {
	if id <= 0 {
		return errors.New("invalid instance id")
	}

	_, err := config.MasterDB.Exec(
		`DELETE FROM instancias_projetion WHERE id=?`,
		id,
	)

	return err
}
