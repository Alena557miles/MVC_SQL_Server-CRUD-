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

func (ac *ArtController) CreateArtDB(artName string) error {
	db, err := database.Connect()
	if err != nil {
		log.Println("SQL DB Connection Failed")
		return err
	}
	defer db.Close()
	database.PingDB(db)
	_, err1 := db.Exec(`INSERT INTO arts (art_name) VALUES (?)`, artName)
	if err != nil {
		log.Fatal(err1)
		return err1
	}
	return nil
}

func (ac *ArtController) FindArtDB(artName string) (*models.Art, error) {
	db, err := database.Connect()
	if err != nil {
		log.Println("SQL DB Connection Failed")
		return nil, err
	}
	defer db.Close()
	database.PingDB(db)
	art := &models.Art{}
	err1 := db.QueryRow(`SELECT arts.id FROM arts WHERE arts.art_name = ?`, artName).Scan(&art.ID)
	if err1 != nil {
		log.Println(err1)
		return nil, err1
	}
	return art, nil
}

func (ac *ArtController) AssignedArtToArtist(art *models.Art, artist *models.Artist) error {
	db, err := database.Connect()
	if err != nil {
		log.Println("SQL DB Connection Failed")
		return err
	}
	defer db.Close()
	database.PingDB(db)
	_, err1 := db.Exec(`INSERT INTO artist_art VALUES (?,?)`, artist.ID, art.ID)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	return nil
}

func (ac *ArtController) ArtCreation(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var artName string = vars["art"]
	art := &models.Art{Name: artName}

	err := ac.CreateArtDB(artName)
	if err != nil {
		log.Println(err)
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

	art, err := ac.FindArtDB(artName)
	artistC := &ArtistController{}
	artist, err := artistC.FindArtistDB(artistName)
	if err != nil {
		log.Fatal(err)
	}

	err = ac.AssignedArtToArtist(art, artist)
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
