package main

import (
    "log"
    "os"
    "tugasminggu3/config"
    "tugasminggu3/database"
)

func main() {
    config.LoadEnv()
    db := database.ConnectDB()
    app := config.NewApp(db)

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "3000"
    }

    log.Fatal(app.Listen(":" + port))
}
