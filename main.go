package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// TopLevel of the output
type TopLevel struct {
	Deliveries []Deliveries `json:"deliveries,omitempty"`
}

// Deliveries from the outpu
type Deliveries struct {
	Params Params `json:"params,omitempty"`
}

// Params from the output
type Params struct {
	Matrix []Matrix `json:"matrix,omitempty"`
}

// Matrix from the output
type Matrix struct {
	Day   string  `json:"day"`
	Hours []Hours `json:"hours,omitempty"`
}

// Hours from the output
type Hours struct {
	Available bool   `json:"available"`
	Hour      string `json:"hour"`
}

func notifySlack(day string, hour string, available bool, webhook string) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// payload to HAProxy Data Plane API
	requestBody, err := json.Marshal(map[string]string{
		"text": "Slot available! " + day + " " + hour + " valandą",
	})
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest("POST", webhook, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	// Get the response body as a string
	dataInBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	pageContent := string(dataInBytes)
	fmt.Println(pageContent)

}

func main() {

	email := flag.String("email", "", "login email address")
	password := flag.String("password", "", "login password")
	webhook := flag.String("webhook", "", "Slack Webhook URL")
	flag.Parse()

	data := url.Values{}
	data.Set("rememberMe", "true")
	data.Set("email", *email)
	data.Set("password", *password)
	formData := strings.NewReader(data.Encode())
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Timeout: 60 * time.Second,
		Jar:     cookieJar,
	}

	request, err := http.NewRequest("POST", "https://www.barbora.lt/api/eshop/v1/user/login", formData)
	if err != nil {
		panic(err)
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	regionCookie := http.Cookie{Name: "region", Value: "barbora.lt", Expires: expiration}

	request.AddCookie(&regionCookie)
	request.Header.Add("Authorization", "Basic YXBpa2V5OlNlY3JldEtleQ==")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0")
	request.Header.Add("Pragma", "no-cache")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Connection", "keep-alive")

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	// Get the response body as a string
	dataInBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	pageContent := string(dataInBytes)

	fmt.Println(pageContent)

	// ------------------------------------------------------------------------------------------------------------------------- //

	request, err = http.NewRequest("GET", "https://www.barbora.lt/api/eshop/v1/cart/deliveries", nil)
	if err != nil {
		panic(err)
	}

	request.AddCookie(&regionCookie)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0")
	request.Header.Add("Pragma", "no-cache")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Authorization", "Basic YXBpa2V5OlNlY3JldEtleQ==")

	response, err = client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	var val TopLevel
	err = decoder.Decode(&val)
	if err != nil {
		log.Fatal(err)
	}

	err = decoder.Decode(&val)

	for _, param := range val.Deliveries {
		for _, matrix := range param.Params.Matrix {
			fmt.Println(matrix.Day)
			for _, hour := range matrix.Hours {
				fmt.Println("Delivery time:", hour.Hour)
				fmt.Println("Available:", hour.Available)
				if hour.Available == true {
					notifySlack(matrix.Day, hour.Hour, hour.Available, *webhook)
				}
				fmt.Println()
			}
			fmt.Println("+------------------------------+")
		}
	}
}
