package models


type InstanceRequest struct {
	ProjectID   int64                  `json:"project_id"`
	ClientName  string                 `json:"client_name"`
	Email       string                 `json:"email"`
	Phone       string                 `json:"phone"`
	Price       float64                `json:"price"`
	PaymentDay  int                    `json:"payment_day"`
	Name        string                 `json:"name"`
	Code        string                 `json:"code"`
	Description string                 `json:"description"`
	Status      string                 `json:"status"`
	Settings    map[string]interface{} `json:"settings"`
}



