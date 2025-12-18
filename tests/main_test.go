package service_test

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("== Starting Achievement Service Tests ==")

	// Setup global test environment
	os.Setenv("APP_ENV", "testing")

	code := m.Run()

	fmt.Println("== All Achievement Tests Done ==")
	os.Exit(code)
}
