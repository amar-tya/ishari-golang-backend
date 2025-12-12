package main

import (
    "fmt"
    "log"

    "ishari-backend/internal/bootstrap"
    "ishari-backend/pkg/config"
)

func main() {
    cfg := config.Load()

    app, err := bootstrap.Build(cfg)
    if err != nil {
        log.Fatalf("bootstrap error: %v", err)
    }
    defer func() {
        if err := app.Cleanup(); err != nil {
            log.Printf("cleanup error: %v", err)
        }
    }()

    addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
    log.Printf("Server starting on %s", addr)
    if err := app.Server.Start(addr); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
