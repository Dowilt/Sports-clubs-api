package models

type Club struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    City      string `json:"city"`
    Titles    int    `json:"titles_count"`
    AvgAge    int    `json:"avg_age"`
}