// file_transfer.rs
use std::fs::File;
use std::io::{self, Read, Write};
use std::net::TcpStream;
use crate::encryption::encrypt;
use aes_gcm::aead::generic_array::GenericArray;
use aes_gcm::Aes256Gcm;
use rand::RngCore;
use aes_gcm::aead::{Aead, KeyInit};

pub fn read_file(file_path: &str) -> io::Result<(Vec<u8>, String)> {
    let mut file = File::open(file_path)?;
    let mut buffer = Vec::new();
    file.read_to_end(&mut buffer)?;

    let file_name = file_path.split('/').last().unwrap_or("unknown").to_string();
    Ok((buffer, file_name))
}


pub fn send_message_or_file(stream: &mut TcpStream, key: &[u8], input: &str) -> io::Result<()> {
    if input.starts_with("/") {
        // File
        send_file(stream, key, input)
    } else {
        // Message
        let encrypted_message = encrypt(key, input.as_bytes()).expect("Message encryption failed");
        stream.write_all(b"MSG\n")?;
        stream.write_all(&encrypted_message)?;
        stream.write_all(b"\n")?;
        Ok(())
    }
}

pub fn send_file(stream: &mut TcpStream, key: &[u8], file_path: &str) -> io::Result<()> {
    let mut file = File::open(file_path)?;
    let mut buffer = Vec::new();
    file.read_to_end(&mut buffer)?;
    let file_name = file_path.split('/').last().unwrap_or("unknown").to_string();
    let file_size = buffer.len();

    // metadataSend
    let metadata = format!("FILENAME:{}\nSIZE:{}\n", file_name, file_size);
    println!("Sending metadata: {}", metadata);
    stream.write_all(metadata.as_bytes())?;
    println!("Metadata sent.");

    // Encrypt
    let cipher = Aes256Gcm::new(GenericArray::from_slice(key));

    // Divide chunks
    let chunk_size = 1024;

    for (i, chunk) in buffer.chunks(chunk_size).enumerate() {
        println!("Encrypting chunk {}", i);
        let mut nonce = [0u8; 12];
        rand::thread_rng().fill_bytes(&mut nonce);
        let nonce = GenericArray::from_slice(&nonce);

        let encrypted_chunk = cipher.encrypt(nonce, chunk)
            .expect("File chunk encryption failed");

        // Sent Nonce and EncryptedChunks
        stream.write_all(nonce)?;
        stream.write_all(&encrypted_chunk)?;
        println!("Chunk {} sent. Nonce: {:?}", i, nonce);
    }

    // Transfer End sign
    println!("Sending end signal.");
    stream.write_all(b"END\n")?;
    println!("File transmission completed.");
    Ok(())
}





// file_transfer.rs





