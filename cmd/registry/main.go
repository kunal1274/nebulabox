package main

import (
    "log"

    "github.com/nebulabox/nebulabox/internal/registry"
)

func main() {
    srv := registry.NewServer()
    if err := srv.Start(); err != nil {
        log.Fatal(err)
    }
}


