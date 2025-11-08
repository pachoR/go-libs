package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var backoffSchedule = []time.Duration{
	1  * time.Second,
	2  * time.Second,
	5  * time.Second,
	10 * time.Second,
}

func Post(url string, payload interface{}) (*http.Response, error) {
	cte := &http.Client{}
	b, _ := json.Marshal(payload)

	log.Println(fmt.Sprintf("bytes: %s", string(b)))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	return cte.Do(req)
}

func PostBody(url string, payload interface{}) ([]byte, error) {
	res, err := Post(url, payload)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return  nil, err
	}

	return body, nil
}

func PostBodyWithRetries(url string, payload interface{}) ([]byte, error) {
	res, err := PostWithRetry(url, payload)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return  nil, err
	}

	return body, nil
}

func GetHTTP(url string) (*http.Response, error) {
	cte := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	return cte.Do(req)
}

func GetBody(url string) ([]byte, error) {
	res, err := GetHTTP(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Get[K any] (url string) (K, error) {
	c := &http.Client{}
	var r K

	rq, _ := http.NewRequest("GET", url, http.NoBody)
	rq.Header.Add("Accept", "application/json")

	res, err := c.Do(rq)
	if err != nil {
		return r, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		var zero K
		return zero, nil
	}

	return r, nil
}

func GetBodyWithRetries(url string) ([]byte, error) {
	res, err := GetWithRetries(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func PostWithRetry(url string, payload interface{}) (*http.Response, error) {
	cte := &http.Client{}
	b, _ := json.Marshal(payload)

	log.Printf("bytes: %s", string(b))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	attempt := 0
	var res *http.Response
	var err error
	for _, backoff := range backoffSchedule {
		res, err = cte.Do(req)
		if err == nil {
			return res, err
		}

		attempt++
		log.Printf("Post failed (Exited with %d): %s, attempt: (%d/4)", res.StatusCode, err.Error(), attempt)
		time.Sleep(backoff)
	}

	return res, fmt.Errorf("Post Request unsuccesfully after 4 retries")
}

func GetWithRetries(url string) (*http.Response, error) {
	cte := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	attempt := 0
	var res *http.Response
	var err error
	for _, backoff := range backoffSchedule {
		res, err = cte.Do(req)
		if err == nil {
			return res, err
		}

		attempt++
		log.Printf("Get failed (Exited with %d): %s, attempt: (%d/4)", res.StatusCode, err.Error(), attempt)
		time.Sleep(backoff)
	}

	return res, fmt.Errorf("Get Request unsuccesfully after 4 retries")
}

func GetWithHeader(url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, http.NoBody)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return client.Do(req)
}


func PostWithHeader(url string, payload interface{}, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return client.Do(req)
}
