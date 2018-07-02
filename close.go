package wporg

import (
	"io"
	"io/ioutil"
)

// checkClose is used to check the return from Close in a defer statement.
func checkClose(c io.Closer, err *error) {
	cErr := c.Close()
	if *err == nil {
		*err = cErr
	}
}

// drainAndClose discards all data from rd and closes it.
func drainAndClose(rd io.ReadCloser, err *error) {
	if rd == nil {
		return
	}
	_, _ = io.Copy(ioutil.Discard, rd)
	cErr := rd.Close()
	if err != nil && *err == nil {
		*err = cErr
	}
}
