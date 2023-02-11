-- To make tables on database use:
CREATE TABLE arts ( art_id INT AUTO_INCREMENT PRIMARY KEY, art_name VARCHAR(50) UNIQUE );

CREATE TABLE artists ( artist_id INT AUTO_INCREMENT PRIMARY KEY, artist_name VARCHAR(50) UNIQUE );

CREATE TABLE galleries ( gallery_id INT AUTO_INCREMENT PRIMARY KEY, gallery_name VARCHAR(50) UNIQUE );

CREATE TABLE artist_art ( artist_id INT not null, art_id INT not null, PRIMARY KEY (artist_id , art_id ) );

CREATE TABLE artist_gallery ( artist_id INT not null, gallery_id INT not null, PRIMARY KEY (artist_id , gallery_id ) );