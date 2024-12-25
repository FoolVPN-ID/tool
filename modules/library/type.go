package library

import "net/http"

type runnerResultStruct struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	Delay  int    `json:"delay"`
	Error  error  `json:"error,omitempty"`
}

type LibraryStruct struct {
	Runner []func(http.Client) runnerResultStruct
	Result []runnerResultStruct
}
