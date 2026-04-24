package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Unparalleled-Calvin/gserver/internal/client"
	"github.com/Unparalleled-Calvin/gserver/internal/settings"
)

func splitInput(line string) (string, string) {
	parts := strings.SplitN(line, " ", 2)
	uri := strings.TrimSpace(parts[0])
	form := ""
	if len(parts) == 2 {
		form = strings.TrimSpace(parts[1])
	}
	return uri, form
}

func basicUrl(url string) string {
	if url[0] == ':' {
		url = "http://localhost" + url
	}
	return url
}

func main() {
	settings.Load()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter URI and JSON form data:")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		uri, form := splitInput(line)
		url := basicUrl(settings.ServerAddr + uri)
		if err := client.Run(url, form); err != nil {
			log.Fatal(err)
			break
		}
	}
}
