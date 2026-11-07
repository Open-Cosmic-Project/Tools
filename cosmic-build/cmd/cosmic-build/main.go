package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Colores ANSI para presentación estética
const (
	ColorReset   = "\033[0m"
	ColorBold    = "\033[1m"
	ColorCyan    = "\033[36m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorRed     = "\033[31m"
	ColorMagenta = "\033[35m"
	ColorGray    = "\033[90m"
)

type Config struct {
	Product   string
	Channel   string
	Target    string
	Action    string
	Jobs      int
	Verbose   bool
	RootDir   string
	OutDir    string
	ToolsDir  string
}

func printBanner() {
	fmt.Println(ColorCyan + ColorBold + `
   ____               _        ____  ____ 
  / ___|___  ___ _ __(_) ___  / __ \/ ___|
 | |   / _ \/ __| '_ \ |/ __|/ / _  \___ \
 | |__| (_) \__ \ | | | | (__| | (_) |___) |
  \____\___/|___/_| |_|_|\___|\ \___/|____/ 
                               \____/      ` + ColorReset)
	fmt.Println(ColorGray + "  Sistema Unificado de Compilación — OpenCosmicProject (OCP)" + ColorReset)
	fmt.Println()
}

func main() {
	var cfg Config

	flag.StringVar(&cfg.Product, "product", "phone", "Producto a compilar: phone, desktop")
	flag.StringVar(&cfg.Channel, "channel", "nebula", "Canal de estabilidad: nebula, pulsar, singularity, quasar, eventhorizon")
	flag.StringVar(&cfg.Target, "target", "arm64", "Arquitectura objetivo: arm64, riscv64")
	flag.IntVar(&cfg.Jobs, "j", runtime.NumCPU(), "Número de trabajos paralelos")
	flag.BoolVar(&cfg.Verbose, "v", false, "Salida verbosa de compilación")

	flag.Usage = func() {
		printBanner()
		fmt.Printf("%sUso:%s ./cosmic-build [opciones] [meta]\n\n", ColorBold, ColorReset)
		fmt.Printf("%sMetas (targets):%s\n", ColorBold, ColorReset)
		fmt.Printf("  %smilestone1%s      Compila crates del Hito 1 (Rust + Go) [Por defecto]\n", ColorGreen, ColorReset)
		fmt.Printf("  %skernel%s          Compila el kernel Linux LTS con defconfig endurecido\n", ColorGreen, ColorReset)
		fmt.Printf("  %sall%s             Compila la imagen completa del sistema\n", ColorGreen, ColorReset)
		fmt.Printf("  %scheck%s           Verifica integridad y sumas SHA-256 de los artefactos\n", ColorGreen, ColorReset)
		fmt.Printf("  %sclean%s           Limpia directorios out/ y target/\n\n", ColorGreen, ColorReset)
		fmt.Printf("%sOpciones:%s\n", ColorBold, ColorReset)
		flag.PrintDefaults()
		fmt.Println()
	}

	flag.Parse()

	// Meta por defecto o argumento posicional
	cfg.Action = "milestone1"
	if flag.NArg() > 0 {
		cfg.Action = flag.Arg(0)
	}

	// Localizar la raíz del workspace
	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError localizando ejecutable:%s %v\n", ColorRed, ColorReset, err)
		os.Exit(1)
	}
	exeDir := filepath.Dir(exePath)
	
	// Buscar la raíz recorriendo hacia arriba hasta encontrar Cargo.toml o .git
	rootDir := exeDir
	for {
		if _, err := os.Stat(filepath.Join(rootDir, "cosmic-manifest")); err == nil {
			break
		}
		parent := filepath.Dir(rootDir)
		if parent == rootDir {
			rootDir, _ = os.Getwd()
			break
		}
		rootDir = parent
	}

	cfg.RootDir = rootDir
	cfg.ToolsDir = filepath.Join(rootDir, "tools", "cosmic-build")
	cfg.OutDir = filepath.Join(rootDir, "out", fmt.Sprintf("%s-%s", cfg.Target, cfg.Product))

	if cfg.Action == "clean" {
		doClean(&cfg)
		return
	}

	printBanner()
	printStatus(&cfg)

	// Ejecutar acción
	startTime := time.Now()
	switch cfg.Action {
	case "milestone1", "all", "kernel":
		err = runNinjaBuild(&cfg)
	case "check":
		err = runCheckArtifacts(&cfg)
	default:
		fmt.Fprintf(os.Stderr, "%sMeta desconocida:%s %s\n", ColorRed, ColorReset, cfg.Action)
		flag.Usage()
		os.Exit(1)
	}

	elapsed := time.Since(startTime)
	if err != nil {
		fmt.Printf("\n%s✗ Falló la ejecución de '%s'%s (tiempo: %v)\n", ColorRed+ColorBold, cfg.Action, ColorReset, elapsed.Round(time.Millisecond))
		os.Exit(1)
	}

	fmt.Printf("\n%s✓ Meta '%s' completada con éxito%s en %v\n\n", ColorGreen+ColorBold, cfg.Action, ColorReset, elapsed.Round(time.Millisecond))
}

func printStatus(cfg *Config) {
	fmt.Printf(" %s⚙ Configuración:%s\n", ColorBold, ColorReset)
	fmt.Printf("   • Producto:     %s%s%s\n", ColorCyan, cfg.Product, ColorReset)
	fmt.Printf("   • Canal:        %s%s%s\n", ColorMagenta, cfg.Channel, ColorReset)
	fmt.Printf("   • Arquitectura: %s%s%s\n", ColorYellow, cfg.Target, ColorReset)
	fmt.Printf("   • Meta:         %s%s%s\n", ColorGreen, cfg.Action, ColorReset)
	fmt.Printf("   • Hilos (jobs): %d\n", cfg.Jobs)
	fmt.Printf("   • Salida:       %s\n", cfg.OutDir)

	// Detección de herramientas de aceleración
	var accels []string
	if _, err := exec.LookPath("ccache"); err == nil {
		accels = append(accels, "ccache (C/C++)")
	}
	if _, err := exec.LookPath("sccache"); err == nil {
		accels = append(accels, "sccache (Rust)")
	}
	if _, err := exec.LookPath("lld"); err == nil {
		accels = append(accels, "lld (Linker)")
	}
	if len(accels) > 0 {
		fmt.Printf("   • Aceleración:  %s%s%s\n", ColorGreen, strings.Join(accels, ", "), ColorReset)
	}
	fmt.Println()
}

func doClean(cfg *Config) {
	fmt.Printf("==> Limpiando directorio de artefactos: %s y target/...\n", cfg.OutDir)
	_ = os.RemoveAll(cfg.OutDir)
	_ = os.RemoveAll(filepath.Join(cfg.RootDir, "target"))
	_ = os.Remove(filepath.Join(cfg.RootDir, ".ninja_log"))
	fmt.Println(ColorGreen + "✓ Limpieza completada con éxito." + ColorReset)
}

func runNinjaBuild(cfg *Config) error {
	// 1. Asegurar generación de build.ninja
	genScript := filepath.Join(cfg.ToolsDir, "ninja", "generate_ninja.sh")
	genCmd := exec.Command(genScript, cfg.Product, cfg.Channel, cfg.Target, cfg.RootDir, cfg.OutDir)
	genCmd.Stdout = os.Stdout
	genCmd.Stderr = os.Stderr
	if err := genCmd.Run(); err != nil {
		return fmt.Errorf("error generando build.ninja: %w", err)
	}

	// 2. Invocar Ninja
	ninjaFile := filepath.Join(cfg.OutDir, "build.ninja")
	args := []string{"-f", ninjaFile, "-j", fmt.Sprintf("%d", cfg.Jobs)}
	if cfg.Verbose {
		args = append(args, "-v")
	}
	args = append(args, cfg.Action)

	fmt.Printf("==> Ejecutando motor Ninja para la meta '%s'...\n", cfg.Action)
	ninjaCmd := exec.Command("ninja", args...)
	ninjaCmd.Stdout = os.Stdout
	ninjaCmd.Stderr = os.Stderr
	ninjaCmd.Env = append(os.Environ(),
		fmt.Sprintf("CARGO_TARGET_DIR=%s", filepath.Join(cfg.RootDir, "target")),
	)

	return ninjaCmd.Run()
}

func runCheckArtifacts(cfg *Config) error {
	fmt.Printf("==> Verificando artefactos compilados en %s...\n", cfg.OutDir)

	if _, err := os.Stat(cfg.OutDir); os.IsNotExist(err) {
		return fmt.Errorf("directorio de salida no existe. Ejecuta './cosmic-build milestone1' primero")
	}

	var count int
	err := filepath.Walk(cfg.OutDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		hashStr, err := calculateSHA256(path)
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(cfg.RootDir, path)
		fmt.Printf("  %s✓%s %-38s %s[%7d B]%s %s%s%s\n",
			ColorGreen, ColorReset,
			relPath,
			ColorGray, info.Size(), ColorReset,
			ColorCyan, hashStr[:16]+"...", ColorReset,
		)
		count++
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("\n==> Total de artefactos verificados: %s%d%s\n", ColorBold, count, ColorReset)
	return nil
}

func calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
