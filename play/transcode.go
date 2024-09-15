//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"time"

	"github.com/clicklord/lms/log"

	"github.com/clicklord/lms/misc"
)

func main() {
	ss := flag.String("ss", "", "")
	t := flag.String("t", "", "")
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("wrong argument count")
	}
	r, err := misc.Transcode(flag.Arg(0), *ss, *t)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		buf := bufio.NewWriterSize(os.Stdout, 1234)
		n, err := io.Copy(buf, r)
		log.Print("copied", n, "bytes")
		if err != nil {
			log.Print(err)
		}
	}()
	time.Sleep(time.Second)
	go r.Close()
	time.Sleep(time.Second)
}
