/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/noctispine/go-email-app/cmd"
	"github.com/noctispine/go-email-app/db"
)

func main() {
	db.NewDB("local_mails")
	cmd.Execute()
}
