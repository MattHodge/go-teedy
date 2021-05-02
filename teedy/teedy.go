package teedy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	url        string
	authCookie string
	httpclient HTTPClient

	Tag      *TagService
	Document *DocumentService
	File     *FileService
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

	if len(cookies) == 0 {
		return nil, fmt.Errorf("no cookie returned, likely auth failure")
	}

	s.authCookie = strings.Join([]string{cookies[0].Name, "=", cookies[0].Value}, "")

	s.Tag = NewTagService(s)
	s.Document = NewDocumentService(s)
	s.File = NewFileService(s)

	return s, nil
}

func NewFakeClient(httpClient HTTPClient) *Client {
	s := &Client{
		url:        "https://fake.local",
		httpclient: httpClient,
	}

	s.Tag = NewTagService(s)
	s.Document = NewDocumentService(s)

	return s
}

func (t *Client) getAPIUrl(apiUrl string) string {
	return strings.Join([]string{t.url, apiUrl}, "/")
}

func (t *Client) multipartUpload(endpoint, method string, params map[string]string, file *os.File) ([]byte, error) {

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	b := new(bytes.Buffer)
	writer := multipart.NewWriter(b)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := t.newRequest(method, endpoint, b)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	return t.doRequest(req)
}

func (t *Client) request(endpoint string, method string, body url.Values) ([]byte, error) {
	req, err := t.newRequest(method, endpoint, strings.NewReader(body.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return t.doRequest(req)
}

func (t *Client) requestUnmarshal(endpoint string, method string, body url.Values, returnType interface{}) (interface{}, error) {
	data, err := t.request(endpoint, method, body)

	if err != nil {
		return nil, err
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

func (t *Client) newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	fullUrl := t.getAPIUrl(endpoint)

	req, err := http.NewRequest(method, fullUrl, body)

	if err != nil {
		return nil, fmt.Errorf("unable to build new request: %v", err)
	}

	req.Header.Set("Cookie", t.authCookie)

	return req, nil
}

func (t *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := t.httpclient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error sending requestUnmarshal: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response code: %s", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	return data, nil
}
