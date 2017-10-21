package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dhowden/tag"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	const addr = "postgresql://luister@localhost:26257/luister?sslmode=disable"
	db, err := gorm.Open("postgres", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	migrate(db)
	importcsv(db)
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&Artist{})
	db.AutoMigrate(&Year{})
	db.AutoMigrate(&Album{})
	db.AutoMigrate(&Song{})
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Play{})
	db.AutoMigrate(&Rating{})
	db.AutoMigrate(&Tag{})

	fmt.Println("migrated")
}

func importcsv(db *gorm.DB) {
	// if an "import.csv" file exists, import from it
	if _, err := os.Stat("import.csv"); os.IsNotExist(err) {
		fmt.Println("no file to import")
	}
	fmt.Println("importing from import.csv...")
	file, err := os.Open("import.csv")
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		filename := record[0]
		sha1 := record[1]
		size := record[2]

		// don't care about CAP
		importFile(db, filename, sha1, size)
	}

	fmt.Println("done importing from import.csv")
}

func importFile(db *gorm.DB, filename, sha1 string, size string) {
	// get the file from hakmes
	ext := filepath.Ext(filename)
	url := hakmesURL(sha1, ext)
	data, err := getFromHakmes(url)
	if err != nil {
		fmt.Println("failed to retrieve file from hakmes", err)
		return
	}

	// parse out id3 tags
	dreader := bytes.NewReader(data)
	fmt.Println(filename, url)
	m, err := tag.ReadFrom(dreader)
	if err != nil {
		fmt.Println("failed to read tags", err)
		return
	}
	fmt.Println("----------------")
	fmt.Println("Format", m.Format())
	fmt.Println("FileType", m.FileType())
	fmt.Println("Title", m.Title())
	fmt.Println("Album", m.Album())
	fmt.Println("Artist", m.Artist())
	fmt.Println("Year", m.Year())
	track, _ := m.Track()
	disc, _ := m.Disc()
	fmt.Println("Track", track)
	fmt.Println("Disc", disc)

	// insert record(s) into DB

	var year Year
	db.FirstOrCreate(&year, Year{Year: fmt.Sprintf("%d", m.Year())})

	var artist Artist
	db.FirstOrCreate(&artist, Artist{Name: m.Artist()})

	var album Album
	db.FirstOrCreate(&album, Album{Name: m.Album(), YearID: year.ID, ArtistID: artist.ID})

	var song Song
	db.FirstOrCreate(&song, Song{
		Title:    m.Title(),
		ArtistID: artist.ID,
		AlbumID:  album.ID,
		YearID:   year.ID,
		Track:    track,
	})

	i, err := strconv.Atoi(size)
	if err != nil {
		i = 0
	}
	var file File
	db.FirstOrCreate(&file, File{
		SongID:   song.ID,
		Filename: filename,
		Format:   string(m.Format()),
		Filetype: string(m.FileType()),
		Hash:     sha1,
		Filesize: i,
	})
}

func hakmesURL(sha1, ext string) string {
	return "http://localhost:9300/file/" + sha1 + "/file" + ext
}

func getFromHakmes(url string) ([]byte, error) {
	c := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)

	if err != nil {
		return nil, err
	}
	if resp.Status != "200 OK" {
		return nil, errors.New("404, probably")
	}
	b, _ := ioutil.ReadAll(resp.Body)
	return b, nil
}
