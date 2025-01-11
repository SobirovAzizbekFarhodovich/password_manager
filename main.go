package main

import (
    "bot/api"
    "bot/api/handler"
    "bot/config"
    "bot/service"
    "bot/storage"
    "log"
)

func main() {
    cfg := config.Load()
    db, err := storage.ConnectDB(cfg)
    if err != nil {
        log.Fatalf("error while connect db, err: %s", err.Error())
    }
    services := service.InitServices(db)
    handler := handler.NewHTTPHandler(services)
    engine := api.NewGin(handler)
    err = engine.Run(cfg.HTTPPort)
    if err != nil {
        log.Fatalf("error while run engine, err: %s", err.Error())
    }
}
