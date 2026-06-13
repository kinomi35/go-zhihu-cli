package config

const ChromeVersion = "145"

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + ChromeVersion + ".0.0.0 Safari/537.36"

func BrowserHeaders(baseURL string) map[string]string {
	return map[string]string{
		"User-Agent":         UserAgent,
		"Accept":             "application/json, text/plain, */*",
		"Accept-Language":    "en-US,en;q=0.9",
		"Referer":            baseURL + "/",
		"sec-ch-ua":          `"Not:A-Brand";v="99", "Google Chrome";v="` + ChromeVersion + `", "Chromium";v="` + ChromeVersion + `"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"Windows"`,
	}
}
