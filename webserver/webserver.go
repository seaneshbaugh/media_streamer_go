package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
)

const version string = "0.0.1"
const serverUA = "Media Streamer (Go Edition v. " + version + ")"
var publicDirectory *string
var listenPort *string
var useGZip *bool
const maxBufferSize = 4096

var mediaDirectory *string

func min(x int64, y int64) int64 {
	if x < y {
		return x
	}

	return y
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", serverUA)

	filePath := path.Join((*publicDirectory), path.Clean(r.URL.Path))

	fileExists, _ := FileExists(filePath);

	if fileExists == true && r.URL.Path != "/" {
		file, _ := os.Open(filePath)

		defer file.Close()

		fileInfo, _ := file.Stat()

		if fileInfo.IsDir() || fileInfo.Mode()&os.ModeSocket != 0 {
			http.Error(w, "403 Forbidden", 403)

			return
		} else {
			ServeFile(filePath, w, r)
		}
	} else {
		pwd := *mediaDirectory

		requestPath := []string{}

		for _, part := range strings.Split(r.URL.Path[1:], "/") {
			if part != "" {
				requestPath = append(requestPath, part)
			}
		}

		switch len(requestPath) {
		case 0:
			w.Header().Set("Content-Type", "text/html")

			artistDirectories, err := ioutil.ReadDir(pwd)

			templateData := make(map[string]interface{})

			if err == nil {
				artists := []string{}

				for _, artistDirectory := range artistDirectories {
					if artistDirectory.IsDir() {
						artists = append(artists, artistDirectory.Name())
					}
				}

				indexTemplate := template.Must(template.ParseFiles("templates/layout.html", "templates/index.html"))

				templateData["pwd"] = pwd
				templateData["Artists"] = artists

				indexTemplate.Execute(w, templateData)

				fmt.Printf("%s requested %s\n", r.RemoteAddr, pwd)
			} else {
				errorTemplate := template.Must(template.ParseFiles("templates/layout.html", "templates/500.html"))

				templateData["Error"] = err

				errorTemplate.Execute(w, templateData)

				fmt.Printf("Error: %s. %s requested %s\n", err, r.RemoteAddr, r.URL)
			}
		case 1:
			w.Header().Set("Content-Type", "text/html")

			pwd = path.Join(pwd, path.Clean(requestPath[0]))

			albumDirectories, err := ioutil.ReadDir(pwd)

			templateData := make(map[string]interface{})

			if err == nil {
				albums := []string{}

				for _, albumDirectory := range albumDirectories {
					if albumDirectory.IsDir() {
						albums = append(albums, albumDirectory.Name())
					}
				}

				artistTemplate := template.Must(template.ParseFiles("templates/layout.html", "templates/artist.html"))

				templateData["pwd"] = pwd
				templateData["Artist"] = requestPath[0]
				templateData["Albums"] = albums

				artistTemplate.Execute(w, templateData)

				fmt.Printf("%s requested %s\n", r.RemoteAddr, pwd)
			} else {
				errorTemplate := template.Must(template.ParseFiles("templates/layout.html", "templates/500.html"))

				templateData["Error"] = err

				errorTemplate.Execute(w, templateData)

				fmt.Printf("Error: %s. %s requested %s\n", err, r.RemoteAddr, r.URL)
			}
		case 2:
			w.Header().Set("Content-Type", "text/html")

			pwd = path.Join(pwd, path.Clean(requestPath[0]))
			pwd = path.Join(pwd, path.Clean(requestPath[1]))

			songFiles, err := ioutil.ReadDir(pwd)

			templateData := make(map[string]interface{})

			if err == nil {
				songs := []string{}

				for _, songFile := range songFiles {
					if !songFile.IsDir() {
						songs = append(songs, songFile.Name())
					}
				}

				albumTemplate := template.Must(template.ParseFiles("templates/layout.html", "templates/album.html"))

				templateData["pwd"] = pwd
				templateData["Artist"] = requestPath[0]
				templateData["Album"] = requestPath[1]
				templateData["Songs"] = songs

				albumTemplate.Execute(w, templateData)

				fmt.Printf("%s requested %s\n", r.RemoteAddr, pwd)
			} else {
				errorTemplate := template.Must(template.ParseFiles("templates/layout.html", "templates/500.html"))

				templateData["Error"] = err

				errorTemplate.Execute(w, templateData)

				fmt.Printf("Error: %s. %s requested %s\n", err, r.RemoteAddr, r.URL)
			}
		case 3:
			pwd = path.Join(pwd, path.Clean(requestPath[0]))
			pwd = path.Join(pwd, path.Clean(requestPath[1]))
			filePath := path.Join(pwd, path.Clean(requestPath[2]))

			ServeFile(filePath, w, r)
		}
	}
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func ServeFile(filePath string, w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(filePath)

	if err != nil {
		http.Error(w, "404 Not Found", 404)

		return
	}

	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		http.Error(w, "500 Internal Server Error", 500)

		return
	}

	if fileInfo.IsDir() || fileInfo.Mode()&os.ModeSocket != 0 {
		http.Error(w, "403 Forbidden", 403)

		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	mimeType := mime.TypeByExtension(path.Ext(filePath))

	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	buffer := make([]byte, min(maxBufferSize, fileInfo.Size()))

	n := 0

	for err == nil {
		n, err = file.Read(buffer)
		w.Write(buffer[0:n])
	}
}

func main() {
	fmt.Printf(">> Starting Media Streamer Webserver (Go Edition v. " + version + ")\n")

	pwd, err := os.Getwd()

	if err != nil {
		fmt.Printf("FATAL: Could not get current working directory!\n")

		return
	}

	publicDirectory = flag.String("d", pwd + "/public", "Public Directory")
	listenPort = flag.String("p", "4568", "Listen Port")
	useGZip = flag.Bool("c", true, "Enable GZip compression")
	mediaDirectory = flag.String("m", "/Users/seshbaugh/Music/iTunes/iTunes Media/", "Media Directory")

	flag.Parse()

	fmt.Printf(">> Go application starting on http://0.0.0.0:" + *listenPort + "\n")
	fmt.Printf(">> ctrl+c to shutdown server\n")
	fmt.Printf(">> pid=" + strconv.Itoa(os.Getpid()) + "\n")

	http.ListenAndServe(":" + *listenPort, http.HandlerFunc(Handler))
}