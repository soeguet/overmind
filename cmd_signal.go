package main

import (
	"fmt"
	"strings"

	"github.com/DarthSim/overmind/v2/utils"
	"github.com/urfave/cli"
)

type cmdSignalHandler struct{ dialer }

func (h *cmdSignalHandler) Run(c *cli.Context) error {
	if len(c.Args()) < 1 {
		utils.Fatal("Signal name is required")
	}

	conn, err := h.Dial()
	utils.FatalOnErr(err)

	fmt.Fprintf(conn, "signal %v\n", strings.Join(c.Args(), " "))

	return nil
}
