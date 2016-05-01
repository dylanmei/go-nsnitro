package nsnitro

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/mgutz/logxi/v1"

	jp "github.com/buger/jsonparser"
)

type Config struct {
	URL        string
	User       string
	Password   string
	HTTPClient *http.Client
}

type Client struct {
	config     *Config
	httpClient *http.Client
	version    string
}

type nsrequest struct {
	Type  string
	Name  string
	Query map[string]string
}

type nsresource struct {
	LBMonitor                    *LBMonitor                    `json:"lbmonitor,omitempty"`
	LBVServer                    *LBVServer                    `json:"lbvserver,omitempty"`
	LBVServerServiceGroupBinding *LBVServerServiceGroupBinding `json:"lbvserver_servicegroup_binding,omitempty"`
	Server                       *Server                       `json:"server,omitempty"`
	ServiceGroup                 *ServiceGroup                 `json:"servicegroup,omitempty"`
	ServiceGroupLBMonitorBinding *ServiceGroupLBMonitorBinding `json:"servicegroup_lbmonitor_binding,omitempty"`
	ServiceGroupServerBinding    *ServiceGroupServerBinding    `json:"servicegroup_servicegroupmember_binding,omitempty"`
}

type nsresult struct {
	ErrorCode int    `json:"errorcode"`
	Message   string `json:"message"`
	Severity  string `json:"severity"`
}

func NewClient(config *Config) *Client {
	httpClient := http.DefaultClient
	if config.HTTPClient != nil {
		httpClient = config.HTTPClient
	}

	return &Client{
		config:     config,
		httpClient: httpClient,
	}
}

func (c *Client) Version() (string, error) {
	if c.version != "" {
		return c.version, nil
	}

	req, err := c.request("GET", "config/nsversion", nil)
	if err != nil {
		return "", err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if log.IsDebug() {
		log.Debug("response", "status", res.Status, "body", string(body))
	}

	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Unexpected HTTP status: %d", res.StatusCode))
	}

	c.version, err = jp.GetString(body, "nsversion", "version")
	return c.version, err
}

func (c *Client) do(verb string, request nsrequest, resource *nsresource) error {
	path := "config/" + request.Type
	if request.Name != "" {
		path = path + "/" + request.Name
	}

	if len(request.Query) > 0 {
		path = path + "?" + querystr(request.Query)
	}

	var buffer io.Reader
	var contentType string
	if resource != nil {
		var b []byte
		b, err := json.Marshal(resource)
		if err != nil {
			return err
		}
		buffer = bytes.NewReader(b)
		contentType = fmt.Sprintf("application/vnd.com.citrix.netscaler.%s+json", request.Type)
	}

	req, err := c.request(verb, path, buffer)
	if err != nil {
		return err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if log.IsDebug() {
		if len(body) == 0 {
			log.Debug("response", "status", res.Status)
		} else {
			log.Debug("response", "status", res.Status, "body", string(body))
		}
	}

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to read HTTP response: %s", err.Error()))
		}

		if len(body) == 0 {
			return nil
		}

		var result nsresult
		err = json.Unmarshal(body, &result)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to parse NetScaler API response: %s", err.Error()))
		}

		if result.Severity == "ERROR" {
			return errors.New(result.Message)
		}

		return nil
	}

	return errors.New(fmt.Sprintf("Unexpected HTTP status: %d", res.StatusCode))
}

func (c *Client) fetch(request nsrequest, result interface{}) error {
	path := "config/" + request.Type
	if request.Name != "" {
		path = path + "/" + request.Name
	}

	if len(request.Query) > 0 {
		path = path + "?" + querystr(request.Query)
	}

	req, err := c.request("GET", path, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if log.IsDebug() {
		log.Debug("response", "status", res.Status, "body", string(body))
	}

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		if res.StatusCode == 400 {
			return nil
		}

		data, _, _, err := jp.Get(body, request.Type)
		if err == jp.KeyPathNotFoundError {
			return nil
		}
		if err != nil {
			return err
		}

		return json.Unmarshal(data, result)
	}

	return errors.New(fmt.Sprintf("Unexpected HTTP status: %d", res.StatusCode))
}

func querystr(kvp map[string]string) string {
	values := []string{}
	for key, value := range kvp {

		if !strings.Contains(value, ":") {
			// example: 'bind' in 'action=bind'
			values = append(values, fmt.Sprintf("%s=%s", key, value))
			continue
		}

		criterias := []string{}
		// example: 'name:/test/' in 'filter=name:/test/'
		for _, criteria := range strings.Split(value, ",") {
			parts := strings.Split(criteria, ":")
			if len(parts) == 1 {
				criterias = append(criterias, criteria)
			}
			if len(parts) == 2 {
				criterias = append(criterias, fmt.Sprintf("%s:%s", parts[0], url.QueryEscape(parts[1])))
			}
		}

		values = append(values, fmt.Sprintf("%s=%s", key, strings.Join(criterias, ",")))
	}

	return strings.Join(values, "&")
}

func (c *Client) request(verb, path string, body io.Reader) (*http.Request, error) {
	uri := fmt.Sprintf("%v/nitro/v1/%s", c.config.URL, path)
	req, err := http.NewRequest(verb, uri, body)
	if err != nil {
		return nil, err
	}
	if c.config.User != "" {
		req.Header.Set("X-NITRO-USER", c.config.User)
	}
	if c.config.Password != "" {
		req.Header.Set("X-NITRO-PASS", c.config.Password)
	}

	log.Debug("request", "verb", verb, "uri", uri)
	return req, nil
}
