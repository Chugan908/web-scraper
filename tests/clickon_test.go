package tests

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"web-scraper/cmd"
)

func TestBasicClickOn_SuccessCase(t *testing.T) {
	command := cmd.RootCmd

	output := new(bytes.Buffer)

	command.SetOut(output)
	command.SetArgs([]string{"click-on", "-p", "https://www.scrapingcourse.com/pagination", "-s", "#pagination-container > div > div > a:nth-child(3)"})

	if err := command.Execute(); err != nil {
		require.NoError(t, err)
	}

	fmt.Println(output.String())
	require.NotEmpty(t, output)
	require.NotContains(t, output.String(), "No CSS selector was provided")
	require.NotContains(t, output.String(), "Can't do both extract and capture after a click on")
	require.NotContains(t, output.String(), "Something went wrong:")
	require.NotContains(t, output.String(), "Provided page is dead")
	require.NotContains(t, output.String(), "provided click-on selector was not found")
	require.NotContains(t, output.String(), "provided capture selector was not found")
	require.NotContains(t, output.String(), "Couldn't click on")
}

func TestCaptureClickOn_SuccessCase(t *testing.T) {
	command := cmd.RootCmd

	output := new(bytes.Buffer)

	command.SetOut(output)
	command.SetArgs([]string{"click-on", "-p", "https://www.scrapingcourse.com/pagination", "-s", "#pagination-container > div > div > a:nth-child(3)", "-c", "#challenge-info"})

	if err := command.Execute(); err != nil {
		require.NoError(t, err)
	}

	fmt.Println(output.String())
	require.NotEmpty(t, output)
	require.NotContains(t, output.String(), "No CSS selector was provided")
	require.NotContains(t, output.String(), "Can't do both extract and capture after a click on")
	require.NotContains(t, output.String(), "Something went wrong:")
	require.NotContains(t, output.String(), "Provided page is dead")
	require.NotContains(t, output.String(), "provided click-on selector was not found")
	require.NotContains(t, output.String(), "provided capture selector was not found")
	require.NotContains(t, output.String(), "Couldn't click on")
}

func TestExtractClickOn_SuccessCase(t *testing.T) {
	command := cmd.RootCmd

	output := new(bytes.Buffer)

	command.SetOut(output)
	command.SetArgs([]string{"click-on", "-p", "https://www.scrapingcourse.com", "-s", "#content-container > div > main > div > div:nth-child(10) > p:nth-child(3) > a", "-e", "#product-catalog > tbody > tr", "-f .product-name"})

	if err := command.Execute(); err != nil {
		require.NoError(t, err)
	}

	fmt.Println(output.String())
	require.NotEmpty(t, output)
	require.NotContains(t, output.String(), "No CSS selector was provided")
	require.NotContains(t, output.String(), "Can't do both extract and capture after a click on")
	require.NotContains(t, output.String(), "Something went wrong:")
	require.NotContains(t, output.String(), "Provided page is dead")
	require.NotContains(t, output.String(), "Couldn't click on")
	require.NotContains(t, output.String(), "provided click-on selector was not found")
	require.NotContains(t, output.String(), "provided extract selector was not found")
}
