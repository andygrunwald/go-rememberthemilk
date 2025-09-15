package rememberthemilk

import (
	"context"
)

// TagService handles communication with the tags related
// methods of the Remember The Milk API.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods.rtm
type TagService service

type TagsGetListResponse struct {
	Tags TagList `json:"tags"`

	BaseResponse
}

type TagList struct {
	Tags []Tag `json:"tag"`
}

type Tag struct {
	Name string `json:"name"`
}

// GetList retrieves a list of tags.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods/rtm.tags.getList.rtm
func (s *TagService) GetList(ctx context.Context) ([]Tag, *Response, error) {
	opts := s.client.addBaseAPIURLOptions("rtm.tags.getList")
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
			TagsGetListResponse
		} `json:"rsp"`
	}
	resp, err := s.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return apiResponse.Response.Tags.Tags, resp, nil
}
