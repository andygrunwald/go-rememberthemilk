package rememberthemilk

import (
	"context"
)

// ContactsService handles communication with the contacts related
// methods of the Remember The Milk API.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods.rtm
type ContactsService service

type ContactsGetListResponse struct {
	Contacts ContactList `json:"contacts"`

	BaseResponse
}

type ContactList struct {
	Contact []Contact `json:"contact"`
}

type Contact struct {
	ID       string `json:"id"`
	FullName string `json:"fullname"`
	Username string `json:"username"`
}

// GetList retrieves a list of contacts.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods/rtm.contacts.getList.rtm
func (s *ContactsService) GetList(ctx context.Context) ([]Contact, *Response, error) {
	opts := s.client.addBaseAPIURLOptions("rtm.contacts.getList")
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
			ContactsGetListResponse
		} `json:"rsp"`
	}
	resp, err := s.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return apiResponse.Response.ContactsGetListResponse.Contacts.Contact, resp, nil
}
