package minioanalytics

type MinioEvent struct {
	EventName string
	Key       string
	Records   []AWSRecord
}

type Analytics struct {
	Filename         string         `json:"file_name"`
	GetRequestCount  int            `json:"get_request_count"`
	HeadRequestCount int            `json:"head_request_count"`
	UserAgentCount   map[string]int `json:"user_agent_count"`
}
