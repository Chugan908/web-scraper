package showcode

import (
	"context"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"math/rand"
	"time"
	"web-scraper/proxies"
)

// TODO: implement headless scraping
// TODO: implement multiple page scraping
// TODO: implement page scrolling

func ShowBaseCode(page string) (htmlContent string, err error) {
	// initialize a new collector object
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	// set random user-agent
	extensions.RandomUserAgent(c)

	c.OnScraped(func(r *colly.Response) {
		htmlContent = string(r.Body)
	})

	if err := c.Visit(page); err == nil {
		return htmlContent, nil
	}

	// set a proxy list
	proxyList := proxies.GetProxies()
	if len(proxyList) == 0 {
		return "", fmt.Errorf("Couldn't get proxies")
	}

	// create Rotating Proxy Switcher
	rp, err := proxy.RoundRobinProxySwitcher(proxyList...)
	if err != nil {
		return "", err
	}

	// set a proxy
	c.SetProxyFunc(rp)

	for i := 0; i < 5; i++ {
		if err := c.Visit(page); err == nil {
			return htmlContent, nil
		}
	}

	return "", err
}

func ShowScrollingCode(times int, loadMore, page string) (htmlContent string, err error) {
	// define the settings
	options := []chromedp.ExecAllocatorOption{
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.UserAgent(browser.Random()),
	}

	// specify a new context set up for NewContext
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// initialize a controllable Chrome instance
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	if loadMore != "" {
		return GetLoadingContent(ctx, times, page, loadMore)
	}

	return GetScrollingContent(ctx, times, page)
}

func ShowScrollingCodeProxy(times int, loadMore, page string) (htmlContent string, err error) {
	// set a proxy list
	proxyList := proxies.GetProxies()

	for i := 0; i < 5; i++ {
		if htmlContent, err := ScrapeWithProxy(times, loadMore, page, proxyList[rand.Intn(len(proxyList))]); err == nil {
			return htmlContent, nil
		}
	}

	return "", err
}
