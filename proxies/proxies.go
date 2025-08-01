package proxies

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Proxy struct {
	Addrs []string
}

func GetProxies() []string {
	rootPath, _ := os.Getwd()

	// Navigate up if needed
	if strings.HasSuffix(rootPath, "\\tests") {
		rootPath = filepath.Dir(rootPath)
	}

	proxyPath := filepath.Join(rootPath, "proxies", "storage", "proxies.json")

	f, _ := os.Open(proxyPath)
	data, _ := io.ReadAll(f)

	var proxies Proxy

	json.Unmarshal(data, &proxies)

	return proxies.Addrs
}

func SetProxies() {
	url := "https://sslproxies.org/"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Load HTML into Goquery
	content, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	// Parsing the Proxy List
	proxies := parseproxies(content)

	rootPath, _ := os.Getwd()

	// Navigate up if needed
	if strings.HasSuffix(rootPath, "\\tests") {
		rootPath = filepath.Dir(rootPath)
	}

	proxyPath := filepath.Join(rootPath, "proxies", "storage", "proxies.json")

	// writing proxies to the file
	f, err := os.Create(proxyPath)
	defer f.Close()
	res, _ := json.Marshal(proxies)
	f.Write(res)
}

func parseproxies(content *goquery.Document) Proxy {
	var proxies []string

	content.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		ip := s.Find("td:nth-child(1)").Text()
		port := s.Find("td:nth-child(2)").Text()

		if net.ParseIP(ip) == nil {
			return
		}

		proxies = append(proxies, "http://"+ip+":"+port)
	})

	return Proxy{proxies}
}
