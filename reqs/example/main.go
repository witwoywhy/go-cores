package main

import (
	"github.com/google/uuid"
	"github.com/witwoywhy/go-cores/apps"
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

func getUsers(l logger.Logger) {
	ll := logs.NewSpanLog(l)

	var response User
	var err Error

	client := reqs.NewClient("integrations.getUsers")
	resp := client.Request().
		AddLogger(ll).
		SetBearerAuthToken("TOKEN").
		SetHeader("X-TEST", "TEST").
		SetResult(&response).
		SetError(&err).
		Do()
	if resp.IsErrorState() {
		ll.Errorf("failed to get users: %v", err)
	} else if resp.Error() != nil {
		ll.Errorf("failed to get users unknow error: %v", resp.Error())
	}
}

func getUsersById(l logger.Logger) {
	ll := logs.NewSpanLog(l)

	var response User
	var err Error

	client := reqs.NewClient("integrations.getUser")
	resp := client.Request().
		AddLogger(ll).
		SetBearerAuthToken("TOKEN").
		SetHeader("X-TEST", "TEST").
		SetPathParam("id", "001").
		SetResult(&response).
		SetError(&err).
		Do()
	if resp.IsErrorState() {
		ll.Errorf("failed to get user: %v", err)
	} else if resp.Error() != nil {
		ll.Errorf("failed to get user unknow error: %v", resp.Error())
	}
}

func updateUserById(l logger.Logger) {
	ll := logs.NewSpanLog(l)

	var body User = User{
		Id:          "002",
		Username:    "GG NN",
		Age:         19,
		Permissions: []string{"DELETE"},
	}
	var err Error

	client := reqs.NewClient("integrations.updateUser")
	resp := client.Request().
		AddLogger(ll).
		SetBearerAuthToken("TOKEN").
		SetHeader("X-TEST", "TEST").
		SetHeader("X-API-KEY", "KEY").
		SetError(&err).
		SetBody(body).
		Do()
	if resp.IsErrorState() {
		ll.Errorf("failed to update user: %v", err)
	} else if resp.Error() != nil {
		ll.Errorf("failed to update user unknow error: %v", resp.Error())
	}
}

func main() {
	l := logs.New(map[string]any{
		apps.TraceID: uuid.New().String(),
		apps.SpanID:  uuid.New().String(),
	})

	l.Info("START")
	getUsers(l)
	getUsersById(l)
	updateUserById(l)
	l.Info("END")
}
