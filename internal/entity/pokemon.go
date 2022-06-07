package entity

type Pokemon struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	Height         int64  `json:"height"`
	BaseExperience int64  `json:"base_experience"`
	Status         bool   `json:"status"`
}
