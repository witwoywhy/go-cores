package main

import (
	"github.com/witwoywhy/go-cores/logger"
	"github.com/witwoywhy/go-cores/logs"
	"github.com/witwoywhy/go-cores/reqs"
	"github.com/witwoywhy/go-cores/vipers"
)

func init() {
	vipers.Init()
}

type User struct {
	Id          string   `json:"id"`
	Username    string   `json:"username"`
	Age         int      `json:"age"`
	Permissions []string `json:"permissions"`
}

type Error struct {
	Code    string
	Message string
}

func getUser(l logger.Logger) {
	l, end := logs.NewSpanLogAction(logs.L, "GET USER")
	defer end()

	client := reqs.NewClient("integrations.get_user")

	var response User
	var err Error

	resp := client.Request(l).
		SetHeader("X-Api-Key", client.Config().ApiKey).
		SetPathParam("id", "001").
		SetResult(&response).
		SetError(&err).
		Do()
	if resp.IsErrorState() {
		l.Errorf("failed to get user: %v", err)
	} else if resp.Error() != nil {
		l.Errorf("failed to get user unknow error: %v", resp.Error())
	}
}

func updateUser(l logger.Logger) {
	l, end := logs.NewSpanLogAction(logs.L, "UPDATE USER")
	defer end()

	client := reqs.NewClient("integrations.update_user")

	var body User = User{
		Id:          "002",
		Username:    "GG NN",
		Age:         19,
		Permissions: []string{"DELETE"},
	}
	var err Error

	resp := client.Request(l).
		SetHeader("X-API-KEY", "KEY").
		SetPathParam("id", "001").
		SetError(&err).
		SetBody(body).
		Do()
	if resp.IsErrorState() {
		l.Errorf("failed to update user: %v", err)
	} else if resp.Error() != nil {
		l.Errorf("failed to update user unknow error: %v", resp.Error())
	}
}

func main() {
	l, end := logs.NewSpanLogAction(logs.L, "")
	defer end()

	getUser(l)
	updateUser(l)
}
