package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"web-scraper/pkg/click"
)

// TODO: Add support of scrolling the page after click

var ClickOnCmd = &cobra.Command{
	Use:   "click-on",
	Short: "This command will trigger a click on a particular HTML element using a CSS selector provided",
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
			fmt.Fprintln(cmd.OutOrStdout(), "No CSS selector was provided")
			return
		}

		extract, _ := cmd.Flags().GetString("extract")

		capture, _ := cmd.Flags().GetString("capture")

		if extract != "" && capture != "" {
			fmt.Fprintln(cmd.OutOrStdout(), "Can't do both extract and capture after a click on")
		}

		fields, _ := cmd.Flags().GetStringSlice("fields")

		if extract != "" && len(fields) > 0 {
			output, err := click.ClickOnExtract(page, sel, extract, fields)
			if err != nil {
				if err.Error() == "provided click-on selector was not found" {
					fmt.Fprintln(cmd.OutOrStdout(), "err: provided click-on selector was not found")
					return
				} else if err.Error() == "error: provided extract selector was not found" {
					fmt.Fprintln(cmd.OutOrStdout(), "err: provided extract selector was not found")
					return
				}
				output, err = click.ClickOnExtractProxy(page, sel, extract, fields)
				if err != nil {
					fmt.Fprintln(cmd.OutOrStdout(), "Something went wrong:", err)
					return
				}
			}

			fmt.Fprintln(cmd.OutOrStdout(), output)
			return
		}

		output, err := click.ClickOn(page, sel, capture)
		if err != nil {
			if err.Error() == "provided click-on selector was not found" {
				fmt.Fprintln(cmd.OutOrStdout(), "err: provided click-on selector was not found")
				return
			} else if err.Error() == "provided capture selector was not found" {
				fmt.Fprintln(cmd.OutOrStdout(), "err: provided capture selector was not found")
				return
			}
			output, err = click.ClickOnProxy(page, sel, capture)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "Couldn't click on the provided element", err)
				return
			}
		}

		fmt.Fprintln(cmd.OutOrStdout(), output)
	},
}

func init() {
	RootCmd.AddCommand(ClickOnCmd)

	// Flags
	ClickOnCmd.PersistentFlags().StringP("page", "p", "", "provide a page you want to explore")

	ClickOnCmd.PersistentFlags().StringP("selector", "s", "", "provide a CSS selector you want to click on")

	ClickOnCmd.PersistentFlags().StringP("extract", "e", "", "provide a subsection from which you want to extract after a click-on")

	ClickOnCmd.PersistentFlags().StringSliceP("fields", "f", []string{}, "provide the fields you want to extract from the provided in extract subsection")

	ClickOnCmd.PersistentFlags().StringP("capture", "c", "", "capture an HTML subsection after a click-on")
}
