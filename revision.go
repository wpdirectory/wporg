package wporg

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
)

const (
	wpRevisionURL = "http://%s.trac.wordpress.org/log/?format=changelog&stop_rev=HEAD"
)

var (
	regexRevision = regexp.MustCompile(`\[(\d+?)\]`)
)

// GetRevision fetches the latest revision of the Plugins/Themes Directories
func (c *Client) GetRevision(dir string) (int, error) {
	var revision int

	// Prepare the URL
	URL := fmt.Sprintf(wpRevisionURL, dir)

	// Make the Request
	resp, err := c.getRequest(URL)
	if err != nil {
		return revision, err
	}

	// Drain body and check Close error
	defer drainAndClose(resp.Body, &err)
	bytes, err := ioutil.ReadAll(resp.Body)
	revs := regexRevision.FindAllStringSubmatch(string(bytes), 1)

	revision, err = strconv.Atoi(revs[0][1])
	if err != nil {
		return 0, err
	}

	return revision, err
}
