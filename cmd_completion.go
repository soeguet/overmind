package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/DarthSim/overmind/v2/start"
	"github.com/DarthSim/overmind/v2/utils"

	"github.com/urfave/cli"
)

//go:embed completions/_overmind
var zshCompletion string

type cmdCompletionHandler struct{}

func (h *cmdCompletionHandler) Run(c *cli.Context) error {
	if c.NArg() != 1 || c.Args().First() != "zsh" {
		utils.Fatal("Specify a supported shell: zsh")
	}

	fmt.Print(zshCompletion)

	return nil
}

type cmdProcfileNamesHandler struct {
	Procfile string
}

func (h *cmdProcfileNamesHandler) Run(c *cli.Context) error {
	if names := start.ProcfileNames(h.Procfile); len(names) > 0 {
		fmt.Println(strings.Join(names, "\n"))
	}

	return nil
}
