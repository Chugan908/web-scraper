package click

import (
	"context"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"math/rand"
	"time"
	"web-scraper/pkg/showcode"
	"web-scraper/proxies"
)

func ClickOn(page, sel, capture string) (string, error) {
	// define the settings
	options := []chromedp.ExecAllocatorOption{
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.UserAgent(browser.Random()),
	}

	// specify a new context set up for NewContext
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 7*time.Second)
	defer cancel()

	// initialize a controllable Chrome instance
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var tasks chromedp.Tasks
	var htmlContent string

	if capture != "" {
		tasks = append(tasks, chromedp.OuterHTML(capture, &htmlContent))
	} else {
		tasks = append(tasks, showcode.GetFullHTML(&htmlContent))
	}

	err := chromedp.Run(ctx,
		chromedp.Navigate(page),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var nodes []*cdp.Node

			chromedp.Nodes(sel, &nodes, chromedp.ByQueryAll, chromedp.AtLeast(0)).Do(ctx)
			if len(nodes) == 0 {
				return fmt.Errorf("provided click-on selector was not found")
			}

			return nil
		}),
		chromedp.Click(sel),
		chromedp.WaitVisible("html"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if capture == "" {
				return nil
			}
			var nodes []*cdp.Node

			chromedp.Nodes(capture, &nodes, chromedp.ByQueryAll, chromedp.AtLeast(0)).Do(ctx)
			if len(nodes) == 0 {
				return fmt.Errorf("provided capture selector was not found")
			}

			return nil
		}),
		tasks,
	)
	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

func ClickOnProxy(page, sel, capture string) (htmlContent string, err error) {
	proxyList := proxies.GetProxies()

	for i := 0; i < 5; i++ {
		if htmlContent, err := ClickOnWithProxy(page, sel, capture, proxyList[rand.Intn(len(proxyList))]); err == nil {
			return htmlContent, nil
		}
	}

	return "", err
}

func ClickOnExtract(page, sel, extract string, fields []string) (result [][]string, err error) {
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

	return DoClickOnExtract(ctx, page, sel, extract, fields)
}

func ClickOnExtractProxy(page, sel, extract string, fields []string) (result [][]string, err error) {
	proxyList := proxies.GetProxies()

	for i := 0; i < 5; i++ {

		if result, err = DoClickOnExtractProxy(page, sel, extract, proxyList[rand.Intn(len(proxyList))], fields); err == nil {
			return result, nil
		}
	}

	return [][]string{}, err
}
