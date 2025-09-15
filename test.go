package rememberthemilk

import (
	"context"
)

// TestService handles communication with the test related
// methods of the Remember The Milk API.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods.rtm
type TestService service

type TestLoginResponse struct {
	User User `json:"user"`

	BaseResponse
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Login represents a testing method which checks if the caller is logged in.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods/rtm.test.login.rtm
func (s *TestService) Login(ctx context.Context) (*User, *Response, error) {
	opts := s.client.addBaseAPIURLOptions("rtm.test.login")
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
			TestLoginResponse
		} `json:"rsp"`
	}
	resp, err := s.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return &apiResponse.Response.TestLoginResponse.User, resp, nil
}
