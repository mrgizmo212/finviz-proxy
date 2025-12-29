package pkg

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

func fetchAllNews(ctx context.Context, isElite bool) ([]byte, error) {
	url := "https://finviz.com/news.ashx"
	if isElite {
		url = "https://elite.finviz.com/news.ashx"
	}

	var html string
	var err error

	for i := 0; i < 5; i++ {
		// check if parent context is done
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// get browser context
		bCtx, err := GetBrowserCtx()
		if err != nil {
			slog.Error("fetchAllNews get browser context failed", "err", err)
			return nil, err
		}

		// create tab
		tabCtx, cancel := chromedp.NewContext(bCtx)

		// set timeout for this attempt
		tCtx, tCancel := context.WithTimeout(tabCtx, 30*time.Second)

		slog.Info("fetchAllNews chromedp start", "url", url)
		err = chromedp.Run(tCtx,
			chromedp.ActionFunc(func(ctx context.Context) error {
				_, _, _, _, err := page.Navigate(url).Do(ctx)
				return err
			}),
			chromedp.WaitVisible(`.news_time-table`, chromedp.ByQuery),
			chromedp.OuterHTML(`html`, &html),
		)

		tCancel()
		cancel()

		if err == nil && html != "" {
			slog.Info("fetchAllNews chromedp finish", "url", url)
			break
		}

		slog.Warn("fetchAllNews attempt failed", "attempt", i+1, "err", err)
		time.Sleep(time.Second)
	}

	if err != nil {
		slog.Error("fetchAllNews chromedp run failed", "err", err)
		return nil, err
	}

	return []byte(html), nil
}

type Record struct {
	Date  string `json:"date"` // Jan-02 2006
	Title string `json:"title"`
	URL   string `json:"url"`
}

func parseLinks(table *goquery.Selection) []Record {
	loc, _ := time.LoadLocation("America/New_York")
	today := time.Now().UTC().In(loc)
	var records []Record
	table.Find("tr.news_table-row").Each(func(i int, tr *goquery.Selection) {
		a := tr.Find("a")
		href, exists := a.Attr("href")
		if !exists {
			return
		}
		date := strings.TrimSpace(tr.Find("td.news_date-cell").Text())
		if strings.HasSuffix(date, "AM") || strings.HasSuffix(date, "PM") {
			// if 05:30AM, format today
			date = today.Format("Jan-02 2006")
		} else {
			date += " " + today.Format("2006") // add year
		}
		text := a.Text()
		records = append(records, Record{
			Date:  date,
			Title: text,
			URL:   href,
		})
	})
	return records
}

func parseNewsAndBlogs(page []byte) ([]Record, []Record, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		slog.Error("failed to parse news and blogs from page", "err", err)
		return nil, nil, err
	}
	var newsTable *goquery.Selection
	var blogsTable *goquery.Selection
	doc.Find("table.styled-table-new").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			newsTable = s
		}
		if i == 1 {
			blogsTable = s
		}
	})
	if newsTable == nil || blogsTable == nil {
		return nil, nil, errors.New("failed to find news and blogs tables")
	}
	return parseLinks(newsTable), parseLinks(blogsTable), nil
}

func FetchAndParseNewsAndBlogs(ctx context.Context, isElite bool) ([]Record, []Record, error) {
	// fetch page
	page, err := fetchAllNews(ctx, isElite)
	if err != nil {
		slog.Error("failed to fetch news and blogs", "err", err)
		return nil, nil, err
	}
	// parse table
	news, blogs, err := parseNewsAndBlogs(page)
	if err != nil {
		slog.Error("failed to parse news and blogs", "err", err)
		return nil, nil, err
	}
	return news, blogs, nil
}
