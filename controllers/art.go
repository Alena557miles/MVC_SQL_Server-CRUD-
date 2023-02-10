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

type ArtController struct {
	arts   []*models.Art
	router *mux.Router
	db     *sql.DB
}

func (ac *ArtController) RegisterRouter(r *mux.Router) {
	ac.router = r
}

func (ac *ArtController) RegisterActions() {
	// CREATE AN ART
	// localhost:8080/createart/blackCat
	ac.router.HandleFunc("/createart/{art}", ac.ArtCreation)

	//ASSIGN AN ART TO THE ARTIST (BY NAME)
	// localhost:8080/artist/assign/Fillip/blackCat
	ac.router.HandleFunc("/artist/assign/{artist}/{art}", ac.AssignArt)

}

func (ac *ArtController) CreateArtDB(db *sql.DB, artName string) error {
	_, err := db.Exec(`INSERT INTO arts (art_name) VALUES (?)`, artName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (ac *ArtController) FindArtDB(db *sql.DB, artName string) (*models.Art, error) {
	art := &models.Art{}
	err := db.QueryRow(`SELECT arts.id FROM arts WHERE arts.art_name = ?`, artName).Scan(&art.ID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return art, nil
}

func (ac *ArtController) AssignedArtToArtist(db *sql.DB, art *models.Art, artist *models.Artist) error {
	// pass data to table artist-art
	_, err := db.Exec(`INSERT INTO artist_art VALUES (?,?)`, artist.ID, art.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (ac *ArtController) ArtCreation(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var artName string = vars["art"]
	art := &models.Art{Name: artName}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("SQL DB Connection Failed")
		return
	}
	defer db.Close()
	database.PingDB(db)

	err = ac.CreateArtDB(db, artName)
	if err != nil {
		log.Fatal(err)
	}
	jsonResp, err := json.Marshal(art)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}

func (ac *ArtController) AssignArt(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var artistName string = vars["artist"]
	var artName string = vars["art"]

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("SQL DB Connection Failed")
		return
	}
	defer db.Close()
	database.PingDB(db)

	art, err := ac.FindArtDB(db, artName)
	artistC := &ArtistController{}
	artist, err := artistC.FindArtistDB(db, artistName)
	if err != nil {
		log.Fatal(err)
	}

	err = ac.AssignedArtToArtist(db, art, artist)
	if err != nil {
		log.Fatal(err)
	}

	resp := make(map[string]string)
	resp["message"] = `Art: ` + artName + ` is assigned to Artist:` + artistName
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}
