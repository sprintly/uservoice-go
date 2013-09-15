package uservoice

import (
	"encoding/xml"
	"fmt"
	"github.com/kurrik/oauth1a"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type UservoiceConfig struct {
	Subdomain  string
	ApiKey     string
	ApiSecret  string
	OauthToken string
}

type UservoiceClient struct {
	Subdomain  string
	service    *oauth1a.Service
	userConfig *oauth1a.UserConfig
	baseUrl    string
	client     *http.Client
}

func (u UservoiceClient) UrlForTicket(number int) string {
	return fmt.Sprintf("https://%s.uservoice.com/admin/tickets/%d/", u.Subdomain, number)
}

func (u UservoiceClient) GetTicketByNumber(number int) (UservoiceTicket, error) {
	ticket := UservoiceTicket{} // not sure why I can't return nil as UservoiceTicket
	resp, err := u.makeRequest("GET", fmt.Sprintf("%s/tickets/search?query=number:%d", u.baseUrl, number), nil)
	if err != nil {
		return ticket, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ticket, err
	}
	defer resp.Body.Close()

	results := SearchResponse{}
	err = xml.Unmarshal(body, &results)
	if err != nil {
		return ticket, err
	}
	return results.Tickets[0], nil
}

func (u UservoiceClient) PostNote(ticket_id int, text string) error {
	vals := url.Values{}
	vals.Set("note[text]", text)
	_, err := u.makeRequest("POST",
		fmt.Sprintf("%s/tickets/%d/notes.json", u.baseUrl, ticket_id),
		strings.NewReader(vals.Encode()),
	)
	return err
}

func (u UservoiceClient) makeRequest(method, endpoint string, body io.Reader) (httpResponse *http.Response, err error) {
	httpRequest, _ := http.NewRequest(method, endpoint, body)
	u.service.Sign(httpRequest, u.userConfig)
	httpResponse, err = u.client.Do(httpRequest)
	return
}

func NewUservoiceClient(config UservoiceConfig) UservoiceApi {
	baseUrl := fmt.Sprintf("https://%s.uservoice.com/api/v1", config.Subdomain)
	service := &oauth1a.Service{
		RequestURL:   fmt.Sprintf("%s/oauth/request_token.json", baseUrl),
		AuthorizeURL: fmt.Sprintf("%s/oauth/authorize.json", baseUrl),
		AccessURL:    fmt.Sprintf("%s/oauth/access_token.json", baseUrl),
		ClientConfig: &oauth1a.ClientConfig{
			ConsumerKey:    config.ApiKey,
			ConsumerSecret: config.ApiSecret,
			CallbackURL:    "oob",
		},
		Signer: new(oauth1a.HmacSha1Signer),
	}

	httpClient := new(http.Client)
	userConfig := &oauth1a.UserConfig{}

	return UservoiceClient{config.Subdomain, service, userConfig, baseUrl, httpClient}
}
