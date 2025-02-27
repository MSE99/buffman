package main

import "time"

type request struct {
	Id        int        `json:"id"`
	Done      bool       `json:"done"`
	Payload   string     `json:"string"`
	CreatedOn *time.Time `json:"createdOn"`
	UpdatedOn *time.Time `json:"updatedOn"`
}
