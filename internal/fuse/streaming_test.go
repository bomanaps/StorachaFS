package fuse

import (
	"io"
	"testing"
)

// TestHttpStreamReader tests the basic functionality of the streaming reader
func TestHttpStreamReader(t *testing.T) {
	// This is a simple test that verifies the reader can be created
	// In a real test environment, you would need a test HTTP server
	// that supports Range requests

	url := "https://httpbin.org/bytes/1024" // 1KB test file
	size := int64(1024)
	debug := true

	reader := newHttpStreamReader(url, size, debug)

	// Test initial state
	if reader.offset != 0 {
		t.Errorf("Expected initial offset to be 0, got %d", reader.offset)
	}

	if reader.size != size {
		t.Errorf("Expected size to be %d, got %d", size, reader.size)
	}

	// Test seeking
	newPos, err := reader.Seek(100, io.SeekStart)
	if err != nil {
		t.Errorf("Seek failed: %v", err)
	}

	if newPos != 100 {
		t.Errorf("Expected seek position to be 100, got %d", newPos)
	}

	// Test reading a small buffer
	buffer := make([]byte, 50)
	n, err := reader.Read(buffer)

	// Note: This test might fail in CI/CD environments without internet access
	if err != nil && err != io.EOF {
		t.Logf("Read failed (expected in test environment): %v", err)
		t.Logf("This is expected if the test environment doesn't have internet access")
	} else {
		if n != 50 {
			t.Errorf("Expected to read 50 bytes, got %d", n)
		}

		// Test stats
		stats := reader.GetStats()
		if stats["current_offset"].(int64) != 150 { // 100 + 50
			t.Errorf("Expected offset to be 150 after read, got %v", stats["current_offset"])
		}
	}
}

// TestHttpStreamReaderSeek tests various seek operations
func TestHttpStreamReaderSeek(t *testing.T) {
	url := "https://httpbin.org/bytes/1024"
	size := int64(1024)
	reader := newHttpStreamReader(url, size, false)

	// Test SeekStart
	pos, err := reader.Seek(500, io.SeekStart)
	if err != nil {
		t.Errorf("SeekStart failed: %v", err)
	}
	if pos != 500 {
		t.Errorf("Expected position 500, got %d", pos)
	}

	// Test SeekCurrent
	pos, err = reader.Seek(100, io.SeekCurrent)
	if err != nil {
		t.Errorf("SeekCurrent failed: %v", err)
	}
	if pos != 600 {
		t.Errorf("Expected position 600, got %d", pos)
	}

	// Test SeekEnd
	pos, err = reader.Seek(-50, io.SeekEnd)
	if err != nil {
		t.Errorf("SeekEnd failed: %v", err)
	}
	if pos != 974 { // 1024 - 50
		t.Errorf("Expected position 974, got %d", pos)
	}

	// Test negative seek
	_, err = reader.Seek(-100, io.SeekStart)
	if err == nil {
		t.Error("Expected error for negative seek, got nil")
	}

	// Test seek beyond file size
	pos, err = reader.Seek(2000, io.SeekStart)
	if err != nil {
		t.Errorf("Seek beyond file size failed: %v", err)
	}
	if pos != size {
		t.Errorf("Expected position to be clamped to file size %d, got %d", size, pos)
	}
}

// BenchmarkHttpStreamReader benchmarks the streaming reader performance
func BenchmarkHttpStreamReader(b *testing.B) {
	url := "https://httpbin.org/bytes/1024"
	size := int64(1024)
	reader := newHttpStreamReader(url, size, false)

	buffer := make([]byte, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset to beginning
		reader.Seek(0, io.SeekStart)

		// Read the entire file
		for {
			n, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Skipf("Skipping benchmark due to network error: %v", err)
				return
			}
			_ = n // Use the result to avoid optimization
		}
	}
}
