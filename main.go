package main

import (
	_ "github.com/lib/pq"
	"github.com/r4f3t/messagesender/cmd"
)

func main() {
	cmd.Execute()
}
