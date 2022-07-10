package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sync"

	"github.com/gocolly/colly/v2"
)

var links sync.Map

type Crawler struct {
	collector *colly.Collector
	options   *Options
}

type Options struct {
	BaseUrl string
	Depth   int
	Eager   bool
	Local   bool
	Output  string
	Threads int
}

func NewCrawler(options *Options) (*Crawler, error) {
	collector, err := NewCollector(options)
	if err != nil {
		return nil, err
	}

	crawler := &Crawler{
		collector: collector,
		options:   options,
	}

	return crawler, nil
}

func (c Crawler) Crawl() []string {
	c.collector.Visit(c.options.BaseUrl)
	c.collector.Wait()

	var linksStr []string
	links.Range(func(key, value interface{}) bool {
		linksStr = append(linksStr, key.(string))
		return true
	})
	return linksStr
}

var LinkRegex = regexp.MustCompile(`https?://(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)`)

func parseJSFIle(body []byte) []string {
	bodyStr := string(body)
	return LinkRegex.FindAllString(bodyStr, -1)
}

func getAllowedDomains(u string) (domains []*regexp.Regexp, err error) {
	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	parsedUrl, err := url.Parse(u)
	if err != nil {
		return
	}

	domains = append(
		domains,
		regexp.MustCompile(fmt.Sprintf(`.*%s.*`, resp.Request.URL.Hostname())),
		regexp.MustCompile(fmt.Sprintf(`.*%s.*`, parsedUrl.Hostname())),
	)

	return
}
