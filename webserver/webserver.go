package webserver

import(
	fileutils "../fileutils"
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

const Version string = "0.0.2"
const ServerUA = "Media Streamer (Go Edition v. " + Version + ")"
var PublicDirectory *string
var ListenPort *string
var MediaDirectory *string

const maxBufferSize = 4096

func min(x int64, y int64) int64 {
	if x < y {
		return x
	}

	return y
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", ServerUA)

	filePath := path.Join((*PublicDirectory), path.Clean(r.URL.Path))

	fileExists, _ := fileutils.FileExists(filePath);

	if fileExists == true && r.URL.Path != "/" {
		file, _ := os.Open(filePath)

		defer file.Close()

		fileInfo, _ := file.Stat()

		if fileInfo.IsDir() || fileInfo.Mode()&os.ModeSocket != 0 {
			http.Error(w, "403 Forbidden", 403)

			return
		} else {
			ServeFile(w, r, filePath)
		}
	} else {
		pwd := *MediaDirectory

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

			ServeFile(w, r, filePath)
		}
	}
}

func RenderErrorPage(w http.ResponseWriter, r *http.Request, statusCode int) {
	w.WriteHeader(statusCode)

	w.Header().Set("Content-Type", "text/html")

	templateData := make(map[string]interface{})

	templateData["StatusCode"] = strconv.FormatInt(int64(statusCode), 10)

	templateData["StatusText"] = http.StatusText(statusCode)

	errorTemplate := template.Must(template.ParseFiles("templates/layout.html", "templates/error.html"))

	errorTemplate.Execute(w, templateData)
}

func ServeFile(w http.ResponseWriter, r *http.Request, filePath string) {
	file, err := os.Open(filePath)

	if err != nil {
		RenderErrorPage(w, r, http.StatusNotFound)

		return
	}

	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		RenderErrorPage(w, r, http.StatusInternalServerError)

		return
	}

	if fileInfo.Mode()&os.ModeDir != 0 || fileInfo.Mode()&os.ModeSocket != 0 || fileInfo.Mode()&os.ModeDevice != 0 {
		RenderErrorPage(w, r, http.StatusForbidden)

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