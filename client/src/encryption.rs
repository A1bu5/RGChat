// encryption.rs

use aes_gcm::aead::{Aead, KeyInit, OsRng};
use aes_gcm::{Aes256Gcm, Nonce};
use aes_gcm::aead::generic_array::GenericArray;
use rand::Rng;
use std::error::Error;

// encryption.rs
pub fn encrypt(key: &[u8], data: &[u8]) -> Result<Vec<u8>, Box<dyn std::error::Error>> {
    let key = GenericArray::from_slice(key);
    let cipher = Aes256Gcm::new(key);

    let mut nonce = [0u8; 12];
    rand::thread_rng().fill(&mut nonce);

    let nonce = Nonce::from_slice(&nonce);
    let ciphertext = cipher.encrypt(nonce, data).map_err(|e| format!("Encryption error: {:?}", e))?;

    println!("Nonce: {:?}", nonce);
    println!("Ciphertext: {:?}", ciphertext);

    let mut result = nonce.to_vec();
    result.extend(ciphertext);
    Ok(result)
}

