# RGChat
## aka. Rust_GO_C Chat, or you can say RipGhostChat

#
A open-source chat application with end-to-end encryption, developed using Rust, Go, and C++. The application supports secure messaging and file transfer between clients and a server. The project is designed with a modular structure, allowing easy extension and integration with different UI frameworks like Qt.

## Features
![Chat Example](./images/1)
![Chat Example](./images/2)
- **Secure Messaging**: Messages sent between clients and server are encrypted using AES-GCM, ensuring privacy and integrity.
- **File Transfer**: Supports encrypted file transfer between clients and the server, with proper handling of file chunks and end-to-end encryption.
- **Modular Design**: The project is structured into multiple modules, making it easy to manage and extend different functionalities.
- **Cross-language Integration**: The core logic is implemented in Rust and Go, with the potential for integration with C++ based GUI using Qt.

## Project Structure

```plaintext

├── client/                 # Rust-based client application
│   ├── src/
│   │   ├── main.rs         # Entry point for the client
│   │   ├── encryption.rs   # Encryption and decryption logic
│   │   └── file_transfer.rs# File transfer handling
│   └── Cargo.toml          # Rust project configuration
│
├── server/                 # Go-based server application
│   ├── main.go             # Entry point for the server
│   └── go.mod              # Go module configuration
│
├── README.md               # Project documentation
└── LICENSE                 # License file
```
# Getting Started
## Prerequisites
(AI generated, its easy to use actually)
- **Rust**: Ensure you have Rust installed. You can install it using rustup.
- **Go**: Ensure you have Go installed. You can download it from the official website.
- **C++**: Sorry, GUI is still in developing
## Building the Project
- **Client (Rust)**
  1.Navigate to the client directory:

		cd client

	2.Build the Rust client:

		cargo build --release
- **Server (Go)**
  1.Navigate to the server directory:

		cd server

	2.Run the Go server:

		go run main.go

- **Run Client**

		cd client


		cargo run
# USAGE
## 1.Running the Server:
Start the server first using the Go command mentioned above. The server listens for incoming connections from clients.
## 2.Running the Client:
After the server is running, start the client. You can send messages or transfer files using the command-line interface. The client will handle encryption before sending data to the server.
## 3.File Transfer:
To send a file, input the file path when prompted. The client will read, encrypt, and send the file to the server in chunks.

# Future Enhancements
## TOO MUCH
## GUI Integration: 
The next step involves integrating the command-line based logic with a C++ GUI using Qt(Maybe not, javascript is much easy to use i think how ever i like C++ as i love u >_<), enabling a richer user experience.
## Group Chats & Channels:
Implement group chat and channel functionalities similar to Telegram and WeChat.
## Extra End-to-End Encryption: (RingSignature)
Maybe
## Contributing
Contributions are welcome! Please fork this repository and submit a pull request with your changes.
