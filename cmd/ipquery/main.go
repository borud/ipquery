package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
)

// Options contains the command line options
type Options struct {
	ListenAddr string `short:"w" long:"web-listen-addr" default:":8000" description:"Web server listen address"`
	Path       string `short:"p" long:"path" default:"/whatismyip" description:"Path for service"`
	Header     string `short:"r" long:"header-name" default:"X-Real-IP" description:"Header from which to fetch IP address"`
}

var opt Options
var parser = flags.NewParser(&opt, flags.Default)

func main() {
	_, err := parser.Parse()
	if err != nil {
		log.Printf(">>> %v", err)
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	log.Printf("Config:")
	log.Printf("  ListenAddr : '%s'", opt.ListenAddr)
	log.Printf("  Path       : '%s'", opt.Path)
	log.Printf("  Header     : '%s'", opt.Header)

	http.HandleFunc(opt.Path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")

		if opt.Header != "" {
			ip := r.Header.Get(opt.Header)
			fmt.Fprintf(w, "%s", ip)
			log.Printf("ip=%s", ip)
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// Don't return anything to the client
			return
		}
		fmt.Fprintf(w, "%s", ip)
		log.Printf("ip=%s", ip)
	})
	http.ListenAndServe(opt.ListenAddr, nil)
}
