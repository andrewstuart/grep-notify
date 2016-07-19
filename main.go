package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/0xAX/notificator"
	"github.com/hpcloud/tail"
)

var (
	file = flag.String("f", "-", "the file to open")
	msg  = flag.String("msg", "Phrase Found", "the message notification title")
)

const (
	si = "/dev/stdin"
)

func init() {
	flag.Parse()
}

func main() {
	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "Must include a search phrase as the last argument.")
		os.Exit(1)
	}

	not := notificator.New(notificator.Options{AppName: "grep-notify"})

	if *file == "-" {
		*file = si
	}

	f, err := tail.TailFile(*file, tail.Config{Follow: *file != si, Pipe: *file == si})

	if err != nil {
		log.Fatal(err)
	}

	for l := range f.Lines {
		if strings.Contains(l.Text, flag.Args()[0]) {
			err = not.Push(*msg, fmt.Sprintf("%s: %s", time.Now(), l.Text), "", notificator.UR_NORMAL)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	if err = f.Err(); err != nil {
		log.Fatal(err)
	}
}
