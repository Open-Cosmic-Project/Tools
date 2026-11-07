#![deny(warnings)]

//! CosmicOS Fastboot — herramienta de flasheo y gestión OEM para dispositivos CosmicOS.
//!
//! Equivalente a `fastboot` de Android, adaptado para el ecosistema CosmicOS.
//! Soporta flasheo de particiones, bloqueo/desbloqueo de bootloader,
//! gestión de particiones y comandos OEM.

mod commands;

use std::env;
use std::process;

/// Banner ASCII de CosmicOS Fastboot.
const BANNER: &str = r#"
   ___          _   _                 _   
  / __\__ _ ___| |_| |__   ___   ___ | |_ 
 / _\/ _` / __| __| '_ \ / _ \ / _ \| __|
/ / | (_| \__ \ |_| |_) | (_) | (_) | |_ 
\/   \__,_|___/\__|_.__/ \___/ \___/ \__|
  CosmicOS Fastboot v0.0.1
"#;

/// Muestra la ayuda de uso del programa.
fn mostrar_ayuda() {
    println!("{BANNER}");
    println!("Uso: cosmic-fastboot <comando> [argumentos]");
    println!();
    println!("Comandos disponibles:");
    println!("  devices                     Listar dispositivos en modo fastboot");
    println!("  flash <partición> <imagen>  Flashear una partición con una imagen");
    println!("  erase <partición>           Borrar una partición");
    println!("  lock                        Bloquear el bootloader");
    println!("  unlock                      Desbloquear el bootloader");
    println!("  reboot                      Reiniciar el dispositivo");
    println!("  reboot-bootloader           Reiniciar al bootloader");
    println!("  getvar <variable>           Obtener variable del dispositivo");
    println!("  oem <comando> [args...]     Ejecutar comando OEM");
    println!("  help                        Mostrar esta ayuda");
    println!();
    println!("Ejemplos:");
    println!("  cosmic-fastboot devices");
    println!("  cosmic-fastboot flash system system.img");
    println!("  cosmic-fastboot getvar serialno");
    println!("  cosmic-fastboot oem unlock-critical");
}

fn main() {
    env_logger::init();

    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        mostrar_ayuda();
        process::exit(0);
    }

    let comando = args[1].as_str();

    let resultado = match comando {
        "devices" => commands::devices::run(),

        "flash" => {
            if args.len() < 4 {
                eprintln!("\x1b[31mError:\x1b[0m el comando 'flash' requiere <partición> e <imagen>");
                eprintln!("Uso: cosmic-fastboot flash <partición> <imagen>");
                process::exit(1);
            }
            commands::flash::run(&args[2], &args[3])
        }

        "erase" => {
            if args.len() < 3 {
                eprintln!("\x1b[31mError:\x1b[0m el comando 'erase' requiere <partición>");
                eprintln!("Uso: cosmic-fastboot erase <partición>");
                process::exit(1);
            }
            commands::erase::run(&args[2])
        }

        "lock" => commands::lock::run(),

        "unlock" => commands::unlock::run(),

        "reboot" => commands::reboot::run(None),

        "reboot-bootloader" => commands::reboot::run(Some("bootloader")),

        "getvar" => {
            if args.len() < 3 {
                eprintln!("\x1b[31mError:\x1b[0m el comando 'getvar' requiere <variable>");
                eprintln!("Uso: cosmic-fastboot getvar <variable>");
                process::exit(1);
            }
            commands::getvar::run(&args[2])
        }

        "oem" => {
            if args.len() < 3 {
                eprintln!("\x1b[31mError:\x1b[0m el comando 'oem' requiere al menos un argumento");
                eprintln!("Uso: cosmic-fastboot oem <comando> [args...]");
                process::exit(1);
            }
            commands::oem::run(&args[2..])
        }

        "help" | "--help" | "-h" => {
            mostrar_ayuda();
            process::exit(0);
        }

        otro => {
            eprintln!("\x1b[31mError:\x1b[0m comando desconocido '{otro}'");
            eprintln!("Ejecuta 'cosmic-fastboot help' para ver los comandos disponibles.");
            process::exit(1);
        }
    };

    if let Err(e) = resultado {
        eprintln!("\x1b[31mError fatal:\x1b[0m {e}");
        process::exit(1);
    }
}
