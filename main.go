package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/odwrtw/transmission"
)

func main() {
	address, username, password := initFlags()
	if address == "" || username == "" || password == "" {
		flag.Usage()
		os.Exit(1)
	}

	var config = transmission.Config{
		Address:  fmt.Sprintf("http://%s:9091/transmission/rpc", address),
		User:     username,
		Password: password,
	}

	t, err := transmission.New(config)
	if err != nil {
		fmt.Println("Error initializing!", err)
		os.Exit(1)
	}

	torrents, err := t.GetTorrents()
	if err != nil {
		fmt.Println("Error getting torrents!", err)
		os.Exit(1)
	}

	if len(torrents) == 0 {
		fmt.Println("No finished torrents")
		os.Exit(0)
	}

	for _, torrent := range torrents {
		if torrent.IsFinished == true && torrent.Status == transmission.StatusSeeding {
			t.RemoveTorrents([]*transmission.Torrent{torrent}, true)
			fmt.Printf("Finished torrent %s (%s) has been removed\n", torrent.Comment, torrent.Name)
		}
	}
}

func initFlags() (string, string, string) {
	address := flag.String("address", "192.168.0.10", "IP or domain name address of Transmission, without port or protocol")
	username := flag.String("username", "", "Username of your Transmission user")
	password := flag.String("password", "", "Password of your Transmission user")
	flag.Parse()

	return *address, *username, *password
}
