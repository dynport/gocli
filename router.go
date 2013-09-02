package gocli

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Router struct {
	Actions   map[string]*Action
	Separator string
}

func NewRouter(mapping map[string]*Action) *Router {
	router := &Router{}
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
	re := regexp.MustCompile(strings.Join(patterns, ".*/"))
	actions := make(map[string]*Action)
	for key, action := range router.Actions {
		if re.MatchString(key) {
			actions[key] = action
		}
	}
	return actions
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

func printActionUsage(parts []string, action *Action, message interface{}) {
	table := NewTable()
	fmt.Println("ERROR:", message)
	AddActionUsage(parts, table, action)
	fmt.Println(table.String())
	os.Exit(1)
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
			defer func(parts []string, action *Action) {
				if r := recover(); r != nil {
					printActionUsage(parts, action, r)
				}
			}(parts, action)
			args := action.Args
			if args == nil {
				args = &Args{}
			}
			e := args.Parse(raw[i:])
			if e == nil {
				e = action.Handler(args)
			}
			if e != nil {
				printActionUsage(parts, action, e.Error())
			}
			return nil
		default:
			keys := []string{}
			for key, _ := range actions {
				keys = append(keys, key)
			}
			fmt.Println(cli.UsageForKeys(keys, ""))
			return nil

		}
	}
	fmt.Println(cli.Usage())
	return nil
}
