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
        // 处理文件传输
        send_file(stream, key, input)
    } else {
        // 处理普通消息
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

    // 发送文件元数据（文件名和文件大小）
    let metadata = format!("FILENAME:{}\nSIZE:{}\n", file_name, file_size);
    println!("Sending metadata: {}", metadata);
    stream.write_all(metadata.as_bytes())?;
    println!("Metadata sent.");

    // 创建加密器
    let cipher = Aes256Gcm::new(GenericArray::from_slice(key));

    // 分块发送文件数据
    let chunk_size = 1024;

    for (i, chunk) in buffer.chunks(chunk_size).enumerate() {
        println!("Encrypting chunk {}", i);
        let mut nonce = [0u8; 12];
        rand::thread_rng().fill_bytes(&mut nonce);
        let nonce = GenericArray::from_slice(&nonce);

        let encrypted_chunk = cipher.encrypt(nonce, chunk)
            .expect("File chunk encryption failed");

        // 发送 Nonce 和加密的数据块
        stream.write_all(nonce)?;
        stream.write_all(&encrypted_chunk)?;
        println!("Chunk {} sent. Nonce: {:?}", i, nonce);
    }

    // 发送文件传输结束标志
    println!("Sending end signal.");
    stream.write_all(b"END\n")?;
    println!("File transmission completed.");
    Ok(())
}





// file_transfer.rs





