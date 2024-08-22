// file_transfer.go
package main

import (
	"fmt"
	"net"
	"os"
)

func receiveFile(conn net.Conn, key []byte) error {
	buffer := make([]byte, 1024) // Custom the datacache size
	var fileData []byte
	var fileName string

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return err
		}

		if string(buffer[:n]) == "END_OF_FILE" {
			break
		}

		if fileName == "" {
			decryptedFileName, err := decrypt(key, buffer[:n])
			if err != nil {
				fmt.Printf("Failed to decrypt file name: %v\n", err)
				continue
			}
			fileName = string(decryptedFileName)
			fmt.Printf("Decrypted file name: %s\n", fileName)
		} else {
			decryptedData, err := decrypt(key, buffer[:n])
			if err != nil {
				fmt.Printf("Failed to decrypt file chunk: %v\n", err)
				continue
			}
			fileData = append(fileData, decryptedData...)
		}
	}

	err := os.WriteFile(fileName, fileData, 0644)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	fmt.Printf("File %s received successfully\n", fileName)
	return nil
}
