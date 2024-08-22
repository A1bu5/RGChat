package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	key := []byte("an example very very secret key.")//Please makesuret the key is same to the clients

	listener, err := net.Listen("tcp", ":7878")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 7878...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, key)
	}
}

func handleConnection(conn net.Conn, key []byte) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected.")
			return
		}

		line = strings.TrimSpace(line)
		if line == "MSG" {
			// Message
			handleIncomingMessage(conn, reader, key)
		} else if strings.HasPrefix(line, "FILENAME:") {
			// File
			handleIncomingFile(conn, reader, key, line)
		}
	}
}

func handleIncomingMessage(conn net.Conn, reader *bufio.Reader, key []byte) {
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}
	data := strings.TrimSpace(line)

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Failed to create cipher:", err)
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("Failed to create GCM:", err)
		return
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		fmt.Println("Failed to decrypt message:", err)
		return
	}

	fmt.Printf("Decrypted message: %s\n", string(plaintext))
	response := strings.ToUpper(string(plaintext))
	conn.Write([]byte(response)) 
}

func handleIncomingFile(conn net.Conn, reader *bufio.Reader, key []byte, metadata string) {
	fmt.Println("Received metadata:", metadata)
	fileName := strings.TrimPrefix(metadata, "FILENAME:")
	fileSizeLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading file size:", err)
		return
	}
	fileSizeLine = strings.TrimSpace(fileSizeLine)
	fmt.Println("Received file size:", fileSizeLine)
	fileSize, err := strconv.Atoi(strings.TrimPrefix(fileSizeLine, "SIZE:"))
	if err != nil {
		fmt.Println("Invalid file size:", err)
		return
	}

	fileData := make([]byte, 0, fileSize) // cache
	receivedBytes := 0

	for {
		fmt.Printf("Reading chunk %d\n", receivedBytes)
		//Read  Nonce
		nonce := make([]byte, 12)
		_, err := io.ReadFull(reader, nonce)
		if err == io.EOF {
			fmt.Println("File transmission ended unexpectedly.")
			break
		}
		if err != nil {
			fmt.Println("Failed to read nonce:", err)
			return
		}
		fmt.Printf("Received Nonce: %x\n", nonce)

		// Read EncryptedChunks
		// DataAdjust is important
		chunkSize := fileSize - receivedBytes
		if chunkSize > 1024 {
			chunkSize = 1024 + 16 // 16 is GCM isgn length
		} else {
			chunkSize += 16 // Check the last chunk sign
		}
		encryptedChunk := make([]byte, chunkSize)
		n, err := io.ReadFull(reader, encryptedChunk)
		if err == io.EOF {
			fmt.Println("File transmission ended unexpectedly.")
			break
		}
		if err != nil {
			fmt.Println("Failed to read file chunk:", err)
			return
		}

		fmt.Println("Decrypting chunk")
		block, err := aes.NewCipher(key)
		if err != nil {
			fmt.Println("Failed to create cipher:", err)
			return
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			fmt.Println("Failed to create GCM:", err)
			return
		}

		decryptedChunk, err := gcm.Open(nil, nonce, encryptedChunk[:n], nil)
		if err != nil {
			fmt.Printf("Failed to decrypt file chunk: %v\n", err)
			return
		}

		fileData = append(fileData, decryptedChunk...)
		receivedBytes += len(decryptedChunk)

		if receivedBytes >= fileSize {
			fmt.Println("File received completely.")
			break
		}
	}

	err = os.WriteFile(fileName, fileData, 0644)
	if err != nil {
		fmt.Println("Failed to save file:", err)
		return
	}

	fmt.Printf("File %s received successfully\n", fileName)
	conn.Write([]byte("File received successfully"))
}
