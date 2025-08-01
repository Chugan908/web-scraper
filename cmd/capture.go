package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"web-scraper/pkg/capture"
)

// TODO: add capture-all flag

var CaptureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Use this command to get a subsection of the HTML of the current page using the CSS selector you specify",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		page, _ := cmd.Flags().GetString("page")

		r, err := http.Head(page)
		if err != nil || r.StatusCode >= 400 {
			fmt.Fprintln(cmd.OutOrStdout(), "Provided page is dead")
			return
		}

		sel, _ := cmd.Flags().GetString("capture")

		if sel == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "No CSS selector was provided")
			return
		}

		output, err := capture.Capture(page, sel)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fmt.Fprintln(cmd.OutOrStdout(), "err: a non-existent selector was provided")
				return
			}
			output, err = capture.CaptureWithProxy(page, sel)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "Couldn't capture the provided selector", err)
				return
			}
		}

		fmt.Fprintln(cmd.OutOrStdout(), output)
	},
}

func init() {
	RootCmd.AddCommand(CaptureCmd)

	// Flags
	CaptureCmd.PersistentFlags().StringP("page", "p", "", "from which page you want to get a subsection")

	CaptureCmd.PersistentFlags().StringP("capture", "c", "", "capture an HTML subsection of the provided page")
}
