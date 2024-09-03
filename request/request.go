package request

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func FetchData() ([]byte, error) {
	url := "https://mycpe.cpe.fr/faces/Planning.xhtml"
	method := "POST"

	// Data to be sent in the POST request
	data := "javax.faces.partial.ajax=true&javax.faces.source=form%3Aj_idt118&javax.faces.partial.execute=form%3Aj_idt118&javax.faces.partial.render=form%3Aj_idt118&form%3Aj_idt118=form%3Aj_idt118&form%3Aj_idt118_start=1725228000000&form%3Aj_idt118_end=1728684000000&form=form&form%3AlargeurDivCenter=742&form%3AidInit=webscolaapp.Planning_515953417222451125&form%3Adate_input=02%2F09%2F2024&form%3Aweek=36-2024&form%3Aj_idt118_view=agendaWeek&form%3AoffsetFuseauNavigateur=-7200000&form%3Aonglets_activeIndex=0&form%3Aonglets_scrollState=0&javax.faces.ViewState=-2413793382658307369:6268075684027325918"

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Adding headers to the request
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0")
	req.Header.Add("Accept", "application/xml, text/xml, */*; q=0.01")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Faces-Request", "partial/ajax")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Origin", "https://mycpe.cpe.fr")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", "https://mycpe.cpe.fr/faces/Planning.xhtml")
	req.Header.Add("Cookie", "JSESSIONID=58FF76D44D7D2FFC45037DC6DEC307C3")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors") // Updated from "no-cors" to "cors" to match the curl
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Priority", "u=0")
	req.Header.Add("TE", "trailers")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")

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
