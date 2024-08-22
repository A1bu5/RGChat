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
	key := []byte("an example very very secret key.")

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
			// 处理普通消息
			handleIncomingMessage(conn, reader, key)
		} else if strings.HasPrefix(line, "FILENAME:") {
			// 处理文件传输
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
	conn.Write([]byte(response)) // 注意这里的 conn 变量
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

	fileData := make([]byte, 0, fileSize) // 动态扩展的缓冲区
	receivedBytes := 0

	for {
		fmt.Printf("Reading chunk %d\n", receivedBytes)
		// 读取 Nonce
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

		// 读取加密的数据块
		// 注意这里的大小要根据实际收到的数据调整
		chunkSize := fileSize - receivedBytes
		if chunkSize > 1024 {
			chunkSize = 1024 + 16 // 16 是 GCM 模式的认证标签长度
		} else {
			chunkSize += 16 // 最后一块数据也包含认证标签
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
