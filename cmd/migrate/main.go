package main

import (
	"github.com/vukieuhaihoa/user-service/internal/infrastructure"
)

func main() {
	_ = infrastructure.CreateSQLDBAndMigration()
}
