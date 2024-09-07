package request

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

func FetchData(start, end, username, password string) ([]byte, error) {
	// Step 1: Get an anonymous session cookie
	sessionCookie, err := getAnonCookie()
	if err != nil {
		return nil, err
	}

	log.Printf("Anon cookie: %s\n", sessionCookie)

	// Step 2: Login and retrieve a new session cookie
	sessionCookie, err = loginAndGetSessionAndViewState(sessionCookie, username, password)
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
	viewState, j_idt, err := accessPlanningLanding(sessionCookie)
	if err != nil {
		return nil, err
	}

	//curl 'https://mycpe.cpe.fr/faces/Planning.xhtml' -X POST -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:130.0) Gecko/20100101 Firefox/130.0' -H 'Accept: application/xml, text/xml, */*; q=0.01' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br, zstd' -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'Faces-Request: partial/ajax' -H 'X-Requested-With: XMLHttpRequest' -H 'Origin: https://mycpe.cpe.fr' -H 'Connection: keep-alive' -H 'Referer: https://mycpe.cpe.fr/faces/Planning.xhtml' -H 'Cookie: JSESSIONID=2B39334F81CCD20D5D60FCA73B5167E0' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-origin' -H 'TE: trailers' --data-raw 'javax.faces.partial.ajax=true&javax.faces.source=form%3Aj_idt119&javax.faces.partial.execute=form%3Aj_idt119&javax.faces.partial.render=form%3Aj_idt119&form%3Aj_idt119=form%3Aj_idt119&form%3Aj_idt119_start=1725228000000&form%3Aj_idt119_end=1725660000000&form=form&form%3AlargeurDivCenter=&form%3AidInit=webscolaapp.Planning_8686754325772970059&form%3Adate_input=02%2F09%2F2024&form%3Aweek=36-2024&form%3Aj_idt119_view=agendaWeek&form%3AoffsetFuseauNavigateur=-7200000&form%3Aonglets_activeIndex=0&form%3Aonglets_scrollState=0&javax.faces.ViewState=-5487281451116740307%3A-2576554943781788046'

	sessionCookie = "JSESSIONID=2B39334F81CCD20D5D60FCA73B5167E0"

	viewState = "-5487281451116740307:-2576554943781788046"

	j_idt = "j_idt119"

	log.Printf("Planning landing viewState: %s\n", viewState)

	// Step 6: Use retrieved sessionCookie and viewState for the final data request
	body, err := makeFinalDataRequest(sessionCookie, viewState, j_idt, start, end)
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

func loginAndGetSessionAndViewState(username string, password string, anonCookie string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
	}

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
			// trow error to stop the redirect
			return errors.New("stop")
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		// return cookie from reponse
		return resp.Header.Get("Set-Cookie"), nil
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
	urlStr := "https://mycpe.cpe.fr/"

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

func accessPlanningLanding(sessionCookie string) (string, string, error) {

	req, err := http.NewRequest("GET", "https://mycpe.cpe.fr/faces/Planning.xhtml", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authority", "mycpe.cpe.fr")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/jxl,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en,fr;q=0.9,en-GB;q=0.8,en-US;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Cookie", sessionCookie)
	req.Header.Set("Dnt", "1")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://mycpe.cpe.fr/")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"117\", \"Not;A=Brand\";v=\"8\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Reading the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Extract the updated ViewState value from the response HTML
	updatedViewState := extractViewState(string(body))

	j_idt := extractj_idt(string(body))

	return updatedViewState, j_idt, nil
}

func makeFinalDataRequest(sessionCookie, viewState, j_idt, start, end string) ([]byte, error) {
	// URL for the final data request
	urlStr := "https://mycpe.cpe.fr/faces/Planning.xhtml"

	// Using url.Values to construct the data
	data := url.Values{
		"javax.faces.partial.ajax":   {"true"},
		"javax.faces.partial.render": {"form:" + j_idt},
		"form:" + j_idt:              {"form:" + j_idt},
		"form:" + j_idt + "_start":   {start},
		"form:" + j_idt + "_end":     {end},
		"javax.faces.ViewState":      {viewState},
	}

	log.Printf("j_idt: %s, start: %s, end: %s", j_idt, start, end)

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
	re := regexp.MustCompile(`id="j_id1:javax.faces.ViewState:0" value="([^"]*)"`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func extractj_idt(html string) string {
	// Extract the j_idt value using a regular expression
	// regex <div id="form:(?<j_idt>.*?)" class="schedule"
	re := regexp.MustCompile(`<div id="form:(?<j_idt>j_idt.{3,5}?)" class="sch`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
