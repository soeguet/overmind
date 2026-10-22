package start

import (
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
)

func TestProcfileNames(t *testing.T) {
	procfile := filepath.Join(t.TempDir(), "Procfile")
	content := "web: bin/rails s\n\n# worker: ignored\nwebpack-dev: bin/webpack-dev-server\njob_queue:   bin/jobs\nbroken:\n"

	if err := os.WriteFile(procfile, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	got := ProcfileNames(procfile)
	want := []string{"web", "webpack-dev", "job_queue"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

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
