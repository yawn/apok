package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/yawn/apok"
)

type Client struct {
	action    func() *http.Request
	client    *http.Client
	resources func() *http.Request
}

func New(creds, userInfo, csrfToken string) *Client {

	init := func(action string) func() *http.Request {

		return func() *http.Request {

			req, err := http.NewRequest("GET", fmt.Sprintf("https://policysim.aws.amazon.com/home/data/%s", action), nil)

			if err != nil {
				panic(err)
			}

			req.Header.Add("Accept", "application/json")
			req.Header.Add("Cookie", fmt.Sprintf("aws-creds=%s; aws-userInfo=%s", creds, userInfo))
			req.Header.Add("X-CSRF-Token", csrfToken)

			return req

		}

	}

	return &Client{
		action:    init("action"),
		client:    &http.Client{},
		resources: init("resource"),
	}

}

func (c *Client) Actions(service apok.Service) ([]apok.Action, []byte, error) {

	req := c.action()

	query := req.URL.Query()
	query.Set("serviceName", service.Name)
	query.Set("servicePrefix", service.ActionPrefix)

	req.URL.RawQuery = query.Encode()

	res, err := c.client.Do(req)

	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to request actions for %q", service.Name)
	}

	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		return nil, nil, fmt.Errorf("unexpected text result for actions request for %q", service.Name)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to read actions response body for %q", service.Name)
	}

	// TODO: add new key tester

	var actions []apok.Action

	if err := json.Unmarshal(body, &actions); err != nil {
		return nil, nil, errors.Wrapf(err, "failed to unmarshal actions for %q", service.Name)
	}

	return actions, body, nil

}

func (c *Client) Services() ([]apok.Service, []byte, error) {

	res, err := c.client.Do(c.resources())

	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to request resources")
	}

	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		return nil, nil, errors.New("unexpected text result for resources request")
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to read resource response body")
	}

	// TODO: add new key tester

	var services []apok.Service

	if err := json.Unmarshal(body, &services); err != nil {
		return nil, nil, errors.Wrapf(err, "failed to unmarshal resources")
	}

	return services, body, nil

}
