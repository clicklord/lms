//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"net/url"

	"github.com/clicklord/lms/log"
)

func main() {
	url_, err := url.Parse("[192:168:26:2::3]:1900")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", url_)
}
