package pkg

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	chromeAllocCtx context.Context
	browserCtx     context.Context
	browserMu      sync.Mutex
)

func GetBrowserCtx() (context.Context, error) {
	browserMu.Lock()
	defer browserMu.Unlock()

	if browserCtx != nil && browserCtx.Err() == nil {
		return browserCtx, nil
	}
	browserCtx, _ = chromedp.NewContext(chromeAllocCtx)
	// start the browser
	if err := chromedp.Run(browserCtx); err != nil {
		return nil, err
	}
	return browserCtx, nil
}

func InitChromeAllocator(ctx context.Context, opts ...chromedp.ExecAllocatorOption) (context.Context, context.CancelFunc) {
	if len(opts) == 0 {
		opts = append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"),
			chromedp.UserDataDir("/tmp/finviz-proxy-chromedp"), // safe for the cookie
		)
	}
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	chromeAllocCtx = allocCtx
	return allocCtx, cancel
}

func fetchFinvizPage(ctx context.Context, params string, isElite bool) ([]byte, error) {
	baseUrl := "https://finviz.com/screener.ashx?"
	if isElite {
		baseUrl = "https://elite.finviz.com/screener.ashx?"
	}
	url := baseUrl + params

	var html string
	var err error

	for i := 0; i < 5; i++ {
		// check if parent context is done
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// create new tab from shared browser
		bCtx, err := GetBrowserCtx()
		if err != nil {
			return nil, err
		}
		tabCtx, cancel := chromedp.NewContext(bCtx)

		// set timeout for this attempt
		tCtx, tCancel := context.WithTimeout(tabCtx, 30*time.Second)

		slog.Info("chromedp start", "url", url)
		err = chromedp.Run(tCtx,
			chromedp.ActionFunc(func(ctx context.Context) error {
				_, _, _, _, err := page.Navigate(url).Do(ctx)
				return err
			}),
			chromedp.WaitVisible(`#screener-content`, chromedp.ByID),
			chromedp.OuterHTML(`html`, &html),
		)

		tCancel()
		cancel() // close tab

		if err == nil {
			slog.Info("chromedp finish", "url", url)
			return []byte(html), nil
		}

		slog.Warn("fetchFinvizPage attempt failed", "attempt", i+1, "err", err)
		time.Sleep(time.Second)
	}

	return nil, err
}
