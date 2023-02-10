package models

type Gallery struct {
	ID      string    `json:"gallery_id"`
	Name    string    `json:"gallery_name"`
	Artists []*Artist `json:"gallery_artists"`
}

func (g *Gallery) DeleteArtist(artist *Artist) []*Artist {
	for i, a := range g.Artists {
		if artist == a {
			g.Artists = append(g.Artists[:i], g.Artists[i+1:]...)
			break
		}
	}
	return g.Artists
}
