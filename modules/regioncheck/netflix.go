package regioncheck

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func Netflix(httpClient http.Client) runnerResultStruct {
	result := runnerResultStruct{
		Name: "Netflix",
	}

	var (
		testURLs = []string{"https://www.netflix.com/title/81280792", "https://www.netflix.com/title/70143836"}
		checkURL = "https://www.netflix.com"

		isTestPassed = false
	)

	for _, testURL := range testURLs {
		req, _ := http.NewRequest("GET", testURL, nil)
		req.Header.Set("host", checkURL)
		req.Header.Set("accept-language", "en-US,en;q=0.9")

		res, _ := httpClient.Do(req)
		if res.StatusCode == 200 {
			isTestPassed = true
			break
		}
	}

	if !isTestPassed {
		result.Error = errors.New("Forbidden")
		return result
	}

	start := time.Now()
	req, _ := http.NewRequest("GET", checkURL, nil)

	res, err := httpClient.Do(req)
	result.Delay = int(time.Since(start).Milliseconds())
	defer res.Body.Close()

	if res.StatusCode == 200 && err == nil {
		finalURL := res.Request.URL.String()
		countryCodePattern := regexp.MustCompile(`com\/(\w{2})`)
		matchResults := countryCodePattern.FindStringSubmatch(finalURL)

		if len(matchResults) == 2 {
			result.Country = strings.ToUpper(matchResults[1])
			result.Region = result.Country
		}

		return result
	}

	result.Error = err
	return result
}
