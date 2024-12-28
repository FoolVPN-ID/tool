package regioncheck

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func YoutubeCDN(httpClient http.Client) runnerResultStruct {
	result := runnerResultStruct{
		Name: "Youtube CDN",
		OK:   false,
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
		iataCodePattern := regexp.MustCompile(`=>\s((\w+-(\w{3}))|(\w{3}))`)
		matchResults := strings.Split(iataCodePattern.FindString(content), " ")
		iataCode := strings.ToUpper(matchResults[len(matchResults)-1])

		result.IATACode = iataCode
		result.Region = cases.Title(language.AmericanEnglish).String(GetRegionFromIATACode(iataCode))
		result.OK = true

		return result
	}

	result.Error = err
	return result
}
