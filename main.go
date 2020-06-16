package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/otiai10/gosseract"
)

// A Shoe is a shoe.
type Shoe struct {
	name   string
	prices map[string]string
}

var c *gosseract.Client

func main() {
	c = gosseract.NewClient()
	defer c.Close()

	/*
		config.txt:
			First line: Shoe name
			Second line: Size
			Third line: Twilio account SID
			Fourth line: Twilio auth token
			Fifth line: Twilio phone number
			Sixth line: Destination phone number
	*/
	f, err := os.Open("config.txt")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	name := lines[0]
	size := lines[1]
	accountSid := lines[2]
	authToken := lines[3]
	srcNum := lines[4]
	dstNum := lines[5]
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("From", srcNum)
	msgData.Set("To", dstNum)
	msgData.Set("Body", "Price of "+name+" (Size "+size+"): "+getShoe(name).prices[size])
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
}
func getShoe(name string) Shoe {
	fmt.Println("Getting price data for shoe named '" + name + "'...")
	urls := map[string]string{
		"Yeezy 350 V2 Cinder NRF":      "https://www.goat.com/sneakers/yeezy-boost-350-v2-cinder-fy2903/available-sizes",
		"AJ1 Obsidian":                 "https://www.goat.com/sneakers/air-jordan-1-retro-high-og-obsidian-555088-140/available-sizes",
		"Yeezy 700 V3 Alvah":           "https://www.goat.com/sneakers/yeezy-700-v3-alvah-h67799/available-sizes",
		"SB Dunk Low Travis Scott":     "https://www.goat.com/sneakers/travis-scott-x-dunk-low-sb-ct5053-101/available-sizes",
		"AJ1 Travis Scott":             "https://www.goat.com/sneakers/travis-scott-x-air-jordan-1-retro-high-og-cd4487-100/available-sizes",
		"Yeezy 350 V2 Cloud White NRF": "https://www.goat.com/sneakers/yeezy-boost-350-v2-cloud-white-fw3042/available-sizes",
		"SB Dunk Low Chunky Dunky":     "https://www.goat.com/sneakers/ben-jerry-s-x-dunk-low-sb-chunky-dunky-cu3244-100/available-sizes",
		"Yeezy 350 V2 Black NRF":       "https://www.goat.com/sneakers/yeezy-boost-350-v2-black-yzy-350-v2-blk/available-sizes",
	}
	if _, ok := urls[name]; !ok {
		fmt.Println("Invalid shoe name. Please try again.")
		return Shoe{}
	}
	cmd := exec.Command("./gowitness", "single", "--url="+urls[name], "--resolution=500,3000")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Screen capture obtained.")

	c.SetImage(getFileName(urls[name]) + ".png")
	out, err := c.Text()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("OCR attempt finished. Begin error checking...")

	shoe := Shoe{name, map[string]string{}}

	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}
		if text[0] != '>' && len(text) > 3 {
			text = correct(text)
			splitText := strings.Split(text, " ")
			if len(splitText) == 2 && strings.Contains(splitText[0], "M") && strings.Contains(splitText[1], "$") {
				if _, ok := shoe.prices[splitText[0]]; !ok {
					shoe.prices[splitText[0]] = splitText[1]
				}
			}
		}
	}
	fmt.Println("Process complete.")
	return shoe
}

func getFileName(in string) string {
	return strings.Replace(strings.Replace(in, "/", "", -1), ":", "-", -1)
}

func correct(out string) string {
	errors := [][]string{
		{"AM", "4M"},
		{"om", "9M"},
		{"‘1M", "11M"},
		{"46M", "16M"},
		{" _ ", " "},
		{"43M", "13M"},
		{"mM", "7M"},
		{"OM", "9M"},
		{"44M", "14M"},
		{"_ ", " "},
		{"ASM", "4.5M"},
		{"am", "9M"},
		{"S750,", "$1,750"},
	}
	for _, v := range errors {
		out = strings.Replace(out, v[0], v[1], -1)
	}
	return out
}
