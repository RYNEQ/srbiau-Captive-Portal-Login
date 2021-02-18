package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	loginurl  = "http://srbiaulogin.net/login"
	statusurl = "http://srbiaulogin.net/status"
)

func main() {
	username := flag.String("username", "", "Account Username")
	var passwd string
	flag.StringVar(&passwd, "password", "", "Account Password")
	flag.Parse()

	if *username == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -username USER [-password PASS]\n", os.Args[0])
		os.Exit(1)
	}

	if passwd == "" {
		fmt.Print("Enter Password: ")
		bytePassword, err := terminal.ReadPassword(0)
		fmt.Println()
		if err != nil {
			log.Fatalln("Error getting password")
			os.Exit(1)
		}
		passwd = string(bytePassword)
	}

	content, status := getPage(loginurl, nil)
	if status == 200 {
		r, _ := regexp.Compile(`hexMD5\('(\\[0-9]{3})'\s*\+\s*document\.login\.password\.value *\+ *'((\\[0-9]{3})+)'`)
		dc := r.FindStringSubmatch(content)
		if len(dc) == 4 {
			passwordTemplate, err := strconv.Unquote(`"` + fmt.Sprintf("%s%s%s", dc[1], passwd, dc[2]) + `"`)
			if err != nil {
				log.Fatalln("Error unquoting password")
				os.Exit(1)
			}

			hash := fmt.Sprintf("%x", md5.Sum([]byte(passwordTemplate)))
			//fmt.Println(hash)
			form := url.Values{}
			form.Add("username", *username)
			form.Add("password", hash)
			form.Add("popup", "true")
			form.Add("dst", "")
			content, status := getPage(loginurl, form)
			if status == 200 {
				//fmt.Println(content)
				if strings.Contains(content, "http://srbiaulogin.net/status") {
					fmt.Println("Login Succeed")
				} else if strings.Contains(content, "Simulation exceed") {
					log.Fatalln("Login Faild (Max session reached!!)")
					os.Exit(1)
				} else {
					log.Fatalln("Login Faild")
					os.Exit(1)
				}
			} else {
				log.Fatalf("Error getting login response (Status: %d)\n", status)
				os.Exit(1)
			}
		} else {
			log.Fatalln("Error extracting hash values")
			os.Exit(1)
		}
	} else {
		log.Fatalf("Error Getting login page (Status: %d)\n", status)
		os.Exit(1)
	}
}

func getPage(url string, form url.Values) (string, int) {
	hc := http.Client{}
	//	resp, err := http.Get("https://httpbin.org/get")
	var req *http.Request
	var err error
	if form == nil {
		req, err = http.NewRequest("GET", url, nil)
	} else {
		req, err = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	}
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:85.0) Gecko/20100101 Firefox/85.0")
	resp, err := hc.Do(req)

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return "", 0
	}

	bodyString := string(bodyBytes)
	return bodyString, resp.StatusCode
	//fmt.Println(resp.StatusCode)

	//fmt.Println(resp.Header)
	//fmt.Println(resp.Header["Content-Type"])
	//fmt.Println(resp.Header["Content-Type"][0])

	//fmt.Println(bodyString)

}
