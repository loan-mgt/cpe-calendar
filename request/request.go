package request

import (
	"cpe/calendar/types"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func FetchData(start, end, username, password string) ([]byte, error) {

	token, err := login(username, password)
	if err != nil {
		return nil, err
	}

	log.Printf("Token: %s\n", token)

	body, err := getCalendar(token, start, end)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func login(username, password string) (types.TokenResponse, error) {
	// Prepare the login request
	urlStr := "https://mycpe.cpe.fr/mobile/login"
	loginData := url.Values{"login": {username}, "password": {password}}

	// Create the request
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(loginData.Encode()))
	if err != nil {
		return types.TokenResponse{}, err
	}

	// Set headers
	req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 15; sdk_gphone64_x86_64 Build/AE3A.240806.005)")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return types.TokenResponse{}, err
	}
	defer resp.Body.Close()

	// Read and unmarshal the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.TokenResponse{}, err
	}

	var formattedResp types.TokenResponse
	if err := json.Unmarshal(body, &formattedResp); err != nil {
		return types.TokenResponse{}, err
	}

	return formattedResp, nil
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

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:130.0) Gecko/20100101 Firefox/130.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Referer", "https://mycpe.cpe.fr/")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", sessionCookie)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("TE", "trailers")

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
