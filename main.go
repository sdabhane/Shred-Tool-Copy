package main

import (
	"crypto/rand" // Importing crypto/rand package to generate cryptographic secure random numbers.
	"fmt"
	"os"
)

// Shred overwrites the file located at the provided path 3 times with random data and then deletes it.
func Shred(path string) error {
	// Open the file with write-only permission and truncate its length to 0.
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	defer file.Close() // Ensure the file gets closed after we are done.

	// Overwrite the file 3 times.
	for i := 0; i < 3; i++ {
		randomData := make([]byte, 1024) // Create a byte slice with 1024 empty bytes.

		// Fill the byte slice with random data.
		_, err := rand.Read(randomData)
		if err != nil {
			return err
		}

		// Write the random data to the file.
		_, err = file.Write(randomData)
		if err != nil {
			return err
		}
	}

	// Explicitly close the file (even though defer is in place). This is for added clarity.
	err = file.Close()
	if err != nil {
		return err
	}

	// Delete the file after overwriting.
	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Check if the user has provided a file path as an argument.
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <file_path>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	// Call the Shred function on the provided file path.
	err := Shred(filePath)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("File shredded and deleted successfully.")
	}
}
