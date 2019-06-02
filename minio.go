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

func MapAnayltics(al []Analytics) map[string]Analytics {
	aggr := map[string]Analytics{}
	for _, a := range al {
		if found, ok := aggr[a.Filename]; ok {
			found.GetRequestCount = found.GetRequestCount + a.GetRequestCount
			found.HeadRequestCount = found.HeadRequestCount + a.HeadRequestCount
			for ua, i := range a.UserAgentCount {
				if _, ok := found.UserAgentCount[ua]; ok {
					found.UserAgentCount[ua] = found.UserAgentCount[ua] + a.UserAgentCount[ua]
				} else {
					found.UserAgentCount[ua] = i
				}
			}
			aggr[a.Filename] = found
		} else {
			aggr[a.Filename] = a
		}
	}
	return aggr
}
