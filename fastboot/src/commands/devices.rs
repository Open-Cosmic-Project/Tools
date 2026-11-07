//! Comando `devices` — lista dispositivos en modo fastboot.

use cosmic_usb::UsbTransport;

/// Escanea y muestra los dispositivos conectados en modo fastboot.
pub fn run() -> Result<(), Box<dyn std::error::Error>> {
    println!("Buscando dispositivos en modo fastboot...");

    let dispositivos = UsbTransport::scan_fastboot_devices()?;

    if dispositivos.is_empty() {
        println!("No se encontraron dispositivos en modo fastboot.");
        return Ok(());
    }

    println!();
    println!("{:<20} {:<15} {:<20}", "SERIAL", "MODO", "PRODUCTO");
    println!("{}", "-".repeat(55));

    for dev in &dispositivos {
        println!("{:<20} {:<15} {:<20}", dev.serial, dev.mode, dev.product);
    }

    println!();
    println!(
        "Total: {} dispositivo(s) encontrado(s).",
        dispositivos.len()
    );

    Ok(())
}
