package extract

import (
	"context"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/chromedp/chromedp"
	"math/rand"
	"web-scraper/proxies"
)

func Extract(page, sel string, fields []string) ([][]string, error) {
	// define the settings
	options := []chromedp.ExecAllocatorOption{
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.UserAgent(browser.Random()),
	}

	// specify a new context set up for NewContext
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	// initialize a controllable Chrome instance
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	return DoExtract(ctx, page, sel, fields, false)
}

func ExtractProxy(page, sel string, fields []string) (result [][]string, err error) {
	proxyList := proxies.GetProxies()

	for i := 0; i < 5; i++ {
		if result, err = DoExtractProxy(page, sel, proxyList[rand.Intn(len(proxyList))], fields); err == nil {
			return result, nil
		}
	}

	return [][]string{}, err
}
