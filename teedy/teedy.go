package teedy

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jarcoal/httpmock"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	url string

	Tag      *TagService
	Document *DocumentService
	File     *FileService
}

func NewClient(url, username, password string) (*Client, error) {
	trimmedUrl := strings.TrimRight(url, "/")

	client := resty.New().SetHostURL(trimmedUrl)

	resp, err := client.R().SetFormData(map[string]string{
		"username": username,
		"password": password,
	}).Post("api/user/login")

	if err != nil {
		return nil, fmt.Errorf("unable to authenticate: %s", err)
	}

	cookies := resp.Cookies()

	if len(cookies) == 0 {
		return nil, fmt.Errorf("no cookie returned, likely auth failure")
	}

	authCookie := strings.Join([]string{cookies[0].Name, "=", cookies[0].Value}, "")

	client.SetCookies([]*http.Cookie{
		{
			Name:  "Cookies",
			Value: authCookie,
		},
	})

	return &Client{
		url:      trimmedUrl,
		Tag:      NewTagService(client, "tag"),
		Document: NewDocumentService(client, "document"),
		File:     NewFileService(client, "file"),
	}, nil
}

func NewFakeClient() *Client {
	client := resty.New().SetHostURL("http://fake")

	httpmock.ActivateNonDefault(client.GetClient())

	return &Client{
		Tag:      NewTagService(client, "tag"),
		Document: NewDocumentService(client, "document"),
		File:     NewFileService(client, "file"),
	}
}

func GetEnvMustExist(key string) string {
	envvar := os.Getenv(key)

	if envvar == "" {
		log.Fatalf("Unable to load value for '%s' from environment.", key)
	}

	return envvar
}
