package request

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func FetchData(start, end string) ([]byte, error) {
	// Step 1: Login and retrieve session cookie and ViewState
	sessionCookie, viewState, err := loginAndGetSessionAndViewState()
	if err != nil {
		return nil, err
	}

	log.Printf("Session cookie: %s, ViewState: %s", sessionCookie, viewState)

	// Step 2: Make a GET request to /faces/Planning.xhtml to get the updated ViewState
	viewState, err = getUpdatedViewState(sessionCookie, viewState)
	if err != nil {
		return nil, err
	}

	//viewState = "5738921265440495234:7732465887243546242"

	//sessionCookie = "JSESSIONID=AE27628EB9907EE32AF5A5B93B6285A1"

	log.Printf("Session cookie: %s, Updated ViewState: %s", sessionCookie, viewState)

	// Step 3: Use retrieved sessionCookie and viewState for the data request
	urlStr := "https://mycpe.cpe.fr/faces/Planning.xhtml"

	// Using url.Values to construct the data
	data := url.Values{
		"javax.faces.partial.ajax":   {"true"},
		"javax.faces.partial.render": {"form:j_idt118"},
		"form:j_idt118":              {"form:j_idt118"},
		"form:j_idt118_start":        {start},
		"form:j_idt118_end":          {end},
		"javax.faces.ViewState":      {viewState},
	}

	// Converting the data to the appropriate format for the request body
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// Adding essential headers to match the curl command
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Cookie", sessionCookie)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Sending the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Reading the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func loginAndGetSessionAndViewState() (string, string, error) {
	// Retrieve username and password from environment variables
	username := os.Getenv("MYCPE_USERNAME")
	password := os.Getenv("MYCPE_PASSWORD")

	// URL for login
	urlStr := "https://mycpe.cpe.fr/login"

	// Prepare form data for login
	loginData := url.Values{
		"username": {username},
		"password": {password},
	}

	// Creating the login request
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(loginData.Encode()))
	if err != nil {
		return "", "", err
	}

	// Adding headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Upgrade-Insecure-Requests", "1")

	// Sending the login request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Retrieve session cookie
	sessionCookie := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "JSESSIONID" {
			sessionCookie = "JSESSIONID=" + cookie.Value
			break
		}
	}

	// Reading the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Extract the ViewState value from the response HTML
	viewState := extractViewState(string(body))

	return sessionCookie, viewState, nil
}

func getUpdatedViewState(sessionCookie, viewState string) (string, error) {
	// URL for the GET request
	urlStr := "https://mycpe.cpe.fr/faces/Planning.xhtml"

	// Creating the GET request
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}

	// Adding headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", sessionCookie)

	// Sending the GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Reading the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Extract the updated ViewState value from the response HTML
	updatedViewState := extractViewState(string(body))

	return updatedViewState, nil
}

// Function to extract javax.faces.ViewState value using regex
func extractViewState(html string) string {
	re := regexp.MustCompile(`name="javax.faces.ViewState" id="j_id1:javax.faces.ViewState:0" value="([^"]+)"`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
