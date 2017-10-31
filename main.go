package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

	if t, err := transmission.New(config); err != nil {
		logger.Fatalln("Error initializing!", err)
	}

	if torrents, err := t.GetTorrents(); err != nil {
		logger.Fatalln("Error getting torrents!", err)
	}

	finished := 0
	for _, torrent := range torrents {
		if torrent.Status == transmission.StatusSeeding {
			finished += 1
			t.RemoveTorrents([]*transmission.Torrent{torrent}, true)
			logger.Printf("Finished torrent %s (%s) has been removed\n", torrent.Comment, torrent.Name)
		}
	}

	if finished == 0 {
		logger.Printf("No finished torrents (of %d total)\n", len(torrents))
		os.Exit(0)
	}
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
	if file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		log.Fatalln("Failed to open log file", err)
	}

	return log.New(file, "TR: ", log.LstdFlags)
}
