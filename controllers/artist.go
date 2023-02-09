package controllers

import (
	"creator/database"
	"creator/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ArtistController struct {
	artists []*models.Artist
	router  *mux.Router
}

func (ac *ArtistController) RegisterRouter(r *mux.Router) {
	ac.router = r
}

func (ac *ArtistController) RegisterActions() {
	// CREATE AN ARTIST
	// localhost:8080/createartist/Fillip
	ac.router.HandleFunc("/createartist/{artist}", ac.Registration)

	//REGISTRATION AN ARTIST ON THE GALLERY
	// localhost:8080/artist/register/Fillip/Tokio
	ac.router.HandleFunc("/artist/register/{artist}/{gallery}", ac.ArtistRegistration)
}

func (ac *ArtistController) CreateArtist(a *models.Artist) {
	ac.artists = append(ac.artists, a)
}

func (ac *ArtistController) FindArtist(name string) *models.Artist {
	for _, artist := range ac.artists {
		if artist.Name == name {
			return artist
		}
	}
	return nil
}

func (ac *ArtistController) Registration(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var artistName string = vars["artist"]
	artist := &models.Artist{Name: artistName, OnGallery: false}
	ac.CreateArtist(artist)

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("SQL DB Connection Failed")
		return
	}
	defer db.Close()
	database.PingDB(db)

	_, err = db.Exec(`INSERT INTO artists (artist_name) VALUES (?)`, artistName)
	if err != nil {
		log.Fatal(err)
	}

	jsonResp, err := json.Marshal(artist)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}

func (ac *ArtistController) ArtistRegistration(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var artistName string = vars["artist"]
	var galleryName string = vars["gallery"]
	artist := ac.FindArtist(artistName)
	if err := ac.FindArtist(artistName); err != nil {
		galleryC := &GalleryController{}
		gallery := galleryC.FindGallery(galleryName)
		galleryC.RegisterArtist(gallery, artist)
	}
	resp := make(map[string]string)
	resp["message"] = `Artist: ` + artistName + `is registered on Gallery:` + galleryName
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}
