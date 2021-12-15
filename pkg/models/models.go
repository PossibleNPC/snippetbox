package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("no matching record found")

type Snippet struct {
	ID int
	Title, Content string
	Created, Expires time.Time
}