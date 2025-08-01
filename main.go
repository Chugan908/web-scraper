package main

import (
	"web-scraper/cmd"
	"web-scraper/proxies"
)

func main() {
	proxies.SetProxies()
	cmd.Execute()
}
