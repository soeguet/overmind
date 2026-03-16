package start

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/DarthSim/overmind/v2/utils"
)

const (
	headerProcess = "PROCESS"
	headerPid     = "PID"
	headerStatus  = "STATUS"
)

type commandCenter struct {
	cmd      *command
	listener net.Listener
	stop     bool

	SocketPath string
	Network    string
}

func newCommandCenter(cmd *command, socket, network string) *commandCenter {
	return &commandCenter{
		cmd:        cmd,
		SocketPath: socket,
		Network:    network,
	}
}

func (c *commandCenter) Start() (err error) {
	if c.listener, err = net.Listen(c.Network, c.SocketPath); err != nil {
		if c.Network == "unix" && strings.Contains(err.Error(), "address already in use") {
			err = fmt.Errorf("it looks like Overmind is already running. If it's not, remove %s and try again", c.SocketPath)
		}
		return
	}

	go func(c *commandCenter) {
		for {
			if conn, err := c.listener.Accept(); err == nil {
				go c.handleConnection(conn)
			}

			if c.stop {
				break
			}
		}
	}(c)

	return nil
}

func (c *commandCenter) Stop() {
	c.stop = true
	c.listener.Close()
}

func (c *commandCenter) handleConnection(conn net.Conn) {
	re := regexp.MustCompile(`\S+`)

	utils.ScanLines(conn, func(b []byte) bool {
		args := re.FindAllString(string(b), -1)

		if len(args) == 0 {
			return true
		}

		cmd := args[0]

		if len(args) > 1 {
			args = args[1:]
		} else {
			args = []string{}
		}

		switch cmd {
		case "restart":
			c.processRestart(args)
		case "stop":
			c.processStop(args)
		case "quit":
			c.processQuit()
		case "kill":
			c.processKill()
		case "get-connection":
			c.processGetConnection(args, conn)
		case "echo":
			c.processEcho(conn)
		case "status":
			c.processStatus(conn)
		case "status-json":
			c.processStatusJSON(conn)
		}

		return true
	})
}

func (c *commandCenter) processRestart(args []string) {
	for _, p := range c.cmd.processes {
		if len(args) == 0 {
			p.Restart()
			continue
		}

		for _, pattern := range args {
			if utils.WildcardMatch(pattern, p.Name) {
				p.Restart()
				break
			}
		}
	}
}

func (c *commandCenter) processStop(args []string) {
	for _, p := range c.cmd.processes {
		if len(args) == 0 {
			p.Stop(true)
			continue
		}

		for _, pattern := range args {
			if utils.WildcardMatch(pattern, p.Name) {
				p.Stop(true)
				break
			}
		}
	}
}

func (c *commandCenter) processKill() {
	for _, p := range c.cmd.processes {
		p.Kill(false)
	}
}

func (c *commandCenter) processQuit() {
	c.cmd.Quit()
}

func (c *commandCenter) processGetConnection(args []string, conn net.Conn) {
	var proc *process

	if len(args) > 0 {
		name := args[0]

		for _, p := range c.cmd.processes {
			if name == p.Name {
				proc = p
				break
			}
		}
	} else {
		proc = c.cmd.processes[0]
	}

	if proc != nil {
		fmt.Fprintf(conn, "%s %s\n", proc.tmux.Socket, proc.WindowID())
	} else {
		fmt.Fprintln(conn, "")
	}
}

func (c *commandCenter) processEcho(conn net.Conn) {
	c.cmd.output.Echo(conn)
}

func processStatusString(p *process) string {
	if p.dead || p.keepingAlive {
		return "dead"
	}
	return "running"
}

func (c *commandCenter) processStatus(conn net.Conn) {
	maxNameLen := 9
	for _, p := range c.cmd.processes {
		if l := len(p.Name); l > maxNameLen {
			maxNameLen = l
		}
	}

	fmt.Fprint(conn, headerProcess)
	for i := maxNameLen - len(headerProcess); i > -1; i-- {
		conn.Write([]byte{' '})
	}
	fmt.Fprint(conn, "PORT      ")
	fmt.Fprint(conn, headerPid)
	fmt.Fprint(conn, "       ")
	fmt.Fprintln(conn, headerStatus)

	for _, p := range c.cmd.processes {
		utils.FprintRpad(conn, p.Name, maxNameLen+1)
		utils.FprintRpad(conn, strconv.Itoa(p.Port), 10)
		utils.FprintRpad(conn, strconv.Itoa(p.pid), 10)
		fmt.Fprintln(conn, processStatusString(p))
	}
	conn.Close()
}

type processStatusEntry struct {
	Name   string `json:"name"`
	Port   int    `json:"port"`
	PID    int    `json:"pid"`
	Status string `json:"status"`
}

func (c *commandCenter) processStatusJSON(conn net.Conn) {
	entries := make([]processStatusEntry, 0, len(c.cmd.processes))
	for _, p := range c.cmd.processes {
		entries = append(entries, processStatusEntry{
			Name:   p.Name,
			Port:   p.Port,
			PID:    p.pid,
			Status: processStatusString(p),
		})
	}
	json.NewEncoder(conn).Encode(entries)
	conn.Close()
}
