package commands

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jpbede/ratiocheck/pkg/ratio"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
)

type request struct {
	URL string `json:"url"`
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

func runListen(c *cli.Context) error {
	port := c.String("port")

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Post("/ratio", ratioHandler)

	log.Printf("Listening on port :%s\n", port)
	return http.ListenAndServe(":"+port, r)
}

func ratioHandler(w http.ResponseWriter, r *http.Request) {
	var input request
	json.NewDecoder(r.Body).Decode(&input)

	if input.URL == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid url"})
		return
	}

	if resp, err := ratio.Get(context.Background(), input.URL); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
	} else {
		render.JSON(w, r, resp)
	}
}
