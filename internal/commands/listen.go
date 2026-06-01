package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jpbede/ratiocheck/pkg/ratio"
	"github.com/urfave/cli/v3"
	"log"
	"net/http"
	"os"
)

type request struct {
	URL  string `json:"url"`
	HTML string `json:"html"`
}

func Listen() *cli.Command {
	return &cli.Command{
		Name:    "listen",
		Aliases: []string{"l"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Value:       "3000",
				DefaultText: "Port to listen on",
			},
		},
		Usage:  "Listen on given port for requests",
		Action: runListen,
	}
}

func runListen(_ context.Context, c *cli.Command) error {
	port := c.String("port")

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Post("/url", urlHandler)
	r.Post("/html", htmlHandler)

	log.Printf("Listening on port :%s\n", port)
	return http.ListenAndServe(":"+port, r)
}

func htmlHandler(w http.ResponseWriter, r *http.Request) {
	var input request
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid json"})
		return
	}

	if input.HTML == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "missing html"})
		return
	}

	file, err := os.CreateTemp("", "ratiocheck-*.html")
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		log.Printf("[%s] %v", r.Header.Get("X-Request-Id"), err)
		return
	}
	filename := file.Name()
	defer os.Remove(filename)

	if _, err := file.WriteString(input.HTML); err != nil {
		file.Close()
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		log.Printf("[%s] %v", r.Header.Get("X-Request-Id"), err)
		return
	}

	if err := file.Close(); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		log.Printf("[%s] %v", r.Header.Get("X-Request-Id"), err)
		return
	}

	if resp, err := ratio.Get(r.Context(), fmt.Sprintf("file://%s", filename)); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		log.Printf("[%s] %v", r.Header.Get("X-Request-Id"), err)
	} else {
		render.JSON(w, r, resp)
	}
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	var input request
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid json"})
		return
	}

	if input.URL == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid url"})
		return
	}

	if resp, err := ratio.Get(r.Context(), input.URL); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
	} else {
		render.JSON(w, r, resp)
	}
}
