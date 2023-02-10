package models

type Artist struct {
	ID        int    `json:"artist_id"`
	Name      string `json:"artist_name"`
	Arts      []*Art `json:"artist_arts"`
	OnGallery bool   `json:"ongallery"`
}
