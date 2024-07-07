package main

import (
	"github.com/google/uuid"
	"github.com/witwoywhy/go-cores/gorms"
	"github.com/witwoywhy/go-cores/logs"
	"github.com/witwoywhy/go-cores/vipers"
	"gorm.io/gorm"
)

func init() {
	vipers.Init()
}

type User struct {
	Id   uuid.UUID
	Name string
}

func mysql() {
	db := gorms.Init("db.mysql")
	l := logs.New(map[string]any{})
	spl := logs.NewSpanLog(l)
	gormLog := gorms.NewGormLog(spl)

	db.AutoMigrate(&User{})

	tx := db.Session(&gorm.Session{Logger: gormLog})

	id := uuid.New()
	u := &User{Id: id, Name: "witwoywhy"}

	tx.Create(u)
	tx.First(u)
	u.Name = "witwoywhy update"
	tx.Updates(u)
	tx.Delete(u)
}

func pg() {
	db := gorms.Init("db.pg")
	l := logs.New(map[string]any{})
	spl := logs.NewSpanLog(l)
	gormLog := gorms.NewGormLog(spl)

	db.AutoMigrate(&User{})

	tx := db.Session(&gorm.Session{Logger: gormLog})

	id := uuid.New()
	u := &User{Id: id, Name: "witwoywhy"}

	tx.Create(u)
	tx.First(u)
	u.Name = "witwoywhy update"
	tx.Updates(u)
	tx.Delete(u)
}

func main() {
	mysql()
	pg()
}
