package models

type Message struct {
	GroupID TaskGroupID `json:"groupId,omitempty"`
	TaskID  TaskID      `json:"taskId,omitempty"`
	Status  string      `json:"status,omitempty"`

	TimeStamp string `json:"timestamp,omitempty"`
	Code      int    `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
}
