package main

import "simple-crud-golang/web"

const HttpPort = "8887"

func main() {
	web.RunWebServer(HttpPort)
}
