package pkg

import (
	"context"
	"log/slog"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func EliteLogin(ctx context.Context, email string, password string) (bool, error) {
	url := "https://finviz.com/login-email?remember=true"
	var err error

	for i := 0; i < 5; i++ {
		// check if parent context is done
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		bCtx, err := GetBrowserCtx()
		if err != nil {
			slog.Error("login get browser context failed", "err", err)
			return false, err
		}

		// create new tab from shared browser
		tabCtx, cancel := chromedp.NewContext(bCtx)

		// set timeout for this attempt
		tCtx, tCancel := context.WithTimeout(tabCtx, 30*time.Second)

		// run login
		slog.Info("chromedp login start", "url", url)
		err = chromedp.Run(tCtx,
			chromedp.ActionFunc(func(ctx context.Context) error {
				_, _, _, _, err := page.Navigate(url).Do(ctx)
				return err
			}),
			chromedp.WaitVisible(`input[name="email"]`, chromedp.ByQuery),
			chromedp.SendKeys(`input[name="email"]`, email, chromedp.ByQuery),
			chromedp.SendKeys(`input[name="password"]`, password, chromedp.ByQuery),
			chromedp.Click(`button[type="submit"]`, chromedp.ByQuery),
			chromedp.WaitVisible(`#account-dropdown`, chromedp.ByID),
		)

		tCancel()
		cancel() // close tab

		if err == nil {
			slog.Info("login successful and cookies synced")
			return true, nil
		}

		slog.Warn("login attempt failed", "attempt", i+1, "err", err)
		time.Sleep(time.Second)
	}

	return false, err
}