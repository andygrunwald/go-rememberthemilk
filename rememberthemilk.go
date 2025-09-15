package rememberthemilk

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	defaultAPIVersion = "2"
	defaultWebBaseURL = "https://www.rememberthemilk.com/services/"
	defaultBaseURL    = "https://api.rememberthemilk.com/services/rest/"

	defaultUserAgent = "go-rememberthemilk"

	StatOK   = "ok"
	StatFail = "fail"

	ResponseFormatJSON = "json"
)

var errNonNilContext = errors.New("context must be non-nil")

// A Client manages communication with the Remember The Milk API.
type Client struct {
	apiKey              string
	sharedSecret        string
	authenticationToken string

	// HTTP client used to communicate with the API.
	client *http.Client

	// Web Base URL for authentication requests. Defaults to the public Remember The Milk API
	// WebBaseURL should always be specified with a trailing slash.
	WebBaseURL *url.URL

	// Base URL for API requests. Defaults to the public Remember The Milk API.
	// WebBaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the Remember The Milk API.
	UserAgent string

	// Reuse a single struct instead of allocating one for each service on the heap.
	common service

	// Services used for talking to different parts of the Remember The Milk API.
	Authentication *AuthenticationService
	Tags           *TagService
	Lists          *ListService
	Contacts       *ContactsService
	Timelines      *TimelineService
	Tasks          *TaskService
	Test           *TestService
}

type service struct {
	client *Client
}

// BaseAPIURLOptions specifies the base options that are included in every API request.
type BaseAPIURLOptions struct {
	Method              string `url:"method,omitempty"`
	APIKey              string `url:"api_key,omitempty"`
	AuthenticationToken string `url:"auth_token,omitempty"`
	Format              string `url:"format,omitempty"`
	Version             string `url:"v,omitempty"`
}

// BaseResponse represents the common fields returned in every API response.
// Every specific response embeds this struct.
type BaseResponse struct {
	Stat  string          `json:"stat"`
	Error FailureResponse `json:"err,omitempty"`
}

type FailureResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

// addOptions adds the parameters in opts as URL query parameters to s.
// opts must be a struct whose fields may contain "url" tags.
func (c *Client) addOptions(s string, opts any) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	// Add the API signature
	signature := c.SignRequest(qs)
	qs.Set("api_sig", signature)

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewClient returns a new Remember the Milk API client. If a nil httpClient is
// provided, a new http.Client will be used.
func NewClient(apiKey, sharedSecret, token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient2 := *httpClient
	c := &Client{
		apiKey:              apiKey,
		sharedSecret:        sharedSecret,
		authenticationToken: token,
		client:              &httpClient2,
	}
	c.initialize()
	return c
}

// initialize sets default values and initializes services.
func (c *Client) initialize() {
	if c.client == nil {
		c.client = &http.Client{}
	}

	if c.WebBaseURL == nil {
		c.WebBaseURL, _ = url.Parse(defaultWebBaseURL)
	}
	if c.BaseURL == nil {
		c.BaseURL, _ = url.Parse(defaultBaseURL)
	}
	if c.UserAgent == "" {
		c.UserAgent = defaultUserAgent
	}

	c.common.client = c
	c.Authentication = (*AuthenticationService)(&c.common)
	c.Tags = (*TagService)(&c.common)
	c.Lists = (*ListService)(&c.common)
	c.Contacts = (*ContactsService)(&c.common)
	c.Timelines = (*TimelineService)(&c.common)
	c.Tasks = (*TaskService)(&c.common)
	c.Test = (*TestService)(&c.common)
}

// SetAuthenticationToken sets the authentication token to be used in API requests.
//
// This token is required for making authenticated requests to the Remember The Milk API.
func (c *Client) SetAuthenticationToken(token string) {
	c.authenticationToken = token
}

// RequestOption represents an option that can modify an http.Request.
type RequestOption func(req *http.Request)

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body any, opts ...RequestOption) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

