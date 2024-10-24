package domaincrawler_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()

	// Clean
	//
	os.Exit(code)
}
