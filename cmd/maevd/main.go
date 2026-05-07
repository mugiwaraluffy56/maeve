package main

import (
	"fmt"
	"os"

	"github.com/puneethaditya/maeve/server"
)

func main() {
	srv := server.New(":7432", server.NewRouter(server.RouterConfig{Version: "dev"}))
	if err := srv.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
