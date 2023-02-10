package controllers

import (
	"creator/database"
	"creator/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"fmt"
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

func (ac *ArtController) CreateArt(a *models.Art) {
	ac.arts = append(ac.arts, a)
}

func (ac *ArtController) FindArt(name string) *models.Art {
	for _, art := range ac.arts {
		if art.Name == name {
			return art
		}
	}
	return nil
}

func (ac *ArtController) AssignedArtToArtist(art *models.Art, artist *models.Artist) *models.Art {
	if art.IsntAssigned() {
		art.Owner = artist.Name
		artist.Arts = append(artist.Arts, art)
		fmt.Println(artist.Name)
		return art
	} else {
		fmt.Println("This art already has an owner! You'd better make your own art")
	}
	return nil
}

func (ac *ArtController) ArtCreation(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var artName string = vars["art"]
	art := &models.Art{Name: artName}
	ac.CreateArt(art)

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("SQL DB Connection Failed")
		return
	}
	defer db.Close()
	database.PingDB(db)

	_, err = db.Exec(`INSERT INTO arts (art_name) VALUES (?)`, artName)
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
	art := ac.FindArt(artName)
	if err := ac.FindArt(artName); err != nil {
		artistC := &ArtistController{}
		artist := artistC.FindArtist(artistName)
		if err := artistC.FindArtist(artistName); err != nil {
			ac.AssignedArtToArtist(art, artist)
		}
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("SQL DB Connection Failed")
		return
	}
	defer db.Close()
	database.PingDB(db)

	// find Art ID in DB
	row, err := db.Query(`SELECT arts.id FROM arts WHERE arts.art_name = ?`, artName)
	if err != nil {
		log.Fatal(err)
	}
	row.Next()
	a := models.Art{}
	err = row.Scan(&a.ID)
	if err != nil {
		log.Fatal(err)
	}
	row.Close()

	// find Artist ID in DB
	rowArtist, err := db.Query(`SELECT artists.id FROM artists WHERE artists.artist_name = ?`, artistName)
	if err != nil {
		log.Fatal(err)
	}
	rowArtist.Next()
	artst := models.Artist{}
	err = rowArtist.Scan(&artst.ID)
	if err != nil {
		log.Fatal(err)
	}
	rowArtist.Close()

	// pass data to table artist-art
	_, err = db.Exec(`INSERT INTO artist_art VALUES (?,?)`, artst.ID, a.ID)
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
