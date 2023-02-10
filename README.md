# Getting Started with MySQL_Server_CRUD  App

In this project i'm implement CRUD app "Art-Artist-Gallery" that manage data on MySQL database.

##  Server will start on 
[http://localhost:8080](http://localhost:8080)

CREATE AN ART on DB use path:
#### `http://localhost:8080/createart/{artName}`

CREATE AN ARTIST on DB use path:
#### `http://localhost:8080/createartist/{artistName}`

CREATE A GALLERY on DB use path:
#### `http://localhost:8080/creategallery/{gallery}`

ASSIGN AN ART TO THE ARTIST (BY NAME) on DB use path:
#### `http://localhost:8080/artist/assign/{artist}/{art}`

REGISTRATION AN ARTIST ON THE GALLERY use path:
#### `http://localhost:8080/artist/register/{artist}/{gallery}`

DELETE AN ARTIST FROM GALLERY use path:
#### `http://localhost:8080/artist/delete/{artist}/{gallery}`

RENAME GALLERY use path:
#### `http://localhost:8080/renamegallery/{gallery}/{newgallery}`



where {art} - name of the Art

where {artist} - name of the Artist

where {gallery} - name of the Gallery

where {newgallery} - new name of the Gallery

#### To make tables on database use:

CREATE TABLE arts
(
art_id     INT AUTO_INCREMENT PRIMARY KEY,
art_name   VARCHAR(50) UNIQUE
);

CREATE TABLE artists
(
artist_id     INT AUTO_INCREMENT PRIMARY KEY,
artist_name   VARCHAR(50) UNIQUE
);

CREATE TABLE galleries
(
gallery_id     INT AUTO_INCREMENT PRIMARY KEY,
gallery_name   VARCHAR(50) UNIQUE
);

CREATE TABLE artist_art
(
artist_id     INT not null,
art_id   INT not null,
PRIMARY KEY (artist_id ,  art_id )
);

CREATE TABLE artist_gallery
(
artist_id     INT not null,
gallery_id   INT not null,
PRIMARY KEY (artist_id ,  gallery_id )
);

