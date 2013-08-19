package main

type Action struct {
	Handler     func(*Args) error
	Args        *Args
	Description string
	Usage       string
}
