package controllers

import (
	"creator/database"
	"creator/models"
	"encoding/json"
	"fmt"
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
	gc.router.HandleFunc("/renamegallery/{gallery}/{newgallery}", gc.RemoveArtistFromGal)
}

func (gc *GalleryController) CreateGallery(g *models.Gallery) {
	gc.Galleries = append(gc.Galleries, g)
}
func (gc *GalleryController) FindGallery(name string) *models.Gallery {
	for _, g := range gc.Galleries {
		if name == g.Name {
			return g
		}
	}
	return nil
}
func (gc *GalleryController) RegisterArtist(gallery *models.Gallery, artist *models.Artist) {
	if len(artist.Arts) > 0 {
		gc.Galleries = append(gc.Galleries, gallery)
		gallery.Artists = append(gallery.Artists, artist)
		artist.OnGallery = true
		return
	}
	if len(artist.Arts) == 0 {
		fmt.Println("We can not register an Artist without Arts")
	}
	if artist.OnGallery {
		fmt.Println("Artist are already on gallery")
	}
}
func (gc *GalleryController) DeleteArtist(gallery *models.Gallery, artist *models.Artist) {
	for _, g := range gc.Galleries {
		if g == gallery {
			gallery.DeleteArtist(artist)
		}
	}
}
func (gc *GalleryController) GalleryCreation(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var galleryName string = vars["gallery"]
	resp := make(map[string]string)
	resp["message"] = `Gallery ` + galleryName + ` created successfully`
	gallery := &models.Gallery{Name: galleryName}
	gc.CreateGallery(gallery)

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("SQL DB Connection Failed")
		return
	}
	defer db.Close()
	database.PingDB(db)

	_, err = db.Exec(`INSERT INTO galleries (gallery_name) VALUES (?)`, galleryName)
	if err != nil {
		log.Fatal(err)
	}

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
	artistC := &ArtistController{}
	artist := artistC.FindArtist(artistName)
	if err := artistC.FindArtist(artistName); err != nil {
		gallery := gc.FindGallery(galleryName)
		gc.DeleteArtist(gallery, artist)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("SQL DB Connection Failed")
		return
	}
	defer db.Close()
	database.PingDB(db)

	// find Gallery ID in DB
	g := models.Gallery{}
	err = db.QueryRow(`SELECT galleries.id FROM galleries WHERE galleries.gallery_name = ?`, galleryName).Scan(&g.ID)
	if err != nil {
		log.Fatal(err)
	}

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

	// delete an Artist from Gallery
	_, err = db.Exec(`DELETE FROM artist_gallery WHERE artist_gallery.artist_id = ?`, artst.ID)
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

	//// find Gallery ID in DB second way
	//g := models.Gallery{}
	//err = db.QueryRow(`SELECT galleries.id FROM galleries WHERE galleries.gallery_name = ?`, galleryName).Scan(&g.ID)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// find GAllery ID in DB
	row, err := db.Query(`SELECT galleries.id FROM galleries WHERE galleries.gallery_name = ?`, galleryName)
	if err != nil {
		log.Fatal(err)
	}
	row.Next()
	g := models.Artist{}
	err = row.Scan(&g.ID)
	if err != nil {
		log.Fatal(err)
	}
	row.Close()
	log.Println(g)
	//// update gallery name
	//_, err = db.Exec(`UPDATE galleries SET gallery_name = ? WHERE gallery_id = ?`, newGalleryName, g.ID)
	//if err != nil {
	//	log.Fatal(err)
	//}

	resp := make(map[string]string)
	resp["message"] = `Gallery ` + galleryName + ` was renamed to ` + newGalleryName + `successfully`
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	rw.Write(jsonResp)
}

func (gc *GalleryController) UpdateGallery(name string) {

}
