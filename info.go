package wporg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"strconv"
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
		"request[fields][compatibility]=0", // Find a use for this, always seems to return []
		"request[fields][versions]=1",
		"request[fields][version]=1",
		"request[fields][screenshots]=1",
		"request[fields][stable_tag]=1",
		"request[fields][download_link]=1", // Not for themes
	}
)

// GetInfo fetches Plugin/Theme info from the API
func (c *Client) GetInfo(dir, name string) ([]byte, error) {
	var info *InfoResponse
	var bytes []byte

	// Main URL Components
	u := &url.URL{
		Scheme: "https",
		Host:   "api.wordpress.org",
		Path:   fmt.Sprintf("%s/info/1.1/", dir),
	}

	// Prepare Query Values
	// TODO: Add ability to overwrite default fields?
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
		return bytes, err
	}

	defer drainAndClose(resp.Body, &err)

	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return bytes, err
	}

	if string(bytes) == "false" {
		return nil, errors.New("No data returned")
	}

	err = json.Unmarshal(bytes, &info)
	if err != nil {
		return bytes, err
	}

	bytes, err = json.Marshal(info)
	if err != nil {
		return bytes, err
	}

	return bytes, err
}

// InfoResponse contains information about a Plugin or Theme
type InfoResponse struct {
	Name                   string     `json:"name"`
	Slug                   string     `json:"slug"`
	Version                string     `json:"version"`
	Author                 string     `json:"author"`
	AuthorProfile          string     `json:"author_profile"`
	Contributors           [][]string `json:"contributors"`
	Requires               string     `json:"requires"`
	Tested                 string     `json:"tested"`
	RequiresPHP            string     `json:"requires_php"`
	Rating                 int        `json:"rating"`
	Ratings                []Rating   `json:"ratings"`
	NumRatings             int        `json:"num_ratings"`
	SupportThreads         int        `json:"support_threads"`
	SupportThreadsResolved int        `json:"support_threads_resolved"`
	ActiveInstalls         int        `json:"active_installs"`
	Downloaded             int        `json:"downloaded"`
	LastUpdated            string     `json:"last_updated"`
	Added                  string     `json:"added"`
	Homepage               string     `json:"homepage"`
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
	//Compatibility          []interface{} `json:"-"`
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
		Version       json.Number `json:"version"`
		AuthorProfile interface{} `json:"author_profile"`
		RequiresPHP   interface{} `json:"requires_php"`
		Contributors  interface{} `json:"contributors"`
		Ratings       interface{} `json:"ratings"`
		NumRatings    interface{} `json:"num_ratings"`
		Screenshots   interface{} `json:"screenshots"`
		Tags          interface{} `json:"tags"`
		Versions      interface{} `json:"versions"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Set Version as string
	r.Version = aux.Version.String()
	/*
		switch v := aux.Version.(type) {
		case string:
			r.Version = v
		case int:
			r.Version = strconv.Itoa(v)
		default:
			r.Version = ""
		}
	*/

	// AuthorProfile can occasionally be a boolean (false)
	switch v := aux.AuthorProfile.(type) {
	case string:
		r.AuthorProfile = v
	default:
		r.AuthorProfile = ""
	}

	// RequiresPHP can occasionally be a boolean (false)
	switch v := aux.RequiresPHP.(type) {
	case string:
		r.RequiresPHP = v
	default:
		r.RequiresPHP = ""
	}

	// RequiresPHP can occasionally be a boolean (false)
	switch v := aux.RequiresPHP.(type) {
	case string:
		r.RequiresPHP = v
	default:
		r.RequiresPHP = ""
	}

	// Parse Contributors
	if aux.Contributors != nil && reflect.TypeOf(aux.Contributors).Kind() == reflect.Map {
		for k, v := range aux.Contributors.(map[string]interface{}) {
			contrib := []string{
				k, v.(string),
			}
			r.Contributors = append(r.Contributors, contrib)
		}
	}

	// Parse Ratings
	if reflect.TypeOf(aux.Ratings).Kind() == reflect.Map {
		for k, v := range aux.Ratings.(map[string]interface{}) {
			var num int
			var err error
			switch t := v.(type) {
			case float64:
				num = int(t)
			case string:
				num, err = strconv.Atoi(t)
				if err != nil {
					num = 0
				}
			default:
				num = 0
			}
			rating := Rating{
				Stars:  k,
				Number: num,
			}
			r.Ratings = append(r.Ratings, rating)
		}
	}

	// NumRatings can be a string "0" when zero
	switch v := aux.NumRatings.(type) {
	case int:
		r.NumRatings = v
	case string:
		num, err := strconv.Atoi(v)
		if err != nil {
			r.NumRatings = 0
		} else {
			r.NumRatings = num
		}
	default:
		r.NumRatings = 0
	}

	// Parse Screenshots
	if reflect.TypeOf(aux.Screenshots).Kind() == reflect.Map {
		for _, v := range aux.Screenshots.(map[string]interface{}) {
			s := v.(map[string]interface{})
			screenshot := Screenshot{
				Src:     s["src"].(string),
				Caption: s["caption"].(string),
			}
			r.Screenshots = append(r.Screenshots, screenshot)
		}
	}

	// Parse Tags
	if reflect.TypeOf(aux.Tags).Kind() == reflect.Map {
		for k, v := range aux.Tags.(map[string]interface{}) {
			tag := []string{
				k, v.(string),
			}
			r.Tags = append(r.Tags, tag)
		}
	}

	// Parse Versions
	if reflect.TypeOf(aux.Versions).Kind() == reflect.Map {
		for k, v := range aux.Versions.(map[string]interface{}) {
			version := []string{
				k, v.(string),
			}
			r.Versions = append(r.Versions, version)
		}
	}

	return nil
}
