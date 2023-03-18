package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"techklepto.dev/views"

	"github.com/gorilla/mux"
)

var (
	homeView *views.View
	contactView *views.View
	imageTrackView *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func imageTrack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(imageTrackView.Render(w, nil))
}

func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	imageTrackView = views.NewView("bootstrap", "views/imageTrack.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/track", imageTrack)

	s := &http.Server{
		Addr: 			":80",
		Handler: 		r,
		IdleTimeout: 	120 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		http.Handle("/", http.FileServer(http.Dir("img")))
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	log.Println("Recieved termination, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)	
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}