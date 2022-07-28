package ratio

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"time"
)

type Result struct {
	ContentArea float64 `json:"content_area"`
	ImageArea   float64 `json:"image_area"`
	Ratio       float64 `json:"ratio"`
}

func Get(parentCtx context.Context, url string) (*Result, error) {
	ctx, cancel := chromedp.NewContext(parentCtx)
	defer cancel()

	chromedp.WithPollingTimeout(5 * time.Second)

	ratio := &Result{}

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// calculate area of whole document
		chromedp.ActionFunc(func(ctx context.Context) error {
			if _, _, _, _, _, contentSize, err := page.GetLayoutMetrics().Do(ctx); err != nil {
				return err
			} else {
				ratio.ContentArea = contentSize.Width * contentSize.Height
				return nil
			}
		}),
		// calculate total image area
		chromedp.QueryAfter("img", func(ctx context.Context, _ runtime.ExecutionContextID, nodes ...*cdp.Node) error {
			for _, node := range nodes {
				if model, _ := dom.GetBoxModel().WithNodeID(node.NodeID).Do(ctx); model != nil {
					ratio.ImageArea += float64(model.Height) * float64(model.Width)
				}
			}
			return nil
		}),
	); err != nil {
		return nil, err
	}

	ratio.Ratio = (100 / ratio.ContentArea) * ratio.ImageArea

	return ratio, nil
}
