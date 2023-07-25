package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func html_document(document string) string {
	return "<!DOCTYPE html>\n<html>\n" + document + "</html>"
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, html_document("<body>\nSome body\n</body>"))
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	SiteName string
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	User     string
	Password string
	Name     string
}

type SearchEngine struct {
	Name       string
	Allow      []string
	Disallow   []string
	CleanParam []string
}

var config Config

func LoadConfig(file string) error {
	configFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)

	if err != nil {
		return err
	}

	return nil
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	searchEngines := []SearchEngine{
		{
			Name: "*",
			Allow: []string{
				"/public/",
				"/news/",
			},
			Disallow: []string{
				"/profile/",
			},
			CleanParam: []string{
				"banner-test-tags",
			},
		},
	}
	w.Header().Set("Content-Type", "text/plain")
	for _, se := range searchEngines {
		fmt.Fprintf(w, "User-agent: %s\n\n", se.Name)
		for _, path := range se.Allow {
			fmt.Fprintf(w, "Allow: %s\n", path)
		}
		fmt.Fprintln(w)
		for _, path := range se.Disallow {
			fmt.Fprintf(w, "Disallow: %s\n", path)
		}
		fmt.Fprintln(w)
		if len(se.CleanParam) > 0 {
			for _, param := range se.CleanParam {
				fmt.Fprintf(w, "Clean-param: %s\n", param)
			}
		}
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Host: %s\n", config.SiteName)
		if len(config.SiteName) > 0 {
			fmt.Fprintf(w, "Sitemap: http://%s/sitemap.xml", config.SiteName)
		}
	}
}

func main() {
	err := LoadConfig("config.json")

	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	address := config.Server.Host + ":" + strconv.FormatUint(uint64(config.Server.Port), 10)

	http.HandleFunc("/", handler)
	http.HandleFunc("/robots.txt", robotsHandler)
	log.Fatal(http.ListenAndServe(address, nil))
}
