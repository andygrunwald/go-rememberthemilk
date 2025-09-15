package rememberthemilk

import (
	"context"
)

// TagService handles communication with the tasks related
// methods of the Remember The Milk API.
//
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/tasks.rtm
type TaskService service

const (
	ParseWithoutSmartAdd = 0
	ParseWithSmartAdd    = 1
)

type TaskInput struct {
	Timeline   string `url:"timeline,omitempty"`
	ListID     string `url:"list_id,omitempty"`
	Name       string `url:"name,omitempty"`
	Parse      int    `url:"parse,omitempty"`
	ParentID   string `url:"parent_task_id,omitempty"`
	ExternalID string `url:"external_id,omitempty"`
	GiveTo     string `url:"give_to,omitempty"`

	BaseAPIURLOptions
}

type Task struct {
	ID           string `json:"id"`
	Due          string `json:"due"`
	HasDueTime   string `json:"has_due_time"`
	Added        string `json:"added"`
	Completed    string `json:"completed"`
	Deleted      string `json:"deleted"`
	Priority     string `json:"priority"`
	Postponed    string `json:"postponed"`
	Estimate     string `json:"estimate"`
	Start        string `json:"start"`
	HasStartTime string `json:"has_start_time"`
}

type Taskseries struct {
	ID           string `json:"id"`
	Created      string `json:"created"`
	Modified     string `json:"modified"`
	Name         string `json:"name"`
	Source       string `json:"source"`
	URL          string `json:"url"`
	LocationID   string `json:"location_id"`
	ParentTaskID string `json:"parent_task_id"`

	// TODO Missing fields
	// - tags
	// - notes
	// - participants

	Task []Task `json:"task"`
}

type TaskList struct {
	ID         string       `json:"id"`
	Taskseries []Taskseries `json:"taskseries"`
}

type TaskAddResponse struct {
	Transaction Transaction `json:"transaction"`
	List        TaskList    `json:"list"`

	BaseResponse
}

// Add adds a task, TaskInput.Name, to the list specified by TaskInput.ListID.
// If TaskInput.ListID is omitted, the task will be added to the Inbox.
// If TaskInput.Parse is ParseWithSmartAdd (1), Smart Add will be used to process the task.
// If TaskInput.ParentTaskID is provided and the user has a Pro account, the new task is created as a sub-task, with the list of the TaskInput.ParentTaskID taking priority over the provided TaskInput.ListID.
//
// This method requires a timeline.
//
// Docs about Smart Add: https://www.rememberthemilk.com/help/?ctx=basics.smartadd.whatis
// Remember The Milk API docs: https://www.rememberthemilk.com/services/api/methods/rtm.tasks.add.rtm
func (s *TaskService) Add(ctx context.Context, task TaskInput) (*TaskAddResponse, *Response, error) {
	task.BaseAPIURLOptions = s.client.addBaseAPIURLOptions("rtm.tasks.add")
	u, err := s.client.addOptions("", task)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var apiResponse struct {
		Response struct {
			TaskAddResponse
		} `json:"rsp"`
	}
	resp, err := s.client.Do(ctx, req, &apiResponse)
	if err != nil {
		return nil, resp, err
	}

	return &apiResponse.Response.TaskAddResponse, resp, nil
}
