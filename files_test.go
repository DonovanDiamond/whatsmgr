package whatsmgr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHashFile(t *testing.T) {
	conn := Connection{}
	data := []byte("hello world")
	expectedHash := "11-2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"

	hash := conn.hashFile(data)
	if hash != expectedHash {
		t.Errorf("Expected hash %s, got %s", expectedHash, hash)
	}
}

func TestWriteFileIfNotExists_NewFile(t *testing.T) {
	conn := Connection{}
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	data := []byte("hello")

	err := conn.writeFileIfNotExists(filePath, data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify file exists and content matches
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != string(data) {
		t.Errorf("Expected content %q, got %q", data, content)
	}
}

func TestWriteFileIfNotExists_FileExists(t *testing.T) {
	conn := Connection{}
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	// Pre-create the file
	if err := os.WriteFile(filePath, []byte("existing"), 0666); err != nil {
		t.Fatalf("Failed to pre-create file: %v", err)
	}

	// Try writing again
	err := conn.writeFileIfNotExists(filePath, []byte("new data"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Content should not have changed
	content, _ := os.ReadFile(filePath)
	if string(content) != "existing" {
		t.Errorf("File content was overwritten. Got %q", content)
	}
}
