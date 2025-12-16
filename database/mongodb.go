package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConn() *mongo.Database {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Pastikan MongoDB Anda berjalan di sini
		log.Println("Peringatan: MONGO_URI tidak disetel. Menggunakan default:", mongoURI)
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	// Membuat konteks dengan batas waktu
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Koneksi ke MongoDB gagal: %v", err)
	}

	// Cek koneksi (Ping)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Ping ke MongoDB gagal: %v", err)
	}

	fmt.Println("🎉 Berhasil terhubung ke MongoDB!")

	database := os.Getenv("DATABASE_NAME")
	if database == "" {
		database = "student_report" // Default database name
		log.Println("Peringatan: DATABASE_NAME tidak disetel. Menggunakan default:", database)
	}
	return client.Database(database)
}
