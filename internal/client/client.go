package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Run(url, form string) error {
	var resp *http.Response
	var err error
	if form != "" {
		resp, err = http.Post(url, "application/json", strings.NewReader(form))
	} else {
		resp, err = http.Get(url)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Print(string(body))
	return nil
}
