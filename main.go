package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(Green("gocli test script"))
	router := NewRouter(
		map[string]*Action{
			"container/start": {
				Description: "Start container",
				Usage:       "<container_id>",
				Handler: func(args *Args) error {
					fmt.Println("ACTION: start container")
					return nil
				},
			},
			"container/stop": {
				Description: "Stop container",
				Usage:       "<container_id>",
				Handler: func(args *Args) error {
					fmt.Println("ACTION: stop container", args.Args)
					return nil
				},
			},
			"image/list": {
				Description: "List Images",
				Handler: func(args *Args) error {
					fmt.Println("ACITON: list images")
					return nil
				},
			},
		},
	)
	router.Separator = " "
	router.Handle(os.Args)
}
