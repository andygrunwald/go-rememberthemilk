package rememberthemilk

import (
	"context"
)

// ListService handles communication with the lists related
// methods of the Remember The Milk API.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods.rtm
type ListService service

type ListsGetListResponse struct {
	Lists ListList `json:"lists"`

	BaseResponse
}

type ListList struct {
	List []List `json:"list"`
}

type List struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Deleted    string `json:"deleted"`
	Locked     string `json:"locked"`
	Archived   string `json:"archived"`
	Position   string `json:"position"`
	Smart      string `json:"smart"`
	SortOrder  string `json:"sort_order"`
	Permission string `json:"permission"`
}

// GetList retrieves a list of lists.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods/rtm.lists.getList.rtm
func (s *ListService) GetList(ctx context.Context) ([]List, *Response, error) {
	opts := s.client.addBaseAPIURLOptions("rtm.lists.getList")
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
			ListsGetListResponse
		} `json:"rsp"`
	}
	resp, err := s.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return apiResponse.Response.ListsGetListResponse.Lists.List, resp, nil
}
