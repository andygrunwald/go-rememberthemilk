package rememberthemilk

import (
	"context"
)

// TimelineService handles communication with the timeline related
// methods of the Remember The Milk API.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/timelines.rtm
type TimelineService service

type TimelinesCreateResponse struct {
	Timeline string `json:"timeline"`

	BaseResponse
}

type Transaction struct {
	ID       string `json:"id"`
	Undoable string `json:"undoable"`
}

// Create returns a new timeline.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods/rtm.timelines.create.rtm
func (s *TimelineService) Create(ctx context.Context) (string, *Response, error) {
	opts := s.client.addBaseAPIURLOptions("rtm.timelines.create")
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
			TimelinesCreateResponse
		} `json:"rsp"`
	}
	resp, err := s.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return "", resp, err
	}

	return apiResponse.Response.TimelinesCreateResponse.Timeline, resp, nil
}
