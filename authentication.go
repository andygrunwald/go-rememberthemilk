package rememberthemilk

import (
	"context"
	"fmt"
	"strings"
)

const (
	// PermissionRead represents the read permission – gives the ability to read task, contact, group and list details and contents.
	PermissionRead = "read"
	// PermissionWrite represents the write permission – gives the ability to add and modify task, contact, group and list details and contents (also allows you to read).
	PermissionWrite = "write"
	// PermissionDelete represents the delete permission – gives the ability to delete tasks, contacts, groups and lists (also allows you to read and write).
	PermissionDelete = "delete"
)

// AuthorizationsService handles communication with the authorization related
// methods of the Remember The Milk API.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/authentication.rtm
type AuthenticationService service

// AuthenticationURLOptions specifies the parameters to the URL to start the authentication process (for web-based applications or desktop applications)
type AuthenticationURLOptions struct {
	// APIKey represents the API key for the request.
	APIKey string `url:"api_key,omitempty"`

	// Permissions represents the level of access being requested. See constants Permission*.
	Permissions string `url:"perms,omitempty"`

	// Frob represents a frob to be associated with the authentication request.
	Frob string `url:"frob,omitempty"`

	// APISignature represents the API signature for the request.
	// See "Signing Requests" in https://www.rememberthemilk.com/services/api/authentication.rtm for more information.
	APISignature string `url:"api_sig,omitempty"`
}

type GetTokenOptions struct {
	Frob string `url:"frob"`

	BaseAPIURLOptions
}

type GetFrobResponse struct {
	Frob string `json:"frob"`

	BaseResponse
}

type GetTokenResponse struct {
	Authentication Authentication `json:"auth"`

	BaseResponse
}

type Authentication struct {
	Permissions string `json:"perms,omitempty"`
	Token       string `json:"token,omitempty"`
	User        struct {
		Fullname string `json:"fullname,omitempty"`
		ID       string `json:"id,omitempty"`
		Username string `json:"username,omitempty"`
	} `json:"user,omitempty"`
}

// GetAuthenticationURL returns the URL to redirect the user to for authentication with the given permission level.
// This is typically the first step in the authentication process for web-based applications.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/authentication.rtm
func (s *AuthenticationService) GetAuthenticationURL(permission string) (string, error) {
	opts := &AuthenticationURLOptions{
		APIKey:      s.client.apiKey,
		Permissions: permission,
	}
	return s.buildAuthenticationURL(opts)
}

// GetAuthenticationURL returns the URL to redirect the user to for authentication with the given permission level.
// This is typically the first step in the authentication process for desktop applications.
// The frob is obtained by calling the GetFrob method first.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/authentication.rtm
func (s *AuthenticationService) GetAuthenticationURLWithFrob(permission, frob string) (string, error) {
	opts := &AuthenticationURLOptions{
		APIKey:      s.client.apiKey,
		Permissions: permission,
		Frob:        frob,
	}
	return s.buildAuthenticationURL(opts)
}

func (s *AuthenticationService) buildAuthenticationURL(opts *AuthenticationURLOptions) (string, error) {
	u := "auth/"
	u, err := s.client.addOptions(u, opts)
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(s.client.WebBaseURL.Path, "/") {
		return "", fmt.Errorf("baseURL must have a trailing slash, but %q does not", s.client.WebBaseURL)
	}

	authenticationURL, err := s.client.WebBaseURL.Parse(u)
	if err != nil {
		return "", err
	}

	return authenticationURL.String(), nil
}

// GetFrob returns a frob to be used during authentication.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods/rtm.auth.getFrob.rtm
func (s *AuthenticationService) GetFrob(ctx context.Context) (string, *Response, error) {
	opts := s.client.addBaseAPIURLOptions("rtm.auth.getFrob")
	u, err := s.client.addOptions("", opts)
	if err != nil {
		return "", nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return "", nil, err
	}

	var apiResponse struct {
		Response struct {
			GetFrobResponse
		} `json:"rsp"`
	}
	resp, err := s.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return "", resp, err
	}

	return apiResponse.Response.Frob, resp, nil
}

// GetToken returns the auth token for the given frob, if one has been attached.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods/rtm.auth.getToken.rtm
func (s *AuthenticationService) GetToken(ctx context.Context, frob string) (*Authentication, *Response, error) {
	opts := &GetTokenOptions{
		Frob:              frob,
		BaseAPIURLOptions: s.client.addBaseAPIURLOptions("rtm.auth.getToken"),
	}
	u, err := s.client.addOptions("", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var apiResponse struct {
		Response struct {
			GetTokenResponse
		} `json:"rsp"`
	}
	resp, err := s.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return &apiResponse.Response.Authentication, resp, nil
}
