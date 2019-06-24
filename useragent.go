package minioanalytics

import (
	"regexp"
	"strings"
)

// User Agent regular expressions
var matches = map[string]string{
	"AndroidDownloadManager": `AndroidDownloadManager*`,
	"AntennaPod":             `AntennaPod*`,
	"BeyondPod":              `BeyondPod*`,
	"iPad":                   `iPad*`,
	"iPhone":                 `iPhone*`,
	"iPod":                   `iPod*`,
	"CastBox":                `CastBox*`,
	"Castro":                 `Castro*`,
	"Chrome":                 `Chrome*`,
	"Overcast":               `Overcast*`,
	"PodcastAddict":          `PodcastAddict*`,
	"Podcasts":               `Podcasts\/*`,
	"Safari":                 `Safari*`,
	"Spotify":                `Spotify*`,
	"Swoot":                  `Swoot*`,
	"watchOS":                `watchOS*`,
	"WindVane":               `WindVane*`,
}

// UserAgentMatcher tries to match user agent strings to a simple description
type UserAgentMatcher struct {
	matches map[string]*regexp.Regexp
}

// NewUserAgentMater retruns a new UserAgentMatcher with compiled regular expressions for matching
func NewUserAgentMater() (*UserAgentMatcher, error) {
	uam := &UserAgentMatcher{}
	return uam, uam.init()
}

func (uam *UserAgentMatcher) init() error {
	m := map[string]*regexp.Regexp{}

	for desc, pat := range matches {
		m[desc] = regexp.MustCompile(pat)
	}

	uam.matches = m
	return nil

}

// Match tries to match a given user agent
func (uam *UserAgentMatcher) Match(ua string) string {
	for desc, regexp := range uam.matches {
		if regexp.MatchString(ua) {
			return desc
		}
	}

	ua = strings.Split(ua, "/")[0]

	return ua
}

// MatchMap tries to match the map keys to their matched user agent descriptions
func (uam *UserAgentMatcher) MatchMap(uas map[string]int) map[string]int {

	res := map[string]int{}
	for ua, count := range uas {
		uaMatched := uam.Match(ua)
		if found, ok := res[uaMatched]; ok {
			res[uaMatched] = found + count
		} else {
			res[uaMatched] = count
		}
	}
	return res
}
