//! cosmic-flash — Herramienta de flasheo de imágenes y bundles de CosmicOS
//!
//! Alcance acordado en CosmicOS spec sección 5.2:
//! - Flasheo unificado de imágenes de disco GPT con UKI (Unified Kernel Image).
//! - Gestiona particiones A/B del sistema base EROFS y slot inactivo.

fn main() {
    println!("[cosmic-flash] Herramienta oficial de flasheo de CosmicOS");
    println!("Uso: cosmic-flash <dispositivo-destino> <imagen.img|bundle.ota>");
}
