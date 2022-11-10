package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"crypto/tls"
	"net/http/cookiejar"
	"os"
)

var (
	// ErrNotOK is returned if any of the Envoy APIs does not return a 200
	ErrNotOK = errors.New("server did not return 200")
)

// Client provides the API for interacting with the Envoy APIs
type Client struct {
	address string
	client  *http.Client
}

// NewClient creates a new Client that will talk to an Envoy unit at *address*, creating its own http.Client underneath.
func NewClient(address string) *Client {
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    jar, err := cookiejar.New(nil)
	if err != nil {
		
	}
	client := &http.Client{ Jar: jar,Transport: tr}

	return &Client{
		address: address,
		client:  client,
	}
}

// NewClientWithHTTP creates a new Client that will talk to an Envoy unit at *address* using the provided http.Client.
func NewClientWithHTTP(address string, client *http.Client) *Client {
	return &Client{
		address: address,
		client:  client,
	}
}

func (c *Client) get(url string, response interface{}, jsonO bool) error {

    var bearer = "Bearer " + os.Getenv("ENPHASE_TOKEN") 

    // Create a new request using http
    req, err := http.NewRequest("GET", fmt.Sprintf("https://%s%s", c.address, url), nil)

    // add authorization header to the req
    req.Header.Add("Authorization", bearer)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrNotOK
	}

	if !jsonO {

		s := []byte(`{"token":"valid"}`)
		data := AuthToken{}
		return json.Unmarshal(s, &data)
		//return json.NewDecoder().Decode(response) 
	}

	return json.NewDecoder(resp.Body).Decode(response) 

}

// JWTCheck returns the status of the token
func (c *Client) JWTCheck() ([]AuthToken, error) {
	var authtoken []AuthToken
	err := c.get("/auth/check_jwt", &authtoken, false)
	return authtoken, err
}

// Inventory returns the list of parts installed in the system and registered with the Envoy unit
func (c *Client) Inventory() ([]Inventory, error) {
	var inventory []Inventory
	err := c.get("/inventory.json?deleted=1", &inventory, true)
	return inventory, err
}

// Production returns the current data for Production and Consumption sensors, if equipped.
func (c *Client) Production() (Production, error) {
	var production Production
	err := c.get("/production.json?details=1", &production, true)
	return production, err
}
