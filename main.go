package main

import (
	"creator/controllers"
	"github.com/Alena557miles/webservergo"
)

func main() {
	artistC := &controllers.ArtistController{}
	artC := &controllers.ArtController{}
	galleryC := &controllers.GalleryController{}
	webservergo.StartServer(artistC, artC, galleryC)
}
