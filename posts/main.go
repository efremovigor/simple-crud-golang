package main

import (
	_ "github.com/lib/pq"
)

const HttpPort = "8887"

func main() {
	runWebServer(HttpPort)
}
