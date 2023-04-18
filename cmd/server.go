package cmd

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"xmltv/internal/xmltv"
)

type Server struct {
	producer xmltv.Producer
	router   chi.Router
}

func (s *Server) Start() {
	s.initRouter()
	s.listen()
}

func (s *Server) initRouter() {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Get("/channels-norway.xml", ChannelListHandler)
		r.Get("/{channelId}_{month}-{day}-{year}.xml.gz", ChannelHandler)
	})

	s.router = router
}

func (s *Server) listen() {
	err := http.ListenAndServe(":8080", s.router)
	if err != nil {
		log.Fatalf("unable to listen: %s", err.Error())
	}
}

func ChannelListHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := xmltv.GetChannelList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(bytes)
}

func ChannelHandler(w http.ResponseWriter, r *http.Request) {
	channelId, err := channelId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := date(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bytes, err := xmltv.GetSchedule(channelId, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	_, _ = w.Write(bytes)
}

/*
*
Finds channelId from Path
Example input: zebrahd.tv2.no_2021-03-16.xml.gz
*/
func channelId(r *http.Request) (string, error) {
	filename := path.Base(r.URL.Path)
	parts := strings.Split(filename, "_")

	if len(parts) == 0 {
		return "", fmt.Errorf("channel name could not be found from input '%s'", filename)
	}
	return parts[0], nil
}

/*
*
Finds time.Time from Path
Example input: zebrahd.tv2.no_2021-03-16.xml.gz
*/
func date(r *http.Request) (time.Time, error) {
	filename := path.Base(r.URL.Path)
	parts := strings.Split(filename, "_")

	if len(parts) <= 1 {
		return time.Time{}, fmt.Errorf("date could not be found from input '%s'", filename)
	}

	parts2 := strings.Split(parts[1], ".")

	date, err := time.Parse("2006-01-02", parts2[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("date could not be parsed: '%s'", err)
	}

	return date, nil
}
