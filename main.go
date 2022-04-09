package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/odwrtw/transmission"
)

func main() {
	address, username, password, logpath := initFlags()
	if address == "" || username == "" || password == "" {
		flag.Usage()
		os.Exit(0)
	}

	logger := initLog(logpath)

	var config = transmission.Config{
		Address:  fmt.Sprintf("http://%s:9091/transmission/rpc", address),
		User:     username,
		Password: password,
	}

	t, err := transmission.New(config)
	if err != nil {
		logger.Fatalln("Error initializing!", err)
	}

	torrents, err := t.GetTorrents()
	if err != nil {
		logger.Fatalln("Error getting torrents!", err)
	}

	removed := 0
	stopped := 0
	for _, torrent := range torrents {
		switch torrent.Status {
		case transmission.StatusSeeding, transmission.StatusSeedPending:
			if err := torrent.Stop(); err == nil {
				stopped++
			} else {
				logger.Printf("Couldn't stop seeding %s: %v", torrent.Name, err)
			}
		case transmission.StatusStopped:
			if err := t.RemoveTorrents([]*transmission.Torrent{torrent}, true); err == nil {
				removed++
				doneAge := time.Since(time.Unix(int64(torrent.DoneDate), 0))
				logger.Printf("Removed %s (completed %s ago)", torrent.Name, doneAge.Round(time.Second).String())
			} else {
				logger.Printf("Couldn't remove %s: %v", torrent.Name, err)
			}
		}
	}

	logger.Printf("Removed %d, stopped %d (of %d total)\n", removed, stopped, len(torrents))
}

func initFlags() (string, string, string, string) {
	address := flag.String("address", "127.0.0.1", "IP or domain name address of Transmission, without port or protocol")
	username := flag.String("username", "", "Username of your Transmission user")
	password := flag.String("password", "", "Password of your Transmission user")
	logpath := flag.String("log", "/var/log/tr.log", "Path to log file")
	flag.Parse()

	return *address, *username, *password, *logpath
}

func initLog(path string) *log.Logger {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}

	return log.New(file, "TR: ", log.LstdFlags)
}
