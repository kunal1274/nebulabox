package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	var (
		packageFlag = flag.String("pkg", "", "Package to benchmark (e.g., ./internal/api)")
		allFlag     = flag.Bool("all", false, "Run all benchmarks")
		timeout     = flag.Duration("timeout", 30*time.Minute, "Benchmark timeout")
	)
	flag.Parse()
	
	if !*allFlag && *packageFlag == "" {
		fmt.Println("Usage: go run cmd/benchmark/main.go [flags]")
		fmt.Println("Flags:")
		fmt.Println("  -pkg <path>   Package to benchmark (e.g., ./internal/api)")
		fmt.Println("  -all          Run all benchmarks")
		fmt.Println("  -timeout      Timeout duration (default: 30m)")
		os.Exit(1)
	}
	
	var packages []string
	if *allFlag {
		packages = []string{
			"./internal/api",
			"./internal/containerd",
			"./internal/registry",
		}
	} else {
		packages = []string{*packageFlag}
	}
	
	fmt.Println("🚀 NebulaBox Benchmark Runner")
	fmt.Println("==============================")
	fmt.Printf("Timeout: %v\n", *timeout)
	fmt.Println()
	
	for _, pkg := range packages {
		fmt.Printf("📊 Benchmarking %s...\n", pkg)
		fmt.Println("----------------------------")
		
		cmd := exec.Command("go", "test", "-bench=.", "-benchmem", "-benchtime=3s", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Benchmark failed for %s: %v\n", pkg, err)
		} else {
			fmt.Printf("✅ Benchmark completed for %s\n", pkg)
		}
		fmt.Println()
	}
	
	fmt.Println("✅ All benchmarks completed!")
}

