package teedy

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type TeedyAPIError struct {
	GetAll error
	Get    error
	Add    error
	Delete error
	Update error
}

func NewTeedyAPIError(api string) *TeedyAPIError {
	return &TeedyAPIError{
		GetAll: fmt.Errorf("error getting all %ss", api),
		Get:    fmt.Errorf("error getting %s", api),
		Add:    fmt.Errorf("error adding %s", api),
		Delete: fmt.Errorf("error deleting %s", api),
		Update: fmt.Errorf("error updating %s", api),
	}
}

func (n *TeedyAPIError) Custom(description string) error {
	return fmt.Errorf("error %s", description)
}

func checkRequestError(resp *resty.Response, respError, errorDescription error) error {
	if respError != nil {
		return fmt.Errorf("%s: %w", errorDescription, respError)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%s: invalid response code: %s", errorDescription, resp.Status())
	}

	return nil
}
