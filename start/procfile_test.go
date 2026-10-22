package start

import (
	"os"
	"path/filepath"
	"reflect"
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
