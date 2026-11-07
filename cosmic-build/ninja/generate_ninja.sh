#!/usr/bin/env bash
#
# generate_ninja.sh — Generador del grafo de compilación build.ninja para CosmicOS
#
# 100% Bash puro, sin dependencias de Python ni runtimes ajenos.
#

set -euo pipefail

PRODUCT="${1:-phone}"
CHANNEL="${2:-nebula}"
TARGET_ARCH="${3:-arm64}"
ROOT_DIR="${4:-.}"
BUILD_DIR="${5:-out/${TARGET_ARCH}-${PRODUCT}}"

mkdir -p "$BUILD_DIR"
NINJA_FILE="$BUILD_DIR/build.ninja"

CROSS_PREFIX=""
if [ "$TARGET_ARCH" = "arm64" ]; then
    CROSS_PREFIX="aarch64-linux-gnu-"
elif [ "$TARGET_ARCH" = "riscv64" ]; then
    CROSS_PREFIX="riscv64-linux-gnu-"
fi

JOBS="$(nproc)"

cat <<EOF > "$NINJA_FILE"
# Grafo de compilación generado por cosmic-build
# Producto: $PRODUCT | Canal: $CHANNEL | Objetivo: $TARGET_ARCH

ninja_required_version = 1.10

include $ROOT_DIR/tools/cosmic-build/ninja/rules.ninja

jobs = $JOBS
arch = $TARGET_ARCH
cross = $CROSS_PREFIX
root_dir = $ROOT_DIR
out_dir = $BUILD_DIR

# Marcador phony para forzar que Cargo/Go verifiquen incrementalmente
build ALWAYS: phony

# --- Motor de logs C++ ---
build $BUILD_DIR/bin/logengine: cpp_binary $ROOT_DIR/tools/cosmic-build/src/log_engine/LogEngine.cpp

# --- Componentes Rust (Hito 1) ---
build $ROOT_DIR/target/release/cinitd: cargo_package | ALWAYS
  package = cinitd

build $ROOT_DIR/target/release/kernel-root: cargo_package | ALWAYS
  package = kernel-root

build $ROOT_DIR/target/release/clauncher: cargo_package | ALWAYS
  package = clauncher

build $ROOT_DIR/target/release/sps: cargo_package | ALWAYS
  package = sps

build $ROOT_DIR/target/release/cosmic-repo: cargo_package | ALWAYS
  package = cosmic-repo

# --- Orquestador en Go ---
build $BUILD_DIR/bin/corchd: go_build | ALWAYS
  source_dir = $ROOT_DIR/core/corchd

# --- Kernel Linux ---
build $ROOT_DIR/kernel/linux/.config: kernel_config
  kernel_dir = $ROOT_DIR/kernel/linux
  config_file = $ROOT_DIR/kernel/configs/cosmicos_defconfig_$TARGET_ARCH

build $BUILD_DIR/boot/Image.gz: kernel_image $ROOT_DIR/kernel/linux/.config
  kernel_dir = $ROOT_DIR/kernel/linux

# --- Metas y Agrupaciones ---
build milestone1: phony \$
  $ROOT_DIR/target/release/cinitd \$
  $ROOT_DIR/target/release/kernel-root \$
  $ROOT_DIR/target/release/clauncher \$
  $ROOT_DIR/target/release/sps \$
  $ROOT_DIR/target/release/cosmic-repo \$
  $BUILD_DIR/bin/corchd \$
  $BUILD_DIR/bin/logengine

build kernel: phony $BUILD_DIR/boot/Image.gz

build all: phony milestone1

default milestone1
EOF
