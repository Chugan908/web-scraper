package capture

import (
	"context"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/chromedp/chromedp"
	"time"
)

func DoCaptureWithProxy(page, sel, proxy string) (string, error) {
	// define the settings
	options := []chromedp.ExecAllocatorOption{
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.UserAgent(browser.Random()),
		chromedp.ProxyServer(proxy),
	}

	// specify a new context set up for NewContext
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 7*time.Second)
	defer cancel()

	// initialize a controllable Chrome instance
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var htmlContent string

	err := chromedp.Run(ctx,
		chromedp.Navigate(page),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		return "", err
	}

	return htmlContent, nil
}
