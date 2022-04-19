package main

import "go_app/web"

const HttpPort = "8887"

func main() {
	web.RunWebServer(HttpPort)
}
