# cosmic-build — Sistema Unificado de Compilación de CosmicOS

`cosmic-build` es la herramienta de compilación y orquestación del ecosistema CosmicOS.

---

## 1. Principios de Diseño

1. **Frontend en Bash:** CLI ligero y determinista (`cosmic-build`) sin dependencias de runtimes pesados ni Python.
2. **Motor Ninja:** Ejecución paralela ultra-rápida del grafo de dependencias mediante Ninja (escrito en C++).
3. **Regla de oro:** **Python NO debe estar en la imagen compilada ni en el runtime de CosmicOS**.
4. **Eficiencia y aceleración:**
   - Detección automática de `ccache` para C y C++.
   - Soporte para `sccache` en Rust.
   - Enlazador rápido `lld` o `mold`.
   - `target/` compartido para evitar recompilaciones redundantes en crates de Rust.
5. **Herramientas auxiliares en C++:** Para operaciones críticas (como verificación de artefactos con `cosmic-verifier`), se emplea C++20 nativo sin dependencias externas.

---

## 2. Uso

```bash
# Compilar los componentes del Hito 1 (por defecto)
./cosmic-build milestone1

# Compilar para arquitectura RISC-V 64
./cosmic-build --target riscv64 milestone1

# Compilar kernel Linux LTS con defconfig endurecido
./cosmic-build --product phone kernel

# Verificar artefactos generados en out/
./cosmic-build check

# Limpiar directorio out/ y target/
./cosmic-build clean
```

---

## 3. Estructura del Repositorio

- **`bin/cosmic-build`:** Frontend principal en Bash.
- **`ninja/rules.ninja`:** Reglas Ninja para Cargo, Go, CMake y Kernel.
- **`ninja/generate_ninja.sh`:** Generador del archivo `build.ninja` en Bash puro.
- **`src/verifier.cpp`:** Verificador de integridad de artefactos en C++20.
- **`Makefile`:** Compilación del verificador C++ con `ccache` y `lld`.
