package work

type Url string

type CrawlResponse struct {
	Urls []Url
	Body string
	Err  error
}
