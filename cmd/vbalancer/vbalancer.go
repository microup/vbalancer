package main

import "vbalancer/internal/app"

func main() {
	var isAppStart = make(chan bool)
	defer close(isAppStart)
	app.Run(isAppStart)
}
