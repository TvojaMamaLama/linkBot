package pocketAPI

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	host             = "https://getpocket.com/v3"
	authorizationUrl = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"

	endpointAdd          = "/add"
	endpointRequestToken = "/oauth/request"
	endpointAuthorize    = "/oauth/authorize"

	xErrorHeader   = "X-Error"
	defaultTimeout = 5 * time.Second
)

type (
	requestTokenRequest struct {
		ConsumerKey string `json:"consumer_key"`
		RedirectUri string `json:"redirect_uri"`
	}

	authorizeRequest struct {
		ConsumerKey string `json:"consumer_key"`
		Code        string `json:"code"`
	}

	AuthorizeResponse struct {
		AccessToken string `json:"access_token"`
		Username    string `json:"username"`
	}

	addRequest struct {
		URL         string `json:"url"`
		Title       string `json:"title,omitempty"`
		Tags        string `json:"tags,omitempty"`
		AccessToken string `json:"access_token"`
		ConsumerKey string `json:"consumer_key"`
	}

	AddInput struct {
		URL         string
		Title       string
		Tags        []string
		AccessToken string
	}
)

func (i AddInput) validate() error {
	if i.URL == "" {
		return errors.New("required URL values is empty")
	}
	if i.AccessToken == "" {
		return errors.New("required AccessToken is empty")
	}
	return nil
}

func (i AddInput) generateRequest(consumerKey string) addRequest {
	return addRequest{
		URL:         i.URL,
		Title:       i.Title,
		Tags:        strings.Join(i.Tags, ","),
		AccessToken: i.AccessToken,
		ConsumerKey: consumerKey,
	}
}

type Client struct {
	client      *http.Client
	consumerKey string
}

func NewCLient(consumerKey string) (*Client, error) {
	if consumerKey == "" {
		return nil, errors.New("required ConsumerKey")
	}
	return &Client{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		consumerKey: consumerKey,
	}, nil
}

func (c *Client) GetRequestToken(ctx context.Context, redirectUrl string) (string, error) {
	inp := &requestTokenRequest{
		ConsumerKey: c.consumerKey,
		RedirectUri: redirectUrl,
	}

	value, err := c.doHTTP(ctx, endpointRequestToken, inp)
	if err != nil {
		return "", err
	}

	if value.Get("code") == "" {
		return "", errors.WithMessage(err, "empty response from API")
	}

	return value.Get("code"), nil
}

func (c *Client) GetAuthorizationUrl(requsetToken, redirectUrl string) (string, error) {
	if requsetToken == "" || redirectUrl == "" {
		return "", errors.New("empty params")
	}
	return fmt.Sprintf(authorizationUrl, requsetToken, redirectUrl), nil
}

func (c *Client) doHTTP(ctx context.Context, endPoint string, body interface{}) (url.Values, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return url.Values{}, errors.WithMessage(err, "failed to marshal body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, host+endPoint, bytes.NewBuffer(b))
	if err != nil {
		return url.Values{}, errors.WithMessage(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF8")

	response, err := c.client.Do(req)
	if err != nil {
		return url.Values{}, errors.WithMessage(err, "failed to send request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := fmt.Sprintf("API Error: %s", response.Header.Get(xErrorHeader))
		return url.Values{}, errors.New(err)
	}

	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return url.Values{}, errors.WithMessage(err, "failed to read response body")
	}

	values, err := url.ParseQuery(string(responseByte))
	if err != nil {
		return url.Values{}, errors.WithMessage(err, "failed to parse url values")
	}
	return values, nil
}

func (c *Client) Authorize(ctx context.Context, requestToken string) (*AuthorizeResponse, error) {
	if requestToken == "" {
		return nil, errors.New("empty request token")
	}

	data := &authorizeRequest{
		Code:        requestToken,
		ConsumerKey: c.consumerKey,
	}

	values, err := c.doHTTP(ctx, endpointAuthorize, data)
	if err != nil {
		return nil, err
	}

	accessToken, username := values.Get("access_token"), values.Get("username")
	log.Println()
	if accessToken == "" {
		return nil, errors.New("empty access token in API response")
	}

	return &AuthorizeResponse{
		AccessToken: accessToken,
		Username:    username,
	}, nil
}

func (c *Client) Add(ctx context.Context, inputLink AddInput) error {
	if err := inputLink.validate(); err != nil {
		return err
	}

	req := inputLink.generateRequest(c.consumerKey)
	_, err := c.doHTTP(ctx, endpointAdd, req)

	return err
}
