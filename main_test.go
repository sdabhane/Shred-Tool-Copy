package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// TestShred checks if a temporary file can be successfully shredded.
func TestShred(t *testing.T) {
	// Create a temporary file for testing.
	tempFile, err := ioutil.TempFile("", "shred_test")
	if err != nil {
		t.Errorf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Ensure the temporary file gets deleted after the test.

	// Populate the temporary file with some content.
	_, err = tempFile.WriteString("This is a test file.")
	if err != nil {
		t.Errorf("Error writing to temporary file: %v", err)
	}

	// Call the Shred function to overwrite and delete the file.
	err = Shred(tempFile.Name())
	if err != nil {
		t.Errorf("Error shredding file: %v", err)
	}

	// Ensure that the file no longer exists after shredding.
	if _, err := os.Stat(tempFile.Name()); err == nil {
		t.Errorf("File still exists after shredding")
	}
	fmt.Println("Test case passed successfully to shred temporary file")
}

// TestShredFileDoesNotExist checks the behavior of Shred when provided a nonexistent file.
func TestShredFileDoesNotExist(t *testing.T) {
	err := Shred("this_file_does_not_exist") // Intentionally shredding a nonexistent file.
	if err == nil {
		t.Errorf("Expected error when shredding non-existent file")
	}
	fmt.Println("Test case passed successfully to shred non-existing file")
}

// TestShredFileNotWritable checks if Shred returns an error when trying to shred a read-only file.
func TestShredFileNotWritable(t *testing.T) {
	// Create a temporary file for testing.
	file, err := ioutil.TempFile("", "shred_test")
	if err != nil {
		t.Errorf("Error creating temporary file: %v", err)
	}
	defer os.Remove(file.Name()) // Ensure the file gets deleted after the test.

	// Make the file read-only.
	err = file.Chmod(0400)
	if err != nil {
		t.Errorf("Error setting file permissions: %v", err)
	}

	// Call the Shred function and expect an error due to file being read-only.
	err = Shred(file.Name())
	if err == nil {
		t.Errorf("Expected error when shredding file that is not writable")
	}
	fmt.Println("Test case passed successfully to try shredding non-writable file")
}

// TestShredFileOpenByAnotherProcess checks if Shred returns an error when trying to shred a file that's currently being read by another process.
func TestShredFileOpenByAnotherProcess(t *testing.T) {
	// Create a temporary file for testing.
	file, err := ioutil.TempFile("", "shred_test")
	if err != nil {
		t.Errorf("Error creating temporary file: %v", err)
	}
	defer os.Remove(file.Name()) // Ensure the file gets deleted after the test.

	// Simulate another process (or goroutine) that reads the file.
	go func() {
		file2, err := os.OpenFile(file.Name(), os.O_RDONLY, 0)
		if err != nil {
			t.Errorf("Error opening file in another process: %v", err)
		}
		defer file2.Close() // Close the file when done reading.

		// Mimic a long-running process by using an infinite loop.
		for {
		}
	}()

	// Call the Shred function and expect an error since the file is being read by another process.
	err = Shred(file.Name())
	if err == nil {
		t.Errorf("Expected error when shredding file that is open by another process")
	}
	fmt.Println("Test case passed successfully to shred already being used file")
}
