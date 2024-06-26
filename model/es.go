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
	Match struct {
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
	Total int     `json:"total"`
	Data  []Eslog `json:"data"`
}

type Eslog struct {
	Log   string `json:"log"`
	Time  int64  `json:"time"`
	Name  string `json:"name"`
	Using string `json:"using"`
	Id    string `json:"id"`
}

type QueryBody struct {
	Query struct {
		Bool struct {
			Must []any `json:"must"`
		} `json:"bool"`
	} `json:"query"`
}
