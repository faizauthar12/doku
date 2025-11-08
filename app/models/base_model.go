package models

type ErrorLog struct {
	Line          string      `json:"line,omitempty"`
	Filename      string      `json:"filename,omitempty"`
	Function      string      `json:"function,omitempty"`
	Message       interface{} `json:"message,omitempty"`
	SystemMessage interface{} `json:"system_message,omitempty"`
	Err           error       `json:"-"`
	StatusCode    int         `json:"-"`
}

type Response struct {
	StatusCode    int         `json:"status_code"`
	StatusMessage string      `json:"status_message"`
	Data          interface{} `json:"data,omitempty"`
	Error         *ErrorLog   `json:"error,omitempty"`
	Page          int         `json:"page,omitempty"`
	PerPage       int         `json:"per_page,omitempty"`
	Total         int64       `json:"total,omitempty"`
}
