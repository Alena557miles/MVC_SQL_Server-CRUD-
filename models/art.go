package models

type Art struct {
	ID    int    `json:"art_id"`
	Name  string `json:"art_name"`
	Owner string `json:"art_owner"`
}

func (a *Art) IsntAssigned() bool {
	return a.Owner == ""
}
