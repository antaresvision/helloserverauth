package api

import "time"

type Response struct {
	Message string		`json:"message"`
	TimeStamp time.Time `json:"time_stamp"`
}
