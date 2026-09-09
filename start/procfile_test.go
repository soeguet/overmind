package start

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestParseProcfileAllowsMissingSpaceAfterColon(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Procfile")
	content := "web:echo hello\nworker: echo world\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	pf := parseProcfile(path, 5000, 100, map[string]int{}, 100, map[string]syscall.Signal{})
	if len(pf) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(pf))
	}
	if pf[0].Name != "web" || pf[0].Command != "echo hello" {
		t.Fatalf("unexpected first entry: %+v", pf[0])
	}
	if pf[1].Name != "worker" || pf[1].Command != "echo world" {
		t.Fatalf("unexpected second entry: %+v", pf[1])
	}
}
