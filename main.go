/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/noctispine/gogmail/cmd"
	"github.com/noctispine/gogmail/db"
)

func main() {
	db.NewDB("local_mails")
	defer db.CloseDB()
	cmd.Execute()
}
