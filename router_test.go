package gocli

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	router := &Router{}
	args := &Args{}
	args.RegisterString("-h", false, "127.0.0.1", "host to use")
	args.RegisterString("-i", true, "", "Image id")
	router.Register(
		"ssh",
		&Action{
			Description: "SSH Into",
			Usage: "<search>",
		},
	)
	router.Register(
		"container/start",
		&Action{
			Description: "start a container",
			Args: args,
			Usage: "<container_id>",
		},
	)
	router.Register(
		"container/stop",
		&Action{
			Description: "stop a container",
			Usage: "<container_id>",
		},
	)
	assert.NotNil(t, router)
	usage := router.Usage()
	assert.Contains(t, usage, "ssh      \t     \t<search>")
	assert.Contains(t, usage, "container\tstart\t<container_id>\tstart a container")
	assert.Contains(t, usage, "container\tstop \t<container_id>\tstop a container")
	assert.Contains(t, usage, `-h DEFAULT: "127.0.0.1" host to use`)
	assert.Contains(t, usage, `-i REQUIRED             Image id`)
}
