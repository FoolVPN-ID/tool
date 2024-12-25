package library

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func YoutubeCDN(httpClient http.Client) runnerResultStruct {
	result := runnerResultStruct{
		Name: "Youtube CDN",
	}

	start := time.Now()
	req, _ := http.NewRequest("GET", "https://redirector.googlevideo.com/report_mapping", nil)
	res, err := httpClient.Do(req)
	result.Delay = int(time.Since(start).Milliseconds())
	defer res.Body.Close()

	if res.StatusCode == 200 && err == nil {
		buf := new(strings.Builder)
		io.Copy(buf, res.Body)

		content := strings.Split(buf.String(), "\n")[0]
		regionPattern := regexp.MustCompile(`=>\s((\w+-(\w{3}))|(\w{3}))`)
		matchResults := strings.Split(regionPattern.FindString(content), " ")
		region := strings.ToUpper(matchResults[len(matchResults)-1])

		result.Region = region

		return result
	}

	result.Error = err
	return result
}