// Response is a Remember the Milk API response. This wraps the standard http.Response
// returned from Remember the Milk and provides convenient access to API specific things.
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// bareDo sends an API request using `caller` http.Client passed in the parameters
// and lets you handle the api response. If an error or API Error occurs, the error
// will contain more information. Otherwise you are supposed to read and close the
// response's Body.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is
// canceled or times out, ctx.Err() will be returned.
func (c *Client) bareDo(ctx context.Context, caller *http.Client, req *http.Request) (*Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	req = req.WithContext(ctx)

	resp, err := caller.Do(req)
	var response *Response
	if resp != nil {
		response = newResponse(resp)
	}

	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return response, ctx.Err()
		default:
		}

		return response, err
	}

	err = CheckResponse(resp)
	if err != nil {
		defer resp.Body.Close()
	}
	return response, err
}

// BareDo sends an API request and lets you handle the api response. If an error
// or API Error occurs, the error will contain more information. Otherwise you
// are supposed to read and close the response's Body.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is
// canceled or times out, ctx.Err() will be returned.
func (c *Client) BareDo(ctx context.Context, req *http.Request) (*Response, error) {
	return c.bareDo(ctx, c.client, req)
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer interface,
// the raw response body will be written to v, without attempting to first
// decode it. If v is nil, and no error happens, the response is returned as is.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*Response, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}
	return resp, err
}

// ErrorResponse reports an error caused by an API request.
//
// Remember the Milk API docs: https://www.rememberthemilk.com/services/api/response.rtm
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response `json:"-"`
	Code     int            `json:"code"`
	Message  string         `json:"message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("[%d] %v", r.Code, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error the response contains `stat`="fail".
func CheckResponse(r *http.Response) error {
	// HTTP error 503 - Service Temporarily Unavailable means "Rate limit hit"
	// See https://www.rememberthemilk.com/services/api/ratelimit.rtm
	if r.StatusCode == http.StatusServiceUnavailable {
		return &ErrorResponse{
			Response: r,
			Code:     r.StatusCode,
			Message:  "Rate limit exceeded. See https://www.rememberthemilk.com/services/api/ratelimit.rtm",
		}
	}

	// The API returns always 200 OK even if an error appears.
	// So we need to parse the response body to check if an error appears.
	errorResponse := &ErrorResponse{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		var apiResponse struct {
			Response struct {
				BaseResponse
			} `json:"rsp"`
		}
		err = json.Unmarshal(data, &apiResponse)
		if err != nil {
			// reset the response as if this never happened
			errorResponse = &ErrorResponse{Response: r}
		}

		if apiResponse.Response.Stat == "fail" {
			errorCode, err := strconv.Atoi(apiResponse.Response.Error.Code)
			if err != nil {
				errorCode = 0
			}
			errorResponse.Code = errorCode
			errorResponse.Message = apiResponse.Response.Error.Message
			return errorResponse
		}

		// Re-populate response body because otherwise the caller won't be able to read it
		r.Body = io.NopCloser(bytes.NewBuffer(data))
	}

	return nil
}

// SignRequest signs a request according to the Remember The Milk API specification.
// It takes a map of parameters, sorts them by key, concatenates them with the shared secret,
// and returns the MD5 hash that should be used as the api_sig parameter.
//
// Remember the Milk API docs: https://www.rememberthemilk.com/services/api/authentication.rtm
func (c *Client) SignRequest(params url.Values) string {
	// Create a slice of keys and sort them
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Concatenate key-value pairs in sorted order
	var paramString strings.Builder
	for _, key := range keys {
		paramString.WriteString(key)
		paramString.WriteString(params.Get(key))
	}

	// Concatenate shared secret with the parameter string
	signString := c.sharedSecret + paramString.String()

	// Calculate MD5 hash
	hash := md5.Sum([]byte(signString))
	return fmt.Sprintf("%x", hash)
}

func (c *Client) addBaseAPIURLOptions(method string) BaseAPIURLOptions {
	opts := BaseAPIURLOptions{
		Method:  method,
		Format:  ResponseFormatJSON,
		Version: defaultAPIVersion,
		APIKey:  c.apiKey,
	}

	if len(c.authenticationToken) > 0 {
		opts.AuthenticationToken = c.authenticationToken
	}

	return opts
}
