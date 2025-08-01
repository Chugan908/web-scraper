package extract

import (
	"context"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"time"
)

func DoExtract(ctx context.Context, page, sel string, fields []string, click bool) ([][]string, error) {
	var Nodes []*cdp.Node

	var err error
	if !click {
		err = chromedp.Run(ctx,
			chromedp.Navigate(page),
			chromedp.Nodes(sel, &Nodes, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		)
	} else {
		err = chromedp.Run(ctx,
			chromedp.Nodes(sel, &Nodes, chromedp.ByQueryAll, chromedp.AtLeast(0)),
		)
	}

	if err != nil {
		return [][]string{}, err
	}

	if len(Nodes) == 0 {
		return [][]string{}, fmt.Errorf("error: provided extract selector was not found")
	}

	output := make([][]string, len(Nodes))

	var tasks chromedp.Tasks

	for i, node := range Nodes {
		output[i] = make([]string, len(fields))
		for j, field := range fields {
			tasks = append(tasks, chromedp.Text(field, &output[i][j], chromedp.ByQuery, chromedp.FromNode(node)))
		}
	}

	tasksCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = chromedp.Run(tasksCtx,
		tasks,
	)
	if err != nil {
		return [][]string{}, err
	}

	return output, err
}

func DoExtractProxy(page, sel, proxy string, fields []string) ([][]string, error) {
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

	return DoExtract(ctx, page, sel, fields, false)
}
