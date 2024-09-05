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
	// Step 1: Get an anonymous session cookie
	sessionCookie, err := getAnonCookie()
	if err != nil {
		return nil, err
	}

	log.Printf("Anon cookie: %s\n", sessionCookie)

	// Step 2: Login and retrieve a new session cookie
	sessionCookie, err = loginAndGetSessionAndViewState(sessionCookie)
	if err != nil {
		return nil, err
	}

	log.Printf("Session cookie: %s\n", sessionCookie)

	// Step 3: Retrieve the ViewState for the logged-in session
	viewState, err := getUpdatedViewState(sessionCookie)
	if err != nil {
		return nil, err
	}

	log.Printf("ViewState: %s\n", viewState)

	// Step 4: Access MainMenuPage.xhtml to maintain session
	err = accessMainMenu(sessionCookie, viewState)
	if err != nil {
		return nil, err
	}

	log.Printf("Accessed MainMenuPage.xhtml\n")

	// Step 5: Access Planning landing page
	viewState, err = accessPlanningLanding(sessionCookie)
	if err != nil {
		return nil, err
	}

	log.Printf("Planning landing viewState: %s\n", viewState)

	// Step 6: Use retrieved sessionCookie and viewState for the final data request
	body, err := makeFinalDataRequest(sessionCookie, viewState, start, end)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getAnonCookie() (string, error) {
	// URL for initial anonymous request
	urlStr := "https://mycpe.cpe.fr/faces/Login.xhtml"

	// Create a GET request
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}

	// Add headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0")

	// Send the request
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Keep the cookies from the redirect chain
			req.Header.Add("Cookie", via[0].Header.Get("Cookie"))
			return nil
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
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

	return sessionCookie, nil
}

func loginAndGetSessionAndViewState(anonCookie string) (string, error) {
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
		return "", err
	}

	// Adding headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", anonCookie)

	// Sending the login request
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Keep the cookies from the redirect chain
			req.Header.Add("Cookie", via[0].Header.Get("Cookie"))
			return nil
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
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

	return sessionCookie, nil
}

func getUpdatedViewState(sessionCookie string) (string, error) {
	// URL for the GET request
	urlStr := "https://mycpe.cpe.fr/faces/Planning.xhtml"

	// Creating the GET request
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}

	// Adding headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0")
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

func accessMainMenu(sessionCookie, viewState string) error {
	// URL for accessing MainMenuPage.xhtml
	urlStr := "https://mycpe.cpe.fr/faces/MainMenuPage.xhtml"

	// Prepare form data for the MainMenuPage request
	mainMenuData := url.Values{
		"form":                  {"form"},
		"form:largeurDivCenter": {"827"},
		"form:idInit":           {"webscolaapp.MainMenuPage_-518408921344646904"},
		"form:sauvegarde":       {""},
		"javax.faces.ViewState": {viewState},
		"form:sidebar":          {"form:sidebar"},
		"form:sidebar_menuid":   {"8"},
	}

	// Creating the MainMenuPage request
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(mainMenuData.Encode()))
	if err != nil {
		return err
	}

	// Adding headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", sessionCookie)

	// Sending the request
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.Header.Add("Cookie", via[0].Header.Get("Cookie"))
			return nil
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func accessPlanningLanding(sessionCookie string) (string, error) {
	// URL for accessing Planning landing page
	urlStr := "https://mycpe.cpe.fr/faces/Planning.xhtml"

	// Creating the Planning landing page request
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}

	// Adding headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0")
	req.Header.Add("Cookie", sessionCookie)

	// Sending the request
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

func makeFinalDataRequest(sessionCookie, viewState, start, end string) ([]byte, error) {
	// URL for the final data request
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

	// Creating the request
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// Adding headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", sessionCookie)

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

func extractViewState(html string) string {
	// Extract the ViewState value using a regular expression
	re := regexp.MustCompile(`id="javax.faces.ViewState" value="([^"]*)"`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
