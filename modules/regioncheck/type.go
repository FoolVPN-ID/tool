package regioncheck

import "net/http"

type runnerResultStruct struct {
	Name     string `json:"name"`
	IATACode string `json:"iata_code"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Delay    int    `json:"delay"`
	Error    error  `json:"error"`
	OK       bool   `json:"ok"`
}

type LibraryStruct struct {
	Runner []func(http.Client) runnerResultStruct
	Result []runnerResultStruct
}
