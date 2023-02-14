package controllers

import (
	"creator/database"
	"creator/models"
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

func (gc *GalleryController) CreateGallery(g *models.Gallery) error {
	db, err := database.Connect()
	if err != nil {
		log.Println("SQL DB Connection Failed")
		return err
	}
	defer db.Close()
	database.PingDB(db)
	_, err1 := db.Exec(`INSERT INTO galleries (gallery_name) VALUES (?)`, g.Name)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	return nil
}
func (gc *GalleryController) FindGalleryDB(galleryName string) (*models.Gallery, error) {
	db, err := database.Connect()
	if err != nil {
		log.Println("SQL DB Connection Failed")
		return nil, err
	}
	defer db.Close()
	database.PingDB(db)
	g := &models.Gallery{}
	err1 := db.QueryRow(`SELECT galleries.id FROM galleries WHERE galleries.gallery_name = ?`, galleryName).Scan(&g.ID)
	if err1 != nil {
		log.Println(err1)
		return nil, err1
	}
	return g, nil
}
func (gc *GalleryController) UpdateGalleryDB(g *models.Gallery, newGalleryName string) error {
	db, err := database.Connect()
	if err != nil {
		log.Println("SQL DB Connection Failed")
		return err
	}
	defer db.Close()
	database.PingDB(db)
	_, err1 := db.Exec(`UPDATE galleries SET gallery_name = ? WHERE id = ?`, newGalleryName, g.ID)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	return nil
}
func (gc *GalleryController) DeleteArtistDB(artist *models.Artist, gallery *models.Gallery) error {
	db, err := database.Connect()
	if err != nil {
		log.Println("SQL DB Connection Failed")
		return err
	}
	defer db.Close()
	database.PingDB(db)
	_, err1 := db.Exec(`DELETE FROM artist_gallery WHERE artist_gallery.artist_id = ? and artist_gallery.gallery_id = ?`, artist.ID, gallery.ID)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	return nil
}

func (gc *GalleryController) GalleryCreation(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var galleryName string = vars["gallery"]
	gallery := &models.Gallery{Name: galleryName}
	gc.CreateGallery(gallery)
	resp := make(map[string]string)
	resp["message"] = `Gallery ` + galleryName + ` created successfully`
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}

func (gc *GalleryController) RemoveArtistFromGal(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var artistName string = vars["artist"]
	var galleryName string = vars["gallery"]
	gallery, err := gc.FindGalleryDB(galleryName)
	if err != nil {
		log.Println(err)
	}
	artistC := &ArtistController{}
	artist, err := artistC.FindArtistDB(artistName)
	if err != nil {
		log.Println(err)
	}
	err1 := gc.DeleteArtistDB(artist, gallery)
	if err1 != nil {
		log.Println(err1)
	}

	resp := make(map[string]string)
	resp["message"] = `Artist:` + artistName + `is deleted from Gallery:` + galleryName
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}

func (gc *GalleryController) GalleryUpdate(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var galleryName string = vars["gallery"]
	var newGalleryName string = vars["newgallery"]

	g, err := gc.FindGalleryDB(galleryName)
	if err != nil {
		log.Println(err)
	}
	err = gc.UpdateGalleryDB(g, newGalleryName)
	if err != nil {
		log.Println(err)
	}

	resp := make(map[string]string)
	resp["message"] = `Gallery ` + galleryName + ` was renamed to ` + newGalleryName + ` successfully`
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}
