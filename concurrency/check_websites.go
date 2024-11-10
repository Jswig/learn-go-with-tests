package concurrency

type WebsiteChecker func(string) bool

type result struct {
	url     string
	isValid bool
}

func CheckWebsites(websiteChecker WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultChannel := make(chan result)

	for _, url := range urls {
		go func(u string) {
			resultChannel <- result{url: u, isValid: websiteChecker(u)}
		}(url)
	}

	for range len(urls) {
		r := <-resultChannel
		results[r.url] = r.isValid
	}

	return results
}
