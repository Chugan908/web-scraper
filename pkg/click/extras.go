package click

import (
	"context"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"time"
	"web-scraper/pkg/extract"
	"web-scraper/pkg/showcode"
)

func ClickOnWithProxy(page, sel, capture, proxy string) (string, error) {
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

	var tasks chromedp.Tasks

	if capture != "" {
		tasks = append(tasks, chromedp.OuterHTML(capture, &htmlContent, chromedp.ByQuery))
	} else {
		tasks = append(tasks, showcode.GetFullHTML(&htmlContent))
	}

	err := chromedp.Run(ctx,
		chromedp.Navigate(page),
		chromedp.Click(sel),
		chromedp.WaitVisible("html"),
		tasks,
	)
	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

func DoClickOnExtract(ctx context.Context, page, sel, extract_sel string, fields []string) (result [][]string, err error) {
	err = chromedp.Run(ctx,
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
		chromedp.Sleep(1*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			result, err = extract.DoExtract(ctx, page, extract_sel, fields, true)
			if err != nil {
				return err
			}

			return nil
		}),
	)
	if err != nil {
		return [][]string{}, err
	}

	return result, nil
}

func DoClickOnExtractProxy(page, sel, extract_sel, proxy string, fields []string) (result [][]string, err error) {
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

	return DoClickOnExtract(ctx, page, sel, extract_sel, fields)
}
