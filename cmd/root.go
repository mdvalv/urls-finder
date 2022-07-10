/*
Copyright © 2022 Marina Valverde mdval.eh@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mdvalv/urls-finder/crawler"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "urls-finder",
		Short: "A crawler for finding URLs in HTML and JavaScript",
		Long:  `A crawler for finding URLs in HTML pages and JavaScript scripts.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			options, err := parseOptions(cmd)
			if err != nil {
				return err
			}
			c, err := crawler.NewCrawler(options)
			if err != nil {
				return err
			}

			links := c.Crawl()
			if !options.Eager {
				sort.Strings(links)
				for _, link := range links {
					fmt.Println(link)
				}
			}

			if options.Output != "" {
				f, err := os.Create(options.Output)
				if err != nil {
					return err
				}
				defer f.Close()
			
				for _, link := range links {
					_, err := f.WriteString(fmt.Sprintf("%s\n", link))
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func parseOptions(cmd *cobra.Command) (*crawler.Options, error) {
	options := crawler.Options{}
	var err error

	options.BaseUrl, err = cmd.Flags().GetString("url")
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(options.BaseUrl, "http") {
		options.BaseUrl = fmt.Sprintf("http://%s", options.BaseUrl)
	}

	options.Threads, err = cmd.Flags().GetInt("threads")
	if err != nil {
		return nil, err
	}

	options.Eager, err = cmd.Flags().GetBool("eager")
	if err != nil {
		return nil, err
	}

	options.Local, err = cmd.Flags().GetBool("local")
	if err != nil {
		return nil, err
	}

	options.Output, err = cmd.Flags().GetString("output")
	if err != nil {
		return nil, err
	}

	options.Depth, err = cmd.Flags().GetInt("depth")
	if err != nil {
		return nil, err
	}

	return &options, nil
}

func init() {
	rootCmd.Flags().BoolP("eager", "e", false, "eager mode: show results as they're found")
	rootCmd.Flags().BoolP("local", "l", false, "list local links only")
	rootCmd.Flags().StringP("output", "o", "", "output file")
	rootCmd.Flags().IntP("depth", "d", 0, "max recursion depth of visited URLs (default 0 = infinite)")
	rootCmd.Flags().IntP("threads", "t", 1, "number of concurrent threads")
	rootCmd.Flags().StringP("url", "u", "", "the target URL")

	rootCmd.MarkFlagRequired("url")
}
