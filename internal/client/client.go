package client

import (
	"fmt"
	"io"
	"net/http"
)

func Run(url string) error {
	resp, err := http.Get(url)
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
