package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Product Page response
type ProductPage struct {
	Products []Product `json:"objects"`
	Next     string    `json:"next"`
}

// Kogan API Client
type Client struct {
	ApiEndpoint url.URL
}

// get list of products via page id path
// Context param as placeholder for future improvements
func (c *Client) GetProductPage(ctx context.Context, path string) (*ProductPage, error) {

	// create a new request query to append the page id to the path
	endPoint := fmt.Sprintf("%s://%s%s", c.ApiEndpoint.Scheme, c.ApiEndpoint.Hostname(), path)

	// get list of products
	resp, err := http.Get(endPoint)
	if err != nil {
		return nil, fmt.Errorf("the get product page request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response data failed with error %s\n", err)
	}

	// unmarshall the response body
	productPage := &ProductPage{}
	err = json.Unmarshal(data, productPage)
	if err != nil {
		return nil, fmt.Errorf("unmarshal process failed with error %s\n", err)
	}

	return productPage, nil
}

// Create a new client
func NewClient(apiEndpoint *url.URL) (*Client, error) {
	var client = &Client{}
	client.ApiEndpoint = *apiEndpoint

	return client, nil
}
