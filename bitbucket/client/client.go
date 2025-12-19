package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-http-utils/headers"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	BBTokenServerUrl = "https://bitbucket.org/site/oauth2/access_token"
	BBApiServerUrl   = "https://api.bitbucket.org/2.0"
	ApplicationJson  = "application/json"
	FormEncoded      = "application/x-www-form-urlencoded"
	Bearer           = "Bearer"
	IdSeparator      = ":"
)

type Client struct {
	accessToken  string
	clientId     string
	clientSecret string
	numRetries   int
	retryDelay   int
	httpClient   *http.Client
}

func NewClient(ctx context.Context, accessToken string, clientId string, clientSecret string, numRetries int, retryDelay int) (*Client, error) {
	c := &Client{
		accessToken:  accessToken,
		clientId:     clientId,
		clientSecret: clientSecret,
		numRetries:   numRetries,
		retryDelay:   retryDelay,
		httpClient:   &http.Client{},
	}
	//Check for client credentials authentication and try to get access token
	if c.accessToken == "" {
		tflog.Info(ctx, "Bitbucket API: Obtaining access token...")
		requestForm := url.Values{
			"grant_type":    []string{"client_credentials"},
			"client_id":     []string{c.clientId},
			"client_secret": []string{c.clientSecret},
		}
		req, err := http.NewRequest(http.MethodPost, BBTokenServerUrl, bytes.NewBufferString(requestForm.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set(headers.ContentType, FormEncoded)
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			tflog.Info(ctx, "Bitbucket API:", map[string]interface{}{"error": err})
		} else {
			tflog.Info(ctx, "Bitbucket API: ", map[string]interface{}{"request": string(requestDump)})
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
		}
		if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
			respBody := new(bytes.Buffer)
			_, err := respBody.ReadFrom(resp.Body)
			if err != nil {
				return nil, &RequestError{StatusCode: resp.StatusCode, Err: err}
			}
			return nil, &RequestError{StatusCode: resp.StatusCode, Err: fmt.Errorf("%s", respBody.String())}
		}
		//Parse body to extract access_token
		token := &OauthToken{}
		err = json.NewDecoder(resp.Body).Decode(token)
		if err != nil {
			return nil, err
		}
		tflog.Info(ctx, "Bitbucket API: Received access token: ", map[string]interface{}{"token": token.AccessToken})
		//Inject token as access_token for client for all future calls
		c.accessToken = token.AccessToken
	}
	return c, nil
}

func (c *Client) HttpRequest(ctx context.Context, method string, path string, query url.Values, headerMap http.Header, body *bytes.Buffer) (*bytes.Buffer, error) {
	req, err := http.NewRequest(method, c.RequestPath(path), body)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	//Handle query values
	if query != nil {
		requestQuery := req.URL.Query()
		for key, values := range query {
			for _, value := range values {
				requestQuery.Add(key, value)
			}
		}
		req.URL.RawQuery = requestQuery.Encode()
	}
	//Handle header values
	if headerMap != nil {
		for key, values := range headerMap {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
	//Handle authentication
	if c.accessToken != "" {
		req.Header.Set(headers.Authorization, Bearer+" "+c.accessToken)
	}
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		tflog.Info(ctx, "Bitbucket API:", map[string]any{"error": err})
	} else {
		tflog.Info(ctx, "Bitbucket API: ", map[string]any{"request": string(requestDump)})
	}
	try := 0
	var resp *http.Response
	for {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
		}
		if (resp.StatusCode == http.StatusTooManyRequests) || (resp.StatusCode >= http.StatusInternalServerError) {
			try++
			if try >= c.numRetries {
				break
			}
			time.Sleep(time.Duration(c.retryDelay) * time.Second)
			continue
		}
		break
	}
	defer resp.Body.Close()
	respBody := new(bytes.Buffer)
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		return nil, &RequestError{StatusCode: resp.StatusCode, Err: err}
	}
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		return nil, &RequestError{StatusCode: resp.StatusCode, Err: fmt.Errorf("%s", respBody.String())}
	}
	return respBody, nil
}

func (c *Client) RequestPath(path string) string {
	return fmt.Sprintf("https://%s/%s", BBApiServerUrl, path)
}
