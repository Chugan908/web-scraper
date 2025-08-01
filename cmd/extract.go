package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"web-scraper/pkg/extract"
)

var ExtractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Use this command to extract necessary data from a page",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		page, _ := cmd.Flags().GetString("page")

		r, err := http.Head(page)
		if err != nil || r.StatusCode >= 400 {
			fmt.Fprintln(cmd.OutOrStdout(), "Provided page is dead")
			return
		}

		sel, _ := cmd.Flags().GetString("selector")

		if sel == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "A selector wasn't provided")
			return
		}

		fields, _ := cmd.Flags().GetStringSlice("fields")

		if len(fields) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Fields weren't provided")
			return
		}

		output, err := extract.Extract(page, sel, fields)
		if err != nil {
			if err.Error() == "error: provided extract selector was not found" {
				fmt.Fprintln(cmd.OutOrStdout(), "err: non-existent selector was provided")
				return
			} else if errors.Is(err, context.DeadlineExceeded) {
				fmt.Fprintln(cmd.OutOrStdout(), "err: non-existent fields were provided")
				return
			}
			output, err = extract.ExtractProxy(page, sel, fields)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "Couldn't extract info from the page", err)
				return
			}
		}

		fmt.Fprintln(cmd.OutOrStdout(), output)
	},
}

func init() {
	RootCmd.AddCommand(ExtractCmd)

	ExtractCmd.Flags().StringP("page", "p", "", "provide a page url")

	ExtractCmd.Flags().StringP("selector", "s", "", "provide a selector you want work with")

	ExtractCmd.Flags().StringSliceP("fields", "f", []string{}, "provide the fields that you want extract from the chosen selector")
}
