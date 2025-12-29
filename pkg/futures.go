package pkg

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/chromedp/chromedp"
)

type FutureQuota struct {
	Label     string  `json:"label"`
	Ticker    string  `json:"ticker"`
	Last      float64 `json:"last"`
	Change    float64 `json:"change"`
	PrevClose float64 `json:"prevClose"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
}

func FetchAllFutures(ctx context.Context, isElite bool) (map[string]FutureQuota, error) {
	url := "https://finviz.com/api/futures_all.ashx?timeframe=NO"
	if isElite {
		url = "https://elite.finviz.com/api/futures_all.ashx?timeframe=NO"
	}

	var bodyContent string
	var err error

	for i := 0; i < 5; i++ {
		// check if parent context is done
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// get browser context
		bCtx, err := GetBrowserCtx()
		if err != nil {
			slog.Error("FetchAllFutures get browser context failed", "err", err)
			return nil, err
		}

		// create tab
		tabCtx, cancel := chromedp.NewContext(bCtx)

		// set timeout for this attempt
		tCtx, tCancel := context.WithTimeout(tabCtx, 30*time.Second)

		slog.Info("FetchAllFutures chromedp start", "url", url)
		err = chromedp.Run(tCtx,
			chromedp.Navigate(url),
			chromedp.Evaluate(`document.body.innerText`, &bodyContent),
		)

		tCancel()
		cancel()

		if err == nil && bodyContent != "" {
			slog.Info("FetchAllFutures chromedp finish", "url", url)
			break
		}

		slog.Warn("FetchAllFutures attempt failed", "attempt", i+1, "err", err)
		time.Sleep(time.Second)
	}

	if err != nil {
		slog.Error("FetchAllFutures chromedp run failed", "err", err)
		return nil, err
	}

	if bodyContent == "" {
		slog.Error("FetchAllFutures received empty body")
		return nil, nil
	}

	// unmarshal to map
	ret := make(map[string]FutureQuota)
	err = json.Unmarshal([]byte(bodyContent), &ret)
	if err != nil {
		slog.Error("FetchAllFutures json decode response", "err", err, "body", bodyContent)
		return nil, err
	}
	return ret, nil
}
