package main

import (
	"context"
	"encoding/json"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"gopkg.in/alecthomas/kingpin.v2"

	"log"
)

type Response struct {
	ContentArea float64 `json:"content_area"`
	ImageArea   float64 `json:"image_area"`
	Ratio       float64 `json:"ratio"`
}

type Request struct {
	URL string `json:"url"`
}

var (
	app     = kingpin.New("ratio-check", "Image-to-content ratio check via REST API")
	verbose = app.Flag("verbose", "Verbose mode").Short('v').Bool()
	port    = app.Flag("port", "Port to listen on").Short('p').Default("3000").String()
)

func main() {
	app.DefaultEnvars().Parse(os.Args[1:])

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Post("/ratio", ratioHandler)

	log.Printf("Listening on port :%s\n", *port)
	http.ListenAndServe(":"+*port, r)
}

func ratioHandler(w http.ResponseWriter, r *http.Request) {
	var input Request
	json.NewDecoder(r.Body).Decode(&input)

	if input.URL == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid url"})
		return
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	chromedp.WithPollingTimeout(5 * time.Second)

	resp := Response{}

	if err := chromedp.Run(ctx,
		chromedp.Navigate(input.URL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if _, _, _, _, _, contentSize, err := page.GetLayoutMetrics().Do(ctx); err != nil {
				return err
			} else {
				resp.ContentArea = contentSize.Width * contentSize.Height
				return nil
			}
		}),
		chromedp.QueryAfter("img", func(ctx context.Context, _ runtime.ExecutionContextID, nodes ...*cdp.Node) error {
			for _, node := range nodes {
				if model, _ := dom.GetBoxModel().WithNodeID(node.NodeID).Do(ctx); model != nil {
					resp.ImageArea += float64(model.Height) * float64(model.Width)
				}
			}
			return nil
		}),
	); err != nil {
		log.Fatal(err)
		render.JSON(w, r, err)
		return
	}

	resp.Ratio = (100 / resp.ContentArea) * resp.ImageArea

	render.JSON(w, r, resp)
}
