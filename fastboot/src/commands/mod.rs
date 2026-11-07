//! Submódulos de comandos para CosmicOS Fastboot.

/// Listar dispositivos en modo fastboot.
pub mod devices;
/// Flashear una partición con una imagen.
pub mod flash;
/// Borrar una partición.
pub mod erase;
/// Bloquear el bootloader.
pub mod lock;
/// Desbloquear el bootloader.
pub mod unlock;
/// Reiniciar el dispositivo.
pub mod reboot;
/// Obtener variables del dispositivo.
pub mod getvar;
/// Ejecutar comandos OEM.
pub mod oem;
