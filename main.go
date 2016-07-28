package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/0xAX/notificator"
	"github.com/hpcloud/tail"
)

var (
	msg = flag.String("msg", "Phrase Found", "the message notification title")
	num = flag.Int64("n", 1, "the number of entries to alert before exiting. <1 indicates all matches (grep-notify must be sent SIGINT (ctrl-c))")
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

	var fPath string

	if len(flag.Args()) > 1 {
		if len(flag.Args()) > 2 {
			fmt.Fprintln(os.Stderr, "Warning: only one filename can be used currently")
		}
		fPath = flag.Args()[1]
	}

	not := notificator.New(notificator.Options{AppName: "grep-notify"})

	if fPath == "-" || fPath == "" {
		fPath = si
	}

	isPipe := false

	st, err := os.Stat(fPath)
	if err != nil {
		log.Fatal(err)
	}

	modePipe := os.ModeNamedPipe | os.ModeCharDevice | os.ModeSocket

	if st.Mode()&modePipe != 0 {
		isPipe = true
	}

	f, err := tail.TailFile(fPath, tail.Config{Follow: !isPipe, Pipe: isPipe})

	if err != nil {
		log.Fatal(err)
	}

	var matches int64

	if *num < 1 {
		handleSigint()
	}

	for l := range f.Lines {
		if strings.Contains(l.Text, flag.Args()[0]) {
			matches++

			fmt.Fprintln(os.Stdout, l.Text)
			err = not.Push(*msg, fmt.Sprintf("%s: %s", time.Now(), l.Text), "", notificator.UR_NORMAL)
			if err != nil {
				log.Fatal(err)
			}

			if matches >= *num && *num > 0 {
				return
			}
		}
	}

	if err = f.Err(); err != nil {
		log.Fatal(err)
	}
}

func handleSigint() {
	ch := make(chan os.Signal)
	go signal.Notify(ch, os.Interrupt)
	go func() {
		select {
		case <-ch:
			os.Exit(0)
		}
	}()
}
