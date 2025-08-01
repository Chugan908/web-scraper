package tests

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
	"web-scraper/cmd"
	"web-scraper/proxies"
)

func init() {
	proxies.SetProxies()
}

func TestBasicShowCode_SuccessCase(t *testing.T) {
	command := cmd.RootCmd

	output := new(bytes.Buffer)

	command.SetOut(output)
	command.SetArgs([]string{"show-code", "-p", "https://www.scrapingcourse.com/infinite-scrolling"})

	if err := command.Execute(); err != nil {
		require.NoError(t, err)
	}

	require.NotEmpty(t, output)
	require.NotContains(t, output.String(), "Couldn't fetch the code")
	require.NotContains(t, output.String(), "Couldn't bypass the protection")
}

func TestBasicShowCode_FailCase(t *testing.T) {
	tests := []struct {
		name        string
		page        string
		ExpectedErr string
	}{
		{
			name:        "Web-Scraping a dead page",
			page:        "fdsfsd",
			ExpectedErr: "Provided page is dead",
		},
		{
			name:        "Web-Scraping a protected page",
			page:        "https://www.scrapingcourse.com/antibot-challenge",
			ExpectedErr: "Couldn't bypass the protection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := cmd.RootCmd

			output := new(bytes.Buffer)

			command.SetOut(output)
			command.SetArgs([]string{"show-code", "-p", tt.page})

			if err := command.Execute(); err != nil {
				require.NoError(t, err)
			}

			require.NotEmpty(t, output)
			require.Contains(t, output.String(), tt.ExpectedErr)
		})
	}
}

func TestScrollingShowCode_SuccessCase(t *testing.T) {
	command := cmd.RootCmd

	output := new(bytes.Buffer)

	command.SetOut(output)
	command.SetArgs([]string{"show-code", "-p", "https://www.scrapingcourse.com/infinite-scrolling", "--scroll-times", "5"})

	if err := command.Execute(); err != nil {
		require.NoError(t, err)
	}

	require.NotEmpty(t, output)
	require.NotContains(t, output.String(), "Couldn't fetch the code")
	require.NotContains(t, output.String(), "Couldn't bypass the protection")
}

func TestLoadingShowCode_SuccessCase(t *testing.T) {
	command := cmd.RootCmd

	output := new(bytes.Buffer)

	command.SetOut(output)
	command.SetArgs([]string{"show-code", "-p", "https://www.scrapingcourse.com/button-click", "--scroll-times", "5", "-l", "#load-more-btn"})

	if err := command.Execute(); err != nil {
		require.NoError(t, err)
	}

	require.NotEmpty(t, output)
	require.NotContains(t, output.String(), "Couldn't fetch the code")
	require.NotContains(t, output.String(), "Couldn't bypass the protection")
}

func TestLoadingShowCode_FailCase(t *testing.T) {
	command := cmd.RootCmd

	output := new(bytes.Buffer)

	command.SetOut(output)
	command.SetArgs([]string{"show-code", "-p", "https://www.scrapingcourse.com/button-click", "--scroll-times", "5", "-l", "#load-more-btnadasdad"})

	if err := command.Execute(); err != nil {
		require.NoError(t, err)
	}

	require.NotEmpty(t, output)
	require.Contains(t, output.String(), "provided selector was not found")
}
