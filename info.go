package wporg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

const (
	wpInfoURL         = "https://api.wordpress.org/%s/info/1.1/"
	timeFormatPlugins = "2006-01-02 3:04pm MST"
	timeFormatThemes  = "2006-01-02"
)

var (
	defaultFields = []string{
		"request[fields][sections]=1",
		"request[fields][description]=0",       // Defaults to off as too long
		"request[fields][short_description]=1", // Not for themes
		"request[fields][tested]=1",
		"request[fields][requires]=1",
		"request[fields][rating]=1",
		"request[fields][ratings]=1",
		"request[fields][downloaded]=1",
		"request[fields][active_installs]=1",
		"request[fields][last_updated]=1",
		"request[fields][homepage]=1",
		"request[fields][tags]=1",
		"request[fields][donate_link]=1",
		"request[fields][contributors]=1",
		"request[fields][compatibility]=1",
		"request[fields][versions]=1",
		"request[fields][version]=1",
		"request[fields][screenshots]=1",
		"request[fields][stable_tag]=1",
		"request[fields][download_link]=1", // Not for themes
	}
)

// GetInfo fetches Plugin/Theme info from the API
func (c *Client) GetInfo(dir, name string) (*InfoResponse, error) {
	var info *InfoResponse

	// Main URL Components
	u := &url.URL{
		Scheme: "https",
		Host:   "api.wordpress.org",
		Path:   fmt.Sprintf("%s/info/1.1/", dir),
	}

	// Prepare Query Values
	values := []string{
		fmt.Sprintf("action=%s_information", dir[:len(dir)-1]),
		"request[slug]=" + name,
	}
	values = append(values, defaultFields...)

	// Add Query Params to URL and return it as a string
	u.RawQuery = strings.Join(values, "&")
	URL := u.String()

	// Make the Request
	resp, err := c.getRequest(URL)
	if err != nil {
		return info, err
	}

	defer drainAndClose(resp.Body, &err)

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return info, err
	}

	//fmt.Printf("Resp: %s\n", string(bytes))

	err = json.Unmarshal(bytes, &info)
	if err != nil {
		return info, err
	}

	return info, err
}

// InfoResponse contains information about a Plugin or Theme
type InfoResponse struct {
	Name                   string        `json:"name"`
	Slug                   string        `json:"slug"`
	Version                string        `json:"version"`
	Author                 string        `json:"author"`
	AuthorProfile          string        `json:"author_profile"`
	Contributors           [][]string    `json:"contributors"`
	Requires               string        `json:"requires"`
	Tested                 string        `json:"tested"`
	RequiresPHP            string        `json:"requires_php"`
	Compatibility          []interface{} `json:"compatibility"`
	Rating                 int           `json:"rating"`
	Ratings                []Rating      `json:"ratings"`
	NumRatings             int           `json:"num_ratings"`
	SupportThreads         int           `json:"support_threads"`
	SupportThreadsResolved int           `json:"support_threads_resolved"`
	ActiveInstalls         int           `json:"active_installs"`
	Downloaded             int           `json:"downloaded"`
	LastUpdated            string        `json:"last_updated"`
	Added                  string        `json:"added"`
	Homepage               string        `json:"homepage"`
	Sections               struct {
		Description string `json:"description"`
		Faq         string `json:"faq"`
		Changelog   string `json:"changelog"`
		Screenshots string `json:"screenshots"`
	} `json:"sections"`
	ShortDescription string       `json:"short_description"`
	DownloadLink     string       `json:"download_link"`
	Screenshots      []Screenshot `json:"screenshots"`
	Tags             [][]string   `json:"tags"`
	StableTag        string       `json:"stable_tag"`
	Versions         [][]string   `json:"versions"`
	DonateLink       string       `json:"donate_link"`
}

// Rating contains information about ratings of a specific star level (0-5)
type Rating struct {
	Stars  string `json:"stars"`
	Number int    `json:"number"`
}

// Screenshot contains the source and caption of a screenshot
type Screenshot struct {
	Src     string `json:"src"`
	Caption string `json:"caption"`
}

// UnmarshalJSON provides custom JSON to struct decoding
// Handles parts of the API response not suited to Go structs
func (r *InfoResponse) UnmarshalJSON(data []byte) error {
	type Alias InfoResponse
	aux := &struct {
		Version      json.Number `json:"version"`
		Contributors interface{} `json:"contributors"`
		Ratings      interface{} `json:"ratings"`
		Screenshots  interface{} `json:"screenshots"`
		Tags         interface{} `json:"tags"`
		Versions     interface{} `json:"versions"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// Set Version as string
	r.Version = aux.Version.String()

	// Unmarshal JSON into interface to parse extra fields
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Parse Contributors
	var contribs [][]string
	for k, v := range raw["contributors"].(map[string]interface{}) {
		contrib := []string{
			k, v.(string),
		}
		contribs = append(contribs, contrib)
	}
	r.Contributors = contribs

	// Parse Ratings
	var ratings []Rating
	for k, v := range raw["ratings"].(map[string]interface{}) {
		rating := Rating{
			Stars:  k,
			Number: int(v.(float64)),
		}
		ratings = append(ratings, rating)
	}
	r.Ratings = ratings

	// Parse Screenshots
	var screenshots []Screenshot
	for _, v := range raw["screenshots"].(map[string]interface{}) {
		s := v.(map[string]interface{})
		screenshot := Screenshot{
			Src:     s["src"].(string),
			Caption: s["caption"].(string),
		}
		screenshots = append(screenshots, screenshot)
	}
	r.Screenshots = screenshots

	// Parse Tags
	var tags [][]string
	for k, v := range raw["tags"].(map[string]interface{}) {
		tag := []string{
			k, v.(string),
		}
		tags = append(tags, tag)
	}
	r.Tags = tags

	// Parse Versions
	var versions [][]string
	for k, v := range raw["versions"].(map[string]interface{}) {
		version := []string{
			k, v.(string),
		}
		versions = append(versions, version)
	}
	r.Versions = versions

	return nil
}
