package route

import (
    "database/sql"
    "github.com/gofiber/fiber/v2"
    "tugasminggu3/app/service"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Selamat, Berjuang, Sukses!!!!")
    })

    app.Get("/:key", func(c *fiber.Ctx) error {
        return service.GetAllAlumniService(c, db)
    })

    app.Get("/check/:key", func(c *fiber.Ctx) error {
        return service.CheckAlumniService(c, db)
    })

	app.Get("/check/:key/:nim", func(c *fiber.Ctx) error {
        nim := c.Params("nim")
        c.Request().PostArgs().Set("nim", nim)
        return service.CheckAlumniService(c, db)
    })
}
