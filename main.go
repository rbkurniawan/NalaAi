package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/rbkurniawan/NalaAi/config"
    "github.com/rbkurniawan/NalaAi/handlers"
    "github.com/rbkurniawan/NalaAi/utils"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize logger
    logger := utils.NewLogger()

    // Initialize handlers
    chatHandler := handlers.NewChatHandler(cfg, logger)

    // Setup routes
    http.HandleFunc("/api/chat", chatHandler.HandleChat)
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // Start server
    addr := fmt.Sprintf(":%s", cfg.ServerPort)
    log.Printf("Server starting on %s", addr)
    
    logger.Log("INFO", fmt.Sprintf("Server starting on %s", addr))
    
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatal(err)
        logger.Log("ERROR", fmt.Sprintf("Server error: %v", err))
    }
}