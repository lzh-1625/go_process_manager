package model

type EsResult struct {
	Took     int    `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   Shards `json:"_shards"`
	Hits     Hits   `json:"hits"`
}
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}
type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}
type Source struct {
	Log   string `json:"log"`
	Name  string `json:"name"`
	Time  int64  `json:"time"`
	Using string `json:"using"`
}
type HitsItem struct {
	Index  string      `json:"_index"`
	ID     string      `json:"_id"`
	Score  interface{} `json:"_score"`
	Source Source      `json:"_source"`
	Sort   []int64     `json:"sort"`
}
type Hits struct {
	Total    Total       `json:"total"`
	MaxScore interface{} `json:"max_score"`
	Hits     []HitsItem  `json:"hits"`
}

type GetLogReq struct {
	FilterName []string `json:"filterName"`
	Match      struct {
		Log   string `json:"log"`
		Name  string `json:"name"`
		Using string `json:"using"`
	} `json:"match"`
	TimeRange struct {
		StartTime int64 `json:"startTime"`
		EndTime   int64 `json:"endTime"`
	} `json:"time"`
	Page struct {
		From int `json:"from"`
		Size int `json:"size"`
	} `json:"page"`
	Sort string `json:"sort"`
}

type EsResp struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore int `json:"max_score"`
		Hits     []struct {
			Index  string `json:"_index"`
			ID     string `json:"_id"`
			Score  int    `json:"_score"`
			Source struct {
				Log   string `json:"log"`
				Name  string `json:"name"`
				Time  int64  `json:"time"`
				Using string `json:"using"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type LogResp struct {
	Total int64        `json:"total"`
	Data  []ProcessLog `json:"data"`
}

type ProcessLog struct {
	Id    int    `json:"id,omitempty" gorm:"primaryKey;autoIncrement;column:id" `
	Log   string `json:"log" gorm:"column:log" type:"text"`
	Time  int64  `json:"time" gorm:"column:time" type:"long"`
	Name  string `json:"name" gorm:"column:name" type:"text"`
	Using string `json:"using" gorm:"column:using" type:"keyword"`
}

func (n *ProcessLog) TableName() string {
	return "process_log"
}

type QueryBody struct {
	MinScore int `json:"min_score,omitempty"`
	Query    struct {
		Bool struct {
			Must               []any `json:"must,omitempty"`
			Should             []any `json:"should,omitempty"`
			MinimumShouldMatch int   `json:"minimum_should_match,omitempty"`
		} `json:"bool,omitempty"`
	} `json:"query,omitempty"`
}
