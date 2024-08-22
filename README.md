# RGChat
##aka. Rust_GO_C Chat

A open-source chat application with end-to-end encryption, developed using Rust, Go, and C++. The application supports secure messaging and file transfer between clients and a server. The project is designed with a modular structure, allowing easy extension and integration with different UI frameworks like Qt.

## Features

- **Secure Messaging**: Messages sent between clients and server are encrypted using AES-GCM, ensuring privacy and integrity.
- **File Transfer**: Supports encrypted file transfer between clients and the server, with proper handling of file chunks and end-to-end encryption.
- **Modular Design**: The project is structured into multiple modules, making it easy to manage and extend different functionalities.
- **Cross-language Integration**: The core logic is implemented in Rust and Go, with the potential for integration with C++ based GUI using Qt.

## Project Structure

```plaintext
.
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
