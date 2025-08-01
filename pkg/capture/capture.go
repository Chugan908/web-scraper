package capture

import (
	"context"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/chromedp/chromedp"
	"math/rand"
	"time"
	"web-scraper/proxies"
)

func Capture(page, sel string) (string, error) {
	var htmlContent string

	// define the settings
	options := []chromedp.ExecAllocatorOption{
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.UserAgent(browser.Random()),
	}

	// specify a new context set up for NewContext
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// initialize a controllable Chrome instance
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(page),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

func CaptureWithProxy(page, sel string) (htmlContent string, err error) {
	proxyList := proxies.GetProxies()

	for i := 0; i < 5; i++ {
		if htmlContent, err := DoCaptureWithProxy(page, sel, proxyList[rand.Intn(len(proxyList))]); err == nil {
			return htmlContent, nil
		}
	}

	return "", err
}
