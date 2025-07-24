package main

import (
	"net/http"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func main() {
	fx.New(
		BuildApp(),
		fx.Invoke(
			func(*grpc.Server) {},
			func(*http.Server) {},
		),
	).Run()
}
