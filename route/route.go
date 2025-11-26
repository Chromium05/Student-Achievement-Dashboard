package route

import (
	"student-report/app/service"
	"student-report/middleware"
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	// Homepage (GET 127.0.0.1:3000)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Halo trainer!")
	})

	// Untuk routing login
	app.Post("/login", func(c *fiber.Ctx) error {
		return service.LoginService(c, db)
	})

	// Implementasi Middleware
	protected := app.Group("", middleware.AuthRequired())

	// Untuk routing alumni
	alumni := protected.Group("/:key/alumni")

	// Menampilkan semua data alumni (GET 127.0.0.1:3000/access/alumni)
	alumni.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllAlumniService(c, db)
	})

	// Menampilkan alumni berdasarkan ID (GET 127.0.0.1:3000/access/alumni/:id)
	alumni.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetAlumniByIDService(c, db)
	})

	// Menambahkan data baru (POST 127.0.0.1:3000/access/alumni)
	alumni.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.PostNewAlumniService(c, db)
	})

	// Update data yang sudah ada (PUT 127.0.0.1:3000/access/alumni/:id)
	alumni.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdateAlumniService(c, db)
	})

	// Hapus data yang sudah ada (DELETE 127.0.0.1:3000/access/alumni/:id)
	alumni.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeleteAlumniService(c, db)
	})

	// Untuk routing pekerjaan alumni
	pekerjaan := protected.Group("/:key/pekerjaan")

	// Menampilkan semua data pekerjaan alumni (GET	127.0.0:3000/access/pekerjaan)
	pekerjaan.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllPekerjaanAlumniService(c, db)
	})

	// Menampilkan pekerjaan alumni berdasarkan ID (GET 127.0.0:3000/access/pekerjaan/:id)
	pekerjaan.Get("/:id", func(c *fiber.Ctx) error {
		return service.GetPekerjaanByIDService(c, db)
	})

	// Menampilkan pekerjaan alumni berdasarkan Alumni_ID (GET 127.0.0:3000/access/pekerjaan/alumni/:alumni_id)
	pekerjaan.Get("/alumni/:alumni_id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.GetPekerjaanByAlumniIDService(c, db)
	})

	// Menambahkan data pekerjaan alumni (POST 127.0.0:3000/access/pekerjaan)
	pekerjaan.Post("/", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.PostNewPekerjaanAlumniService(c, db)
	})

	// Update data pekerjaan alumni yang sudah ada (PUT 127.0.0:3000/access/pekerjaan/:id)
	pekerjaan.Put("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.UpdatePekerjaanAlumniService(c, db)
	})

	// SOFT DELETE pekerjaan alumni
	pekerjaan.Delete("/sofdel/:id", middleware.AuthorizePekerjaanAlumniOwnership(db), func(c *fiber.Ctx) error {
		return service.SoftDeletePekerjaanAlumniService(c, db)
	})

	// HARD DELETE data pekerjaan alumni yang sudah ada (DELETE 127.0.0:3000/access/pekerjaan/:id)
	pekerjaan.Delete("/:id", middleware.AdminOnly(), func(c *fiber.Ctx) error {
		return service.DeletePekerjaanAlumniService(c, db)
	})

	employed := protected.Group("/:key/employed")

	// Ambil view employed alumni
	employed.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllEmployedAlumniService(c, db)
	})

	// Ambil view employed alumni kurang dari 3 tahun
	employed.Get("/3years", func(c *fiber.Ctx) error {
		return service.GetEmployedAlumniLessThreeYearsService(c, db)
	})

	// Untuk routing trash
	trash := protected.Group("/:key/trash")

	// Menampilkan semua data pekerjaan alumni yang dihapus (GET 127.0.0.1:3000/access/trash)
	trash.Get("/", func(c *fiber.Ctx) error {
		return service.GetAllTrashService(c, db)
	})

	// Restore data dari trash (PUT 127.0.0.1:3000/access/trash/restore/:id)
	trash.Put("/restore/:id", middleware.RestoreOwnData(db), func(c *fiber.Ctx) error {
		return service.RestoreDataService(c, db)
	})

	// Hapus data dari trash permanen (DELETE http://127.0.0.1:3000/access/trash/delete/:id)
	trash.Delete("/delete/:id", middleware.RestoreOwnData(db), func(c *fiber.Ctx) error {
		return service.PermanentDeleteService(c, db)
	})
}
