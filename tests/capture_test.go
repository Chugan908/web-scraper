package tests

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
	"web-scraper/cmd"
)

func TestCapture_SuccessCase(t *testing.T) {
	command := cmd.RootCmd

	output := new(bytes.Buffer)

	command.SetOut(output)
	command.SetArgs([]string{"capture", "-p", "https://www.scrapingcourse.com/table-parsing", "-c", "#product-catalog"})

	if err := command.Execute(); err != nil {
		require.NoError(t, err)
	}

	require.NotEmpty(t, output)
	require.NotContains(t, output.String(), "a non-existent selector was provided")
	require.NotContains(t, output.String(), "Couldn't capture the provided selector")
	require.NotContains(t, output.String(), "No CSS selector was provided")
	require.NotContains(t, output.String(), "Provided page is dead")
}

func TestCapture_FailCase(t *testing.T) {
	tests := []struct {
		name        string
		page        string
		capture     string
		ExpectedErr string
	}{
		{
			name:        "Web-Scraping a dead page",
			page:        "fdsfsd",
			ExpectedErr: "Provided page is dead",
		},
		{
			name:        "Empty capture selector",
			page:        "https://www.scrapingcourse.com/table-parsing",
			capture:     "",
			ExpectedErr: "No CSS selector was provided",
		},
		{
			name:        "Non-existent selector",
			page:        "https://www.scrapingcourse.com/table-parsing",
			capture:     "dsfsfs",
			ExpectedErr: "a non-existent selector was provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := cmd.RootCmd

			output := new(bytes.Buffer)

			command.SetOut(output)
			command.SetArgs([]string{"capture", "-p", tt.page, "-c", tt.capture})

			if err := command.Execute(); err != nil {
				require.NoError(t, err)
			}

			require.NotEmpty(t, output)
			require.Contains(t, output.String(), tt.ExpectedErr)
		})
	}
}
