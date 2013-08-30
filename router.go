package gocli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Router struct {
	Actions   map[string]*Action
	Separator string
	writer io.Writer
}

func NewRouter(mapping map[string]*Action) *Router {
	router := &Router{}
	router.Actions = map[string]*Action{}
	for path, action := range mapping {
		router.Register(path, action)
	}
	return router
}

func (cli *Router) Register(path string, action *Action) {
	if cli.Actions == nil {
		cli.Actions = make(map[string]*Action)
	}
	cli.Actions[path] = action
}

func (router *Router) Search(patterns []string) map[string]*Action {
	actions := make(map[string]*Action)
	for key, action := range router.Actions {
		if keyMatches(key, patterns) {
			actions[key] = action
		}
	}
	return actions
}

func keyMatches(key string, pattern []string) bool {
	for i, key := range strings.Split(key, "/") {
		if i >= len(pattern) {
			return true
		}
		if !strings.HasPrefix(key, pattern[i]) {
			return false
		}
	}
	return true
}

func (cli *Router) Usage() string {
	keys := []string{}
	for key := range cli.Actions {
		keys = append(keys, key)
	}
	return cli.UsageForKeys(keys, "")
}

func (cli *Router) UsageForKeys(keys []string, pattern string) string {
	sort.Strings(keys)
	table := NewTable()
	if cli.Separator != "" {
		table.Separator = cli.Separator
	}
	maxParts := 0
	selected := []string{}
	for _, key := range keys {
		partsCount := len(strings.Split(key, "/"))
		if partsCount > maxParts {
			maxParts = partsCount
		}
		selected = append(selected, key)
	}
	for _, key := range selected {
		parts := strings.Split(key, "/")
		action := cli.Actions[key]

		// fill up parts
		for i := (maxParts - len(parts)); i > 0; i-- {
			parts = append(parts, "")
		}

		parts = append(parts, action.Usage, action.Description)
		table.AddStrings(parts)
		if action.Args != nil {
			usage := action.Args.Usage()
			if usage != "" {
				lines := strings.Split(usage, "\n")
				for _, line := range lines {
					usageParts := []string{}
					for j := 0; j < 3; j++ {
						usageParts = append(usageParts, "")
					}
					current := append(usageParts, line)
					table.AddStrings(current)
				}
			}
		}
	}
	out := []string{"USAGE"}
	out = append(out, table.String())
	return strings.Join(out, "\n")
}

func AddActionUsage(parts []string, table *Table, action *Action) {
	parts = append(parts, action.Usage, action.Description)
	table.AddStrings(parts)
	if action.Args != nil {
		usage := action.Args.Usage()
		if usage != "" {
			lines := strings.Split(usage, "\n")
			for _, line := range lines {
				usageParts := []string{}
				for j := 0; j < 3; j++ {
					usageParts = append(usageParts, "")
				}
				current := append(usageParts, line)
				table.AddStrings(current)
			}
		}
	}
}

func (router *Router) SetWriter(writer io.Writer) {
	router.writer = writer
}

func (router *Router) Writer() io.Writer {
	if router.writer != nil {
		return router.writer
	}
	return os.Stdout
}

func (cli *Router) Handle(raw []string) error {
	for i := len(raw); i > 0; i-- {
		parts := raw[1:i]
		actions := cli.Search(parts)
		switch len(actions) {
		case 0:
			continue
		case 1:
			var action *Action
			for k, a := range actions {
				parts = strings.Split(k, "/")
				action = a
			}
			args := action.Args
			if args == nil {
				args = &Args{}
			}
			e := args.Parse(raw[i:])
			if e == nil {
				e = action.Handler(args)
			}
			if e != nil {
				table := NewTable()
				fmt.Fprintln(cli.Writer(), "ERROR: " + e.Error())
				AddActionUsage(parts, table, action)
				fmt.Fprintln(cli.Writer(), table.String())
				os.Exit(1)
			}
			return nil
		default:
			keys := []string{}
			for key, _ := range actions {
				keys = append(keys, key)
			}
			fmt.Fprintln(cli.Writer(), cli.UsageForKeys(keys, ""))
			return nil

		}
	}
	fmt.Fprintln(cli.Writer(), cli.Usage())
	return nil
}
