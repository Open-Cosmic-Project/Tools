//! cosmic-repo — Infraestructura de repositorio central y gestión de TUF
//!
//! Alcance acordado en CosmicOS spec sección 8.1, 8.2 y 8.4:
//! - Generación y firma de metadatos TUF para repositorios Main, Space y MTL.
//! - Registro de procedencia SLSA y atestación cosign.
//! - Soporte para compilaciones reproducibles e inspección de hashes.

fn main() {
    println!("[cosmic-repo] Gestor del repositorio central de CosmicOS (TUF / SLSA)");
    println!("Uso: cosmic-repo <comando> [opciones]");
}
