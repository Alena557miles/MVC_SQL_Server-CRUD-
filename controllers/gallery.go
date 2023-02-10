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

type GalleryController struct {
	Galleries []*models.Gallery
	router    *mux.Router
}

func (gc *GalleryController) RegisterRouter(r *mux.Router) {
	gc.router = r
}

func (gc *GalleryController) RegisterActions() {
	// CREATE GALLERY
	// localhost:8080/creategallery/Tokio
	gc.router.HandleFunc("/creategallery/{gallery}", gc.GalleryCreation)

	// DELETE AN ARTIST FROM GALLERY
	// localhost:8080/artist/delete/Fillip/Tokio
	gc.router.HandleFunc("/artist/delete/{artist}/{gallery}", gc.RemoveArtistFromGal)

	// RENAME GALLERY
	// localhost:8080/renamegallery/Tokio/JapaneTreasure
	gc.router.HandleFunc("/renamegallery/{gallery}/{newgallery}", gc.GalleryUpdate)
}

func (gc *GalleryController) CreateGallery(db *sql.DB, g *models.Gallery) {
	_, err := db.Exec(`INSERT INTO galleries (gallery_name) VALUES (?)`, g.Name)
	if err != nil {
		log.Fatal(err)

	}
}
func (gc *GalleryController) FindGalleryDB(db *sql.DB, galleryName string) (*models.Gallery, error) {
	g := &models.Gallery{}
	err := db.QueryRow(`SELECT galleries.id FROM galleries WHERE galleries.gallery_name = ?`, galleryName).Scan(&g.ID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return g, nil
}
func (gc *GalleryController) UpdateGalleryDB(db *sql.DB, g *models.Gallery, newGalleryName string) error {
	_, err := db.Exec(`UPDATE galleries SET gallery_name = ? WHERE id = ?`, newGalleryName, g.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
func (gc *GalleryController) DeleteArtistDB(db *sql.DB, artist *models.Artist, gallery *models.Gallery) error {
	_, err := db.Exec(`DELETE FROM artist_gallery WHERE artist_gallery.artist_id = ? and artist_gallery.gallery_id = ?`, artist.ID, gallery.ID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (gc *GalleryController) GalleryCreation(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var galleryName string = vars["gallery"]

	gallery := &models.Gallery{Name: galleryName}
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("SQL DB Connection Failed")
		return
	}
	defer db.Close()
	database.PingDB(db)

	gc.CreateGallery(db, gallery)

	resp := make(map[string]string)
	resp["message"] = `Gallery ` + galleryName + ` created successfully`
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}

func (gc *GalleryController) RemoveArtistFromGal(rw http.ResponseWriter, r *http.Request) {
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

	gallery, err := gc.FindGalleryDB(db, galleryName)
	if err != nil {
		log.Fatal(err)
	}
	artistC := &ArtistController{}
	artist, err := artistC.FindArtistDB(db, artistName)
	if err != nil {
		log.Fatal(err)
	}
	err = gc.DeleteArtistDB(db, artist, gallery)
	if err != nil {
		log.Fatal(err)
	}

	resp := make(map[string]string)
	resp["message"] = `Artist:` + artistName + `is deleted from Gallery:` + galleryName
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}

func (gc *GalleryController) GalleryUpdate(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var galleryName string = vars["gallery"]
	var newGalleryName string = vars["newgallery"]

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("SQL DB Connection Failed")
		return
	}
	defer db.Close()
	database.PingDB(db)

	g, err := gc.FindGalleryDB(db, galleryName)
	if err != nil {
		log.Fatal(err)
	}
	err = gc.UpdateGalleryDB(db, g, newGalleryName)
	if err != nil {
		log.Fatal(err)
	}

	resp := make(map[string]string)
	resp["message"] = `Gallery ` + galleryName + ` was renamed to ` + newGalleryName + ` successfully`
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}
