package threat_pattern

import (
	"os"
	"testing"
	"threatreg/internal/testutil"
)

func TestMain(m *testing.M) {
	// Setup database once for all tests in this package  
	dummyT := &testing.T{}
	cleanup := testutil.SetupTestDatabase(dummyT)
	
	// Run all tests
	code := m.Run()
	
	// Cleanup database
	cleanup()
	
	// Exit with the same code as the tests
	os.Exit(code)
}