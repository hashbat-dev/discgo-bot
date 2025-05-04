package helpers

import (
	"fmt"
	"net/http"
)

func DoesLinkWork(link string) (bool, error) {
	resp, err := http.Head(link)
	if err != nil {
		return false, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return true, nil
}
