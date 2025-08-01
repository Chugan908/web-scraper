package showcode

import (
	"context"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"time"
)

func ScrollPage(times int, wait time.Duration) chromedp.Tasks {
	// initialize tasks for chromedp instance
	var tasks chromedp.Tasks
	for i := 0; i < times; i++ {
		tasks = append(tasks,
			// scroll to bottom
			chromedp.Evaluate("window.scrollTo(0, document.body.scrollHeight)", nil),

			// wait a bit for new data to load
			chromedp.Sleep(wait),
		)
	}

	return tasks
}

func ScrollLoadMorePage(times int, loadMore string, wait time.Duration) chromedp.Tasks {
	// initialize tasks for chromedp instance
	var tasks chromedp.Tasks
	for i := 0; i < times; i++ {
		tasks = append(tasks,
			// scroll to bottom
			chromedp.Evaluate("window.scrollTo(0, document.body.scrollHeight)", nil),

			// check for button to load more content
			chromedp.ActionFunc(func(ctx context.Context) error {
				var nodes []*cdp.Node

				chromedp.Nodes(loadMore, &nodes, chromedp.AtLeast(0)).Do(ctx)

				chromedp.Click(loadMore).Do(ctx)

				return nil
			}),

			// wait a bit for new data to load
			chromedp.Sleep(wait),
		)
	}

	return tasks
}

func GetFullHTML(htmlContent *string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// select the root node on the page
		rootNode, err := dom.GetDocument().Do(ctx)
		if err != nil {
			return err
		}

		*htmlContent, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)

		return err
	})
}

func GetScrollingContent(ctx context.Context, times int, page string) (string, error) {
	var htmlContent string

	return htmlContent, chromedp.Run(ctx,
		// Navigate to the page
		chromedp.Navigate(page),

		// Perform scrolls
		ScrollPage(times, 1*time.Second),

		// Get the HTML after scrolling
		GetFullHTML(&htmlContent),
	)
}

func GetLoadingContent(ctx context.Context, times int, page, loadMore string) (string, error) {
	var htmlContent string

	return htmlContent, chromedp.Run(ctx,
		// Navigate to the page
		chromedp.Navigate(page),

		// check whether selector exists
		chromedp.ActionFunc(func(ctx context.Context) error {

			chromedp.Evaluate("window.scrollTo(0, document.body.scrollHeight)", nil)

			var nodes []*cdp.Node
			chromedp.Nodes(loadMore, &nodes, chromedp.AtLeast(0)).Do(ctx)

			if len(nodes) == 0 {
				return fmt.Errorf("provided selector was not found")
			}

			return nil
		}),

		// Perform scrolls
		ScrollLoadMorePage(times, loadMore, 1*time.Second),

		// Get the HTML after scrolling
		GetFullHTML(&htmlContent),
	)
}

func ScrapeWithProxy(times int, loadMore, page, proxy string) (string, error) {
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

	// initialize a controllable Chrome instance
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	if loadMore != "" {
		return GetLoadingContent(ctx, times, page, loadMore)
	}

	return GetScrollingContent(ctx, times, page)
}
