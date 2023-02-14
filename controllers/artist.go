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

func (ac *ArtistController) CreateArtistDB(artistName string) error {
	db, err := database.Connect()
	if err != nil {
		log.Println("SQL DB Connection Failed")
		return err
	}
	defer db.Close()
	database.PingDB(db)
	_, err1 := db.Exec(`INSERT INTO artists (artist_name) VALUES (?)`, artistName)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	return nil
}

func (ac *ArtistController) FindArtistDB(artistName string) (*models.Artist, error) {
	db, err := database.Connect()
	if err != nil {
		log.Println("SQL DB Connection Failed")
		return nil, err
	}
	defer db.Close()
	database.PingDB(db)
	artist := &models.Artist{}
	err1 := db.QueryRow(`SELECT artists.id FROM artists WHERE artists.artist_name = ?`, artistName).Scan(&artist.ID)
	if err1 != nil {
		log.Println(err1)
		return nil, err1
	}
	return artist, nil
}

func (ac *ArtistController) RegisterArtistToGallery(artist *models.Artist, gallery *models.Gallery) error {
	db, err := database.Connect()
	if err != nil {
		log.Println("SQL DB Connection Failed")
		return err
	}
	defer db.Close()
	database.PingDB(db)
	_, err1 := db.Exec(`INSERT INTO artist_gallery VALUES (?,?)`, artist.ID, gallery.ID)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	return nil
}

func (ac *ArtistController) Registration(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var artistName string = vars["artist"]
	artist := &models.Artist{Name: artistName, OnGallery: false}

	err := ac.CreateArtistDB(artistName)
	if err != nil {
		log.Println(err)
	}

	jsonResp, err := json.Marshal(artist)
	if err != nil {
		log.Println("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}

func (ac *ArtistController) ArtistRegistration(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var artistName string = vars["artist"]
	var galleryName string = vars["gallery"]

	galleryC := &GalleryController{}
	gallery, err := galleryC.FindGalleryDB(galleryName)
	artist, err := ac.FindArtistDB(artistName)
	if err != nil {
		log.Println(err)
	}
	err1 := ac.RegisterArtistToGallery(artist, gallery)
	if err1 != nil {
		log.Println(err1)
	}

	resp := make(map[string]string)
	resp["message"] = `Artist: ` + artistName + `is registered on Gallery:` + galleryName
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}
