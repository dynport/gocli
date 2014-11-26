package main

import (
	"github.com/dynport/gocli"
	"fmt"
	"os"
)

func main() {
	fmt.Println(gocli.Green("gocli test script"))
	router := gocli.NewRouter(
		map[string]*gocli.Action{
			"container/start": {
				Description: "Start container",
				Usage:       "<container_id>",
				Handler: func(args *gocli.Args) error {
					fmt.Println("ACTION: start container")
					return nil
				},
			},
			"container/stop": {
				Description: "Stop container",
				Usage:       "<container_id>",
				Handler: func(args *gocli.Args) error {
					fmt.Println("ACTION: stop container", args.Args)
					return nil
				},
			},
			"image/list": {
				Description: "List Images",
				Handler: func(args *gocli.Args) error {
					fmt.Println("ACTION: list images")
					return nil
				},
			},
		},
	)
	router.Separator = " "
	router.Handle(os.Args)
}
