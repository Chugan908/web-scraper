package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"web-scraper/pkg/showcode"
)

// TODO: implement proxy receiver
// TODO: implement cloudflare bypasser
// TODO: implement javascript rendering

var ShowCodeCmd = &cobra.Command{
	Use:   "show-code",
	Short: "Use this command to list the HTML code of the provided page",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		page, _ := cmd.Flags().GetString("page")

		// checking if page is alive
		r, err := http.Head(page)
		if err != nil || r.StatusCode >= 400 {
			fmt.Fprintln(cmd.OutOrStdout(), "Provided page is dead")
			return
		}

		scrollTimes, _ := cmd.Flags().GetInt("scroll-times")

		if scrollTimes > 0 {
			btn, _ := cmd.Flags().GetString("load-more")
			htmlContent, err := showcode.ShowScrollingCode(scrollTimes, btn, page)
			if err != nil {
				if err.Error() == "provided selector was not found" {
					fmt.Fprintln(cmd.OutOrStdout(), "Couldn't fetch the code", page, err)
					return
				}
				htmlContent, err = showcode.ShowScrollingCodeProxy(scrollTimes, btn, page)
				if err != nil {
					fmt.Fprintln(cmd.OutOrStdout(), "Couldn't fetch the code with proxy...", page, err)
					return
				} else if len(htmlContent) == 0 {
					fmt.Fprintln(cmd.OutOrStdout(), "Couldn't bypass the protection", page)
					return
				}
			}

			fmt.Fprintln(cmd.OutOrStdout(), htmlContent)

			return
		}

		htmlContent, err := showcode.ShowBaseCode(page)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "Couldn't fetch the code from", page, err)
			return
		} else if len(htmlContent) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Couldn't bypass the protection", page)
		}

		fmt.Fprintln(cmd.OutOrStdout(), htmlContent)
	},
}

func init() {
	RootCmd.AddCommand(ShowCodeCmd)

	// Flags
	ShowCodeCmd.PersistentFlags().StringP("page", "p", "", "show code of the page")

	ShowCodeCmd.PersistentFlags().Int("scroll-times", 0, "define the amount of times the page should be scrolled")

	ShowCodeCmd.PersistentFlags().StringP("load-more", "l", "", "provide the name of load-more element")
}
