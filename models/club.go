package models

type Club struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	City   string `json:"city"`
	CityID int    `json:"city_id"`
	Titles int    `json:"titles"`
	AvgAge int    `json:"avgAge"`
}

type UpdateClubRequest struct {
	Name   *string `json:"name,omitempty"`
	City   *string `json:"city,omitempty"`
	Titles *int    `json:"titles_count,omitempty"`
	AvgAge *int    `json:"avg_age,omitempty"`
}
