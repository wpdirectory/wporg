package wporg

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

const (
	wpChangelogURL = "https://%s.trac.wordpress.org/log/?verbose=on&mode=follow_copy&format=changelog&rev=%d&limit=%d"
)

var (
	// TODO: Update this to html parsing?
	regexChangelog = regexp.MustCompile(`[a-z]+-[a-z]+\/`)
)

// GetChangeLog fetches a list of updated Plugins/Themes from between the provided revisions
func (c *Client) GetChangeLog(dir string, current, latest int) ([]string, error) {
	var list []string
	log.Printf("Current: %d Latest: %d\n", current, latest)

	for current < latest {
		URL := fmt.Sprintf(wpChangelogURL, dir, current, 100)
		log.Printf("URL: %s\n", URL)
		items, err := c.doChangeLog(URL)
		if err != nil {
			return list, err
		}
		list = append(list, items...)
		current += 100
	}

	// We are less than 100 updates behind, make one request
	URL := fmt.Sprintf(wpChangelogURL, dir, latest, 100)
	log.Printf("URL: %s\n", URL)
	items, err := c.doChangeLog(URL)
	if err != nil {
		return list, err
	}
	list = append(list, items...)

	return list, err
}

func (c *Client) doChangeLog(URL string) ([]string, error) {
	var list []string

	// Make the Request
	resp, err := c.getRequest(URL)
	if err != nil {
		return list, err
	}

	// Drain body and check Close error
	defer drainAndClose(resp.Body, &err)
	bytes, err := ioutil.ReadAll(resp.Body)
	matches := regexChangelog.FindAllString(string(bytes), -1)

	found := make(map[string]bool)
	// Get the desired substring match and remove duplicates
	for _, match := range matches {
		match = strings.TrimRight(match, "/")
		if !found[match] {
			found[match] = true
			list = append(list, match)
		}
	}

	return list, err
}
