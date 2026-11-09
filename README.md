# Tools

Herramientas del ecosistema.

| Herramienta | Qué hace |
|---|---|
| `cosmic-build` | Sistema de construcción unificado (frontend, Ninja y registro de eventos) |
| `cosmic-repo` | Herramienta de repositorios |
| `cosmic-flash` | Flasheo del dispositivo desde el PC |
| `fastboot` | `cosmic-fastboot`: cliente fastboot parcialmente compatible con el de Android (`--oem-unlock`, `--oem-get-data`, `--oem-lock`) |
| `cad` | CosmicOS Advanced Driver: el ADB de CosmicOS |

`fastboot` y `cad` usan `libs/cosmic-usb`.

## Construir el sistema

Desde la raíz del árbol:

```sh
./cosmic-build milestone1
```

## Estado

`cosmic-build` y `cosmic-repo` forman parte del hito 1. `fastboot`, `cad` y `cosmic-flash` están aplazados.

## Licencia

GPL-3.0-or-later.
