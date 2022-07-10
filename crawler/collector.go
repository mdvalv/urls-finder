package crawler

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

func NewCollector(options *Options) (*colly.Collector, error) {
	collector := colly.NewCollector(
		colly.Async(true),
	)

	if options.Local {
		domains, err := getAllowedDomains(options.BaseUrl)
		if err != nil {
			return nil, err
		}
		collector.URLFilters = domains
	}

	collector.MaxDepth = options.Depth

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: options.Threads,
	})

	collector.OnHTML("[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absLink := e.Request.AbsoluteURL(link)
		e.Request.Visit(absLink)
	})
	collector.OnHTML("[src]", func(e *colly.HTMLElement) {
		link := e.Attr("src")
		absLink := e.Request.AbsoluteURL(link)
		e.Request.Visit(absLink)
	})

	collector.OnRequest(func(r *colly.Request) {
		absLink := r.URL.String()
		if options.Eager {
			fmt.Println(absLink)
		}
		links.Store(absLink, true)
	})

	collector.OnResponse(func(r *colly.Response) {
		absLink := r.Request.URL.String()
		if strings.HasSuffix(absLink, ".js") {
			links := parseJSFIle(r.Body)
			for _, link := range links {
				r.Request.Visit(link)
			}
		}
	})

	return collector, nil
}
