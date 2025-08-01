package tests

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"web-scraper/cmd"
)

func TestExtract_SuccessCase(t *testing.T) {
	command := cmd.RootCmd

	output := new(bytes.Buffer)

	command.SetOut(output)
	command.SetArgs([]string{"extract", "-p", "https://www.scrapingcourse.com/table-parsing", "-s", "#product-catalog > tbody > tr", "-f=" + ".product-name,.product-category,.product-price"})

	if err := command.Execute(); err != nil {
		require.NoError(t, err)
	}

	require.NotEmpty(t, output)
	require.NotContains(t, output.String(), "Couldn't extract info")
	require.NotContains(t, output.String(), "Provided page is dead")
	require.NotContains(t, output.String(), "A selector wasn't provided")
	require.NotContains(t, output.String(), "Fields weren't provided")
	require.NotContains(t, output.String(), "provided selector was not found")
	require.NotContains(t, output.String(), "non-existent fields were provided")
}

func TestExtract_FailCase(t *testing.T) {
	tests := []struct {
		name        string
		page        string
		selector    string
		fields      []string
		ExpectedErr string
	}{
		{
			name:        "Dead page",
			page:        "dsfdsf",
			selector:    "",
			fields:      []string{},
			ExpectedErr: "Provided page is dead",
		},
		{
			name:        "Empty selector",
			page:        "https://www.scrapingcourse.com/table-parsing",
			selector:    "",
			fields:      []string{},
			ExpectedErr: "A selector wasn't provided",
		},
		{
			name:        "Empty fields",
			page:        "https://www.scrapingcourse.com/table-parsing",
			selector:    "#product-catalog > tbody > tr",
			fields:      []string{},
			ExpectedErr: "Fields weren't provided",
		},
		{
			name:        "Bad selector",
			page:        "https://www.scrapingcourse.com/table-parsing",
			selector:    "sdfsdfsd",
			fields:      []string{".product-name", ".product-category", ".product-price"},
			ExpectedErr: "non-existent selector was provided",
		},
		{
			name:        "Bad fields",
			page:        "https://www.scrapingcourse.com/table-parsing",
			selector:    "#product-catalog > tbody > tr",
			fields:      []string{"dfdgfdgd"},
			ExpectedErr: "non-existent fields were provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := cmd.RootCmd
			output := new(bytes.Buffer)

			command.SetOut(output)
			flags := []string{"extract", "-p", tt.page, "-s", tt.selector}
			for i := 0; i < len(tt.fields); i++ {
				flags = append(flags, "-f")
				flags = append(flags, tt.fields[i])
			}

			command.SetArgs(flags)

			if err := command.Execute(); err != nil {
				require.NoError(t, err)
			}

			fmt.Println(output.String())
			require.NotEmpty(t, output)
			require.Contains(t, output.String(), tt.ExpectedErr)
		})
	}

}
