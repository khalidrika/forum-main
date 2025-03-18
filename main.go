package main

import "forum/server"

func main() {
	if !server.Initialise() {
		return
	}
	server.StartServer()
	server.DB.Close()
}
