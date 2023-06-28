package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client provides a connection to the Confluence API
type Client struct {
	client    *http.Client
	baseURL   *url.URL
	basePath  string
	publicURL *url.URL
}

// NewClientInput provides information to connect to the Confluence API
type NewClientInput struct {
	Hostname string
	Port     int
	UseTls   bool
	Username string
	Password string
}

// ErrorResponse describes why a request failed
type ErrorResponse struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

// NewClient returns an authenticated client ready to use
func NewClient(input *NewClientInput) *Client {
	publicURL := url.URL{
		Scheme: ifThenElse(input.UseTls, "https", "http").(string),
		Host:   fmt.Sprintf(`%s:%d`, input.Hostname, input.Port),
	}

	// Default
	basePath := "/api"

	baseURL := url.URL{
		Scheme: ifThenElse(input.UseTls, "https", "http").(string),
		Host:   fmt.Sprintf(`%s:%d`, input.Hostname, input.Port),
	}
	baseURL.User = url.UserPassword(input.Username, input.Password)
	return &Client{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL:   &baseURL,
		basePath:  basePath,
		publicURL: &publicURL,
	}
}

// GetString uses the client to send a GET request and returns a string
func (c *Client) GetString(path string) (string, error) {
	body := new(bytes.Buffer)
	responseBody, err := c.doRaw("GET", path, "", body)
	if err != nil {
		return "", err
	}
	result := responseBody.String()
	return result, nil
}

// Get uses the client to send a GET request
func (c *Client) Get(path string, result interface{}) error {
	body := new(bytes.Buffer)
	return c.do("GET", path, "", body, result)
}

func (c *Client) GetRaw(path string) (*bytes.Buffer, error) {
	body := new(bytes.Buffer)
	return c.doRaw("GET", path, "", body)
}

// Delete uses the client to send a DELETE request
func (c *Client) Delete(path string) error {
	body := new(bytes.Buffer)
	return c.do("DELETE", path, "", body, nil)
}

// Post uses the client to send a POST request
func (c *Client) Post(path string, body interface{}, result interface{}, itemsToRemove []string) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if itemsToRemove != nil && len(itemsToRemove) > 0 {
		RemoveKeysFromJSONObjectBytes(&bodyBytes, itemsToRemove)
	}
	b := bytes.NewBuffer(bodyBytes)
	if err != nil {
		return err
	}
	return c.do("POST", path, "application/json", b, result)
}

// Put uses the client to send a PUT request
func (c *Client) Put(path string, body interface{}, result interface{}, itemsToRemove []string) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if itemsToRemove != nil && len(itemsToRemove) > 0 {
		RemoveKeysFromJSONObjectBytes(&bodyBytes, itemsToRemove)
	}
	b := bytes.NewBuffer(bodyBytes)
	if err != nil {
		return err
	}
	return c.do("PUT", path, "application/json", b, result)
}

func JsonBytesBuffer(body interface{}) (*bytes.Buffer, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(bodyBytes), nil
}

func bytesBytesBuffer(bodyBytes []byte) *bytes.Buffer {
	return bytes.NewBuffer(bodyBytes)
}

func bytesBufferJSON(bodyBytes *bytes.Buffer, result interface{}) error {
	if result == nil {
		return nil
	}
	reader := bytes.NewReader(bodyBytes.Bytes())
	return json.NewDecoder(reader).Decode(&result)
}

func (c *Client) do(method, path, contentType string, body *bytes.Buffer, result interface{}) error {
	responseBody, err := c.doRaw(method, path, contentType, body)
	if err != nil {
		return err
	}
	return bytesBufferJSON(responseBody, result)
}

// do uses the client to send a specified request
func (c *Client) doRaw(method, path, contentType string, body *bytes.Buffer) (*bytes.Buffer, error) {
	fullPath := c.basePath + path
	u, err := c.baseURL.Parse(fullPath)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Add("kbn-xsrf", "monitoring")
	// LOGGING LOCALLY FOR DEBUGGIN PURPOSES. Uncomment to view raw requests.
	// f, err := os.OpenFile("/tmp/httputildebug.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	// f.WriteString("[doRaw]\n")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()
	// dump, err := httputil.DumpRequestOut(req, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _, err2 := f.Write(dump)
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }
	// f.WriteString("\n")
	//
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var expectedStatusCode = map[string][]int{
		"POST":   {200, 201},
		"PUT":    {200},
		"GET":    {200},
		"DELETE": {200, 204},
	}
	if !contains(expectedStatusCode[method], resp.StatusCode) {
		var responseBody string
		var errResponse ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errResponse)
		if err != nil {
			responseBody = "Could not decode error"
		} else {
			responseBody = errResponse.String()
		}
		s := body.String()
		return nil, fmt.Errorf("%s\n\n%s %s\n%s\n\n%s",
			resp.Status, method, fullPath, s, responseBody)
	}
	result := new(bytes.Buffer)
	_, err = result.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *ErrorResponse) String() string {
	return fmt.Sprintf("%s\nCode: %d",
		e.Message, e.StatusCode)
}

// URL returns the public URL for a given path
func (c *Client) URL(path string) string {
	u, err := c.publicURL.Parse(path)
	if err != nil {
		return ""
	}
	return u.String()
}
