mod encryption;
mod file_transfer;

use std::io::{self, Write, Read};
use std::net::TcpStream;
use crate::file_transfer::send_message_or_file;

fn main() -> io::Result<()> {
    let key = b"an example very very secret key.";
    let mut stream = TcpStream::connect("127.0.0.1:7878")?;
    println!("Connected to server!");

    loop {
        println!("Enter a message or the path to the file you want to send (type 'exit' to quit):");
        let mut input = String::new();
        io::stdin().read_line(&mut input)?;
        let input = input.trim();

        if input == "exit" {
            break;
        }

        send_message_or_file(&mut stream, key, input)?;

        let mut buffer = [0; 512];
        let n = stream.read(&mut buffer)?;
        println!("Server: {}", String::from_utf8_lossy(&buffer[..n]));
    }

    Ok(())
}
