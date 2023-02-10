package controllers

import (
	"creator/database"
	"creator/models"
	"database/sql"
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

func (ac *ArtistController) CreateArtistDB(db *sql.DB, artistName string) error {
	_, err := db.Exec(`INSERT INTO artists (artist_name) VALUES (?)`, artistName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (ac *ArtistController) FindArtistDB(db *sql.DB, artistName string) (*models.Artist, error) {
	artist := &models.Artist{}
	err := db.QueryRow(`SELECT artists.id FROM artists WHERE artists.artist_name = ?`, artistName).Scan(&artist.ID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return artist, nil
}

func (ac *ArtistController) RegisterArtistToGallery(db *sql.DB, artist *models.Artist, gallery *models.Gallery) error {
	_, err := db.Exec(`INSERT INTO artist_gallery VALUES (?,?)`, artist.ID, gallery.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (ac *ArtistController) Registration(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var artistName string = vars["artist"]
	artist := &models.Artist{Name: artistName, OnGallery: false}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("SQL DB Connection Failed")
		return
	}
	defer db.Close()
	database.PingDB(db)
	err = ac.CreateArtistDB(db, artistName)
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

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("SQL DB Connection Failed")
		return
	}
	defer db.Close()
	database.PingDB(db)

	galleryC := &GalleryController{}
	gallery, err := galleryC.FindGalleryDB(db, galleryName)
	artist, err := ac.FindArtistDB(db, artistName)
	if err != nil {
		log.Fatal(err)
	}
	err = ac.RegisterArtistToGallery(db, artist, gallery)
	if err != nil {
		log.Fatal(err)
	}

	resp := make(map[string]string)
	resp["message"] = `Artist: ` + artistName + `is registered on Gallery:` + galleryName
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}
