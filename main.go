package main

import (
	"flag"
	"github.com/catalinc/ipmon/lib"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	defaultNetConf  = "net.json"
	defaultMailConf = "mail.json"
	defaultInterval = 30
)

var (
	interval uint
	runOnce  bool
	netConf  string
	mailConf string
)

func init() {
	flag.UintVar(&interval, "interval", defaultInterval, "interval between checks (in seconds)")
	flag.BoolVar(&runOnce, "runOnce", false, "run once and stop")
	flag.StringVar(&netConf, "netConf", defaultNetConf, "network configuration file")
	flag.StringVar(&mailConf, "mailConf", defaultMailConf, "mail configuration file")
}

// run implements the network configuration check
func run() error {
	log.Println("Checking network configuration...")

	crtConf, err := lib.GetNetConfig()
	if err != nil {
		return err
	}

	if exists(netConf) {
		log.Println("Found previous configuration")

		prevConf, err := lib.LoadNetConfig(netConf)
		if err != nil {
			return err
		}

		if crtConf.IsChanged(prevConf) {
			log.Println("Network configuration changed")
			log.Println("Sending mail...")
			mc, err := lib.LoadMailConfig(mailConf)
			if err != nil {
				return err
			}
			diffs := lib.Report(crtConf, prevConf)
			err = lib.SendMailSSL(mc, "Network configuration changed on "+crtConf.Hostname, diffs)
			if err != nil {
				return err
			}
		} else {
			log.Println("No changes detected")
		}
	} else {
		log.Println("Previous configuration not found")
	}

	err = crtConf.Save(netConf)
	if err != nil {
		return err
	}
	log.Println("Current configuration saved")

	return nil
}

// exists reports whether the named file or directory exists
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func main() {
	flag.Parse()

	log.Println("Running...")

	if runOnce {
		if err := run(); err != nil {
			log.Fatalf("Error: %v\n", err)
		}
	} else {
		done := make(chan bool)
		quit := make(chan bool)
		trap := make(chan os.Signal, 1)

		ticker := time.NewTicker(time.Duration(interval) * time.Second)

		signal.Notify(trap, os.Interrupt)
		go func() {
			select {
			case sig := <-trap:
				log.Printf("Aborting... Got %v\n", sig)
				quit <- true
			}
		}()

		go func() {
			for {
				select {
				case <-ticker.C:
					if err := run(); err != nil {
						log.Printf("Error: %v", err)
					}
				case <-quit:
					ticker.Stop()
					done <- true
					return
				}
			}
		}()

		<-done
		log.Println("Bye")
	}
}
