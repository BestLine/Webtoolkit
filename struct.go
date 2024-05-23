package main

type Project struct {
	Name string `json:"project"`
}

type TestStatus struct {
	Status      int         `json:"status"`
	State       string      `json:"state"`
	TestData    TestData2   `json:"test_data"`
	RequestData RequestData `json:"request_data"`
}

type RequestData struct {
	Gitlab   string `json:"gitlab"`
	Count    int    `json:"count"`
	Resource string `json:"resource"`
	Data     []Data `json:"data"`
	Testplan string `json:"testplan"`
}

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TestData2 struct {
	ConfID      int    `json:"conf_id"`
	Application string `json:"application"`
	Bucket      string `json:"bucket"`
	Delivery    string `json:"delivery"`
	Type        string `json:"type"`
	TimeStart   int64  `json:"time_start"`
	TimeEnd     int64  `json:"time_end"`
}

type UserProject struct {
	User    string `json:"user"`
	Project string `json:"project"`
}

type Version struct {
	Value string `json:"version"`
}

type MethodicSet struct {
	Bucket  string `json:"bucket"`
	Version string `json:"version"`
	Page    string `json:"page"`
}

type info struct {
	application       string
	startTestTime     string
	startTestTimeUnix int64
	endTestTime       string
	endTestTimeUnix   int64
	transaction       []string
	bucketInfluxdb    string
	regid             int
}

type TestDataItem struct {
	Application string `json:"application"`
	Bucket      string `json:"bucket"`
	CfURL       string `json:"cfurl"`
}

type StatusTestDataItem struct {
	Project   string `json:"project"`
	Bucket    string `json:"bucket"`
	StartTime string `json:"starttime"`
	Status    string `json:"status"`
	Type      string `json:"type"`
}

type TestsTableData struct {
	Count int            `json:"count"`
	Data  []TestDataItem `json:"data"`
}

type CurrentTestsTableData struct {
	Count int                  `json:"count"`
	Data  []StatusTestDataItem `json:"data"`
}

type Test struct {
	Application string `json:"application"`
	Bucket      string `json:"bucket"`
	State       string `json:"state"`
}

type TestData struct {
	Tests []Test `json:"tests"`
}

type GitEnvData struct {
	Data []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"data"`
	TestPlan []string `json:"testplan"`
}

type GitLabUrl struct {
	Gitlab string `json:"gitlab"`
}
