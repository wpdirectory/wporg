package wporg

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

const (
	wpListURL = "http://%s.svn.wordpress.org/"
)

var (
	// TODO: Update this to html parsing?
	regexList = regexp.MustCompile(`.+?\>(\S+?)\/\<`)
)

// GetList fetches a list of Plugins/Themes from the WordPress Directories
func (c *Client) GetList(dir string) ([]string, error) {
	var list []string

	// Prepare the URL
	URL := fmt.Sprintf(wpListURL, dir)

	// Make the Request
	resp, err := c.getRequest(URL)
	if err != nil {
		return list, err
	}

	// Drain body and check Close error
	defer drainAndClose(resp.Body, &err)
	bytes, err := ioutil.ReadAll(resp.Body)
	matches := regexList.FindAllStringSubmatch(string(bytes), -1)

	// Add all matches to extension list
	for _, match := range matches {
		list = append(list, match[1])
	}

	return list, err
}
