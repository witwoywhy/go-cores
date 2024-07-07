package main

import (
	"errors"

	"github.com/google/uuid"
	"github.com/witwoywhy/go-cores/logs"
)

type User struct {
	Id         uuid.UUID `json:"id"`
	FirstName  string    `json:"firstName"`
	Age        int       `json:"age"`
	Permission []string  `json:"permission"`
	Selected   []string  `json:"selected"`
}

func main() {
	l := logs.New(map[string]any{
		"information": map[string]any{
			"userId": uuid.NewString(),
		},
	})

	l.Info("test info")
	l.Infof("test %s %s", "format", "info")

	l.Debug("test debug")
	l.Debugf("test %s %s", "format", "debug")

	l.Warn("test warn")
	l.Warnf("test %s %s", "format", "warn")

	l.Error(errors.New("test error"))
	l.Errorf("test %s %s", "format", "error")

	user := User{
		Id:         uuid.New(),
		FirstName:  "Test Test",
		Age:        16,
		Permission: nil,
		Selected:   []string{"WWW"},
	}
	l.Info(user)

	ll, end := logs.NewSpanLogAction(l, "GET USER")
	defer end()

	ll.Info("USER")
}
