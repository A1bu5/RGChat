use aes_gcm::aead::{Aead, KeyInit, OsRng};
use aes_gcm::{Aes256Gcm, Nonce}; // Or `Aes128Gcm`
use aes_gcm::aead::generic_array::GenericArray;
use rand::Rng;

pub fn encrypt(key: &[u8], data: &[u8]) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
    let key = GenericArray::from_slice(key);
    let cipher = Aes256Gcm::new(key);

    let mut nonce = [0u8; 12]; // GCM模式中Nonce的大小通常为12字节
    rand::thread_rng().fill(&mut nonce);

    let nonce = Nonce::from_slice(&nonce); // Nonce需要转换为slice类型
    let ciphertext = cipher.encrypt(nonce, data).map_err(|e| format!("Encryption error: {:?}", e))?;

    let mut result = nonce.to_vec();
    result.extend(ciphertext);
    Ok(result)
}

pub fn decrypt(key: &[u8], data: &[u8]) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
    let key = GenericArray::from_slice(key);
    let cipher = Aes256Gcm::new(key);

    let (nonce, ciphertext) = data.split_at(12); // Nonce长度为12字节
    let nonce = Nonce::from_slice(nonce);
    let plaintext = cipher.decrypt(nonce, ciphertext).map_err(|e| format!("Decryption error: {:?}", e))?;

    Ok(plaintext)
}