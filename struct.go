package main

import "time"

type Role struct {
	ID   int
	Name string
}

type Project struct {
	Name string `json:"project"`
}

type UserProject struct {
	User    string `json:"user"`
	Project string `json:"project"`
}

type NewBucket struct {
	Host   string `json:"host"`
	Bucket string `json:"bucket"`
}

type NewRootPage struct {
	Page string `json:"page"`
}

type Version struct {
	Value string `json:"version"`
}

type MethodicSet struct {
	Bucket  string `json:"bucket"`
	Version string `json:"version"`
	Page    string `json:"page"`
}

type User struct {
	ID       int
	Username string
	Password string
	Email    string
	Roles    []Role // Добавляем поле для хранения ролей пользователя
}

type ResponseG []struct {
	ID          int    `json:"id"`
	UID         string `json:"uid"`
	OrgID       int    `json:"orgId"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	TypeName    string `json:"typeName"`
	TypeLogoURL string `json:"typeLogoUrl"`
	Access      string `json:"access"`
	URL         string `json:"url"`
	Password    string `json:"password"`
	User        string `json:"user"`
	Database    string `json:"database"`
	BasicAuth   bool   `json:"basicAuth"`
	IsDefault   bool   `json:"isDefault"`
	JSONData    struct {
		DefaultDatabase string `json:"defaultDatabase"`
	} `json:"jsonData,omitempty"`
	ReadOnly bool `json:"readOnly"`
}

type responseDataS struct {
	Datasource struct {
		ID                int    `json:"id"`
		UID               string `json:"uid"`
		OrgID             int    `json:"orgId"`
		Name              string `json:"name"`
		Type              string `json:"type"`
		TypeLogoURL       string `json:"typeLogoUrl"`
		Access            string `json:"access"`
		URL               string `json:"url"`
		Password          string `json:"password"`
		User              string `json:"user"`
		Database          string `json:"database"`
		BasicAuth         bool   `json:"basicAuth"`
		BasicAuthUser     string `json:"basicAuthUser"`
		BasicAuthPassword string `json:"basicAuthPassword"`
		WithCredentials   bool   `json:"withCredentials"`
		IsDefault         bool   `json:"isDefault"`
		JSONData          struct {
		} `json:"jsonData"`
		SecureJSONFields struct {
		} `json:"secureJsonFields"`
		Version  int  `json:"version"`
		ReadOnly bool `json:"readOnly"`
	} `json:"datasource"`
	ID      int    `json:"id"`
	Message string `json:"message"`
	Name    string `json:"name"`
}
type requestV struct {
	STARTTIMEEPOCH string      `json:"START_TIME_EPOCH"`
	ENDTIMEEPOCH   string      `json:"END_TIME_EPOCH"`
	DASHBOARDUUID  string      `json:"DASHBOARD_UUID"`
	DASHBOARDNAME  string      `json:"DASHBOARD_NAME"`
	DASHBOARDGROUP interface{} `json:"DASHBOARD_GROUP"`
	DATASOURCE     string      `json:"DATA_SOURCE"`
	TESTNAME       string      `json:"TEST_NAME"`
	LISTMETOD      []string    `json:"LIST_METOD"`
	PAGEIDCONFL    string      `json:"PAGE_ID_CONFL"`
	HEADREPORT     struct {
		HeadingReport string `json:"heading_report"`
		TextReport    string `json:"text_report"`
		TableStep     string `json:"table_step"`
	} `json:"HEAD_REPORT"`
}

type getPage struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Title   string `json:"title"`
	Version struct {
		By struct {
			Type           string `json:"type"`
			Username       string `json:"username"`
			UserKey        string `json:"userKey"`
			ProfilePicture struct {
				Path      string `json:"path"`
				Width     int    `json:"width"`
				Height    int    `json:"height"`
				IsDefault bool   `json:"isDefault"`
			} `json:"profilePicture"`
			DisplayName string `json:"displayName"`
			Links       struct {
				Self string `json:"self"`
			} `json:"_links"`
			Expandable struct {
				Status string `json:"status"`
			} `json:"_expandable"`
		} `json:"by"`
		When      time.Time `json:"when"`
		Number    int       `json:"number"`
		MinorEdit bool      `json:"minorEdit"`
		Hidden    bool      `json:"hidden"`
		Links     struct {
			Self string `json:"self"`
		} `json:"_links"`
		Expandable struct {
			Content string `json:"content"`
		} `json:"_expandable"`
	} `json:"version"`
	Body struct {
		Storage struct {
			Value          string `json:"value"`
			Representation string `json:"representation"`
			Expandable     struct {
				Content string `json:"content"`
			} `json:"_expandable"`
		} `json:"storage"`
		Expandable struct {
			Editor              string `json:"editor"`
			View                string `json:"view"`
			ExportView          string `json:"export_view"`
			StyledView          string `json:"styled_view"`
			AnonymousExportView string `json:"anonymous_export_view"`
		} `json:"_expandable"`
	} `json:"body"`
	Extensions struct {
		Position int `json:"position"`
	} `json:"extensions"`
	Links struct {
		Webui      string `json:"webui"`
		Edit       string `json:"edit"`
		Tinyui     string `json:"tinyui"`
		Collection string `json:"collection"`
		Base       string `json:"base"`
		Context    string `json:"context"`
		Self       string `json:"self"`
	} `json:"_links"`
	Expandable struct {
		Container    string `json:"container"`
		Metadata     string `json:"metadata"`
		Operations   string `json:"operations"`
		Children     string `json:"children"`
		Restrictions string `json:"restrictions"`
		History      string `json:"history"`
		Ancestors    string `json:"ancestors"`
		Descendants  string `json:"descendants"`
		Space        string `json:"space"`
	} `json:"_expandable"`
}
type confluenceResponseId struct {
	Message string `json:"message"`
	PageID  string `json:"page_id"`
	URL     string `json:"url"`
}
type pgDataTest struct {
	id          int
	application string
	time_start  string
	time_end    string
	confid      string
	host        string
	bucket      string
}
type pgStageResult struct {
	application string
	tag         string
	stage       int
	transaction string
	rpsMean     float64
	rpsMeanRes  float64
	perc95      float64
	perc95Res   float64
	perc99      float64
	perc99Res   float64
	percEr      float64
	percErRes   float64
}
type errorData struct {
	stage          []string
	responseCode   []string
	responseMesage []string
	transactionEr  []string
	count          []string
}
type pushCF struct {
	req  []byte
	stat info
	host string
}
type Report struct {
	textReport string
	host       []string
	pageID     string
	TSS        []string
	TES        []string
	stageCount string
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

type addconfluence struct {
	ID     int    `json:"id"`
	Bucket string `json:"bucket"`
}

type BucketData struct {
	Bucket string `json:"Bucket"`
}
type dataCompareRest struct {
	Application  string `json:"application"`
	ApplicationC string `json:"applicationC"`
	Bucket       string `json:"bucket"`
}
type createBucketApi struct {
	Host   string `json:"host"`
	Bucket string `json:"bucket"`
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
}

type TestData struct {
	Tests []Test `json:"tests"`
}

type AutoGenerated struct {
	Status           string `json:"status"`
	Time             int    `json:"time"`
	Rph              string `json:"rph"`
	BaselineRampup   string `json:"baseline_rampup"`
	BaselinePercent  string `json:"baseline_percent"`
	BaselineDuration string `json:"baseline_duration"`
	StepRampup       string `json:"step_rampup"`
	StepPercent      string `json:"step_percent"`
	StepDuration     string `json:"step_duration"`
	Minpacing        string `json:"minpacing"`
}

type GitEnvData struct {
	Data []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"data"`
}
