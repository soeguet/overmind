package main

import (
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestZshCompletionCoversCommands(t *testing.T) {
	for _, cmd := range newApp().Commands {
		if cmd.Hidden {
			continue
		}

		if !strings.Contains(zshCompletion, "'"+cmd.Name+":") {
			t.Errorf("command %s is not listed", cmd.Name)
		}

		pattern := strings.Join(append([]string{cmd.Name}, cmd.Aliases...), "|")
		if !regexp.MustCompile(`[(|]` + regexp.QuoteMeta(pattern) + `[|)]`).MatchString(zshCompletion) {
			t.Errorf("command %s has no case branch matching %s", cmd.Name, pattern)
		}
	}
}

func TestZshCompletionCoversFlags(t *testing.T) {
	for _, cmd := range newApp().Commands {
		if cmd.Hidden {
			continue
		}

		for _, flag := range cmd.Flags {
			for _, name := range strings.Split(flag.GetName(), ",") {
				name = "--" + strings.TrimSpace(name)
				if len(name) == 3 {
					name = name[1:]
				}

				if !regexp.MustCompile(`[{,']` + regexp.QuoteMeta(name) + `[=\[},]`).MatchString(zshCompletion) {
					t.Errorf("flag %s of command %s is not completed", name, cmd.Name)
				}
			}
		}
	}
}

func TestZshCompletionSyntax(t *testing.T) {
	zsh, err := exec.LookPath("zsh")
	if err != nil {
		t.Skip("zsh is not installed")
	}

	cmd := exec.Command(zsh, "-n")
	cmd.Stdin = strings.NewReader(zshCompletion)

	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
}
