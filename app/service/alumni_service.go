package service

import (
    "database/sql"
    "os"
    "github.com/gofiber/fiber/v2"
    "tugasminggu3/app/repository"
)

func CheckAlumniService(c *fiber.Ctx, db *sql.DB) error {
    key := c.Params("key")
    if key != os.Getenv("API_KEY") {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "API Key tidak valid",
            "success": false,
        })
    }

    // Coba ambil nim dari form body dulu
    nim := c.FormValue("nim")

    // Kalau kosong, ambil dari URL param
    if nim == "" {
        nim = c.Params("nim")
    }

    if nim == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "NIM wajib diisi",
            "success": false,
        })
    }
    
    alumni, err := repository.CheckAlumniByNim(db, nim)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(fiber.StatusOK).JSON(fiber.Map{
                "message": "Mahasiswa bukan alumni",
                "success": true,
                "isAlumni": false,
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal cek alumni karena " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Ini get salah satu aja",
        "success": true,
        "isAlumni": true,
        "alumni": alumni,
    })
}

func GetAllAlumniService(c *fiber.Ctx, db *sql.DB) error {
    key := c.Params("key")
    if key != os.Getenv("API_KEY") {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "message": "API Key tidak valid",
            "success": false,
        })
    }

    alumni, err := repository.GetAllAlumni(db)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(fiber.StatusOK).JSON(fiber.Map{
                "message": "Tidak ada data alumni",
                "success": true,
                "alumni": nil,
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal mendapatkan data alumni karena " + err.Error(),
            "success": false,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Ini get semua alumni",
        "success": true,
        "alumni": alumni,
    })
}