package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	fmt.Println("chatterday started! ", time.Now().Format("2006-01-02 15:04:05"))

	mux := http.NewServeMux()

	fileServer := http.FileServer(safeStaticFileSystem{http.Dir("./static")})

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))

		tmpl.Execute(w, nil)
	})

	mux.HandleFunc("/wsreload", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err)
			return
		}

		ws.WriteMessage(websocket.TextMessage, []byte("Ready."))
	})

	log.Fatal(http.ListenAndServe(":8000", mux))

}

type safeStaticFileSystem struct {
	fs http.FileSystem
}

func (ssfs safeStaticFileSystem) Open(path string) (http.File, error) {
	f, err := ssfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := ssfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
