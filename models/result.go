package models

type QueryResult struct {
	Page        int         `json:"page"`
	PageSize    int         `json:"pagesize"`
	TotalRecord int64       `json:"total_record"`
	Data        interface{} `json:"data"`
}

func NewQueryResult(page int, pagesize int, total int64, data interface{}) *QueryResult {
	if data == nil {
		data = []string{}
	}
	return &QueryResult{
		Page:        page,
		PageSize:    pagesize,
		TotalRecord: total,
		Data:        data,
	}
}
