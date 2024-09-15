//go:build ignore
// +build ignore

package main

import (
	"flag"

	"github.com/anacrolix/ffprobe"
	"github.com/clicklord/lms/log"
)

func main() {
	flag.Parse()
	for _, path := range flag.Args() {
		i, err := ffprobe.Probe(path)
		log.Printf("%#v %#v", i, err)
	}
}
