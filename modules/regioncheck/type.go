package regioncheck

import "net/http"

type runnerResultStruct struct {
	Name     string `json:"name"`
	IATACode string `json:"iata_code,omitempty"`
	Region   string `json:"region,omitempty"`
	Country  string `json:"country,omitempty"`
	Delay    int    `json:"delay"`
	Error    error  `json:"error,omitempty"`
}

type LibraryStruct struct {
	Runner []func(http.Client) runnerResultStruct
	Result []runnerResultStruct
}
