package teedy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	url        string
	authCookie string
	httpclient HTTPClient

	Tag *TagService
}

func NewClient(httpClient HTTPClient, teedyUrl, username, password string) (*Client, error) {
	trimmedUrl := strings.TrimRight(teedyUrl, "/")

	s := &Client{
		url: trimmedUrl,
	}

	if httpClient == nil {
		s.httpclient = &http.Client{}
	}

	formData := url.Values{
		"username": {username},
		"password": {password},
	}

	tokenReqUrl := strings.Join([]string{trimmedUrl, "api/user/login"}, "/")

	resp, err := http.PostForm(tokenReqUrl, formData)

	if err != nil {
		log.Fatalf("Unable to get token from %s: %v", tokenReqUrl, err)
	}

	cookies := resp.Cookies()

	s.authCookie = strings.Join([]string{cookies[0].Name, "=", cookies[0].Value}, "")

	s.Tag = NewTagService(s)

	return s, nil
}

func NewFakeClient(httpClient HTTPClient) *Client {
	s := &Client{
		url:        "https://fake.local",
		httpclient: httpClient,
	}

	s.Tag = NewTagService(s)

	return s
}

func (t *Client) getAPIUrl(apiUrl string) string {
	return strings.Join([]string{t.url, apiUrl}, "/")
}

func (t *Client) request(endpoint string, method string, body url.Values, returnType interface{}) (interface{}, error) {
	fullUrl := t.getAPIUrl(endpoint)

	req, err := http.NewRequest(method, fullUrl, strings.NewReader(body.Encode()))

	if err != nil {
		return nil, fmt.Errorf("error building request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))

	// build the auth header
	req.Header.Set("Cookie", t.authCookie)

	resp, err := t.httpclient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response code: %s", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	return unmarshalResponse(data, returnType)
}

func unmarshalResponse(data []byte, t interface{}) (interface{}, error) {
	err := json.Unmarshal(data, t)

	if err != nil {
		fmt.Printf("can't unmarshall json = %v \n", err)
	}

	return t, nil
}
