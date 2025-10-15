package tests

import "os"

// SetupTestEnvironment sets up required environment variables for testing
func SetupTestEnvironment() {
	// Set a dummy Stripe key for testing if not already set
	if os.Getenv("STRIPE_SECRET_KEY") == "" {
		os.Setenv("STRIPE_SECRET_KEY", "sk_test_dummy_key_for_testing")
	}

	// Set other required environment variables
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test_jwt_secret_key_for_testing")
	}

	if os.Getenv("SMTP_HOST") == "" {
		os.Setenv("SMTP_HOST", "localhost")
	}

	if os.Getenv("SMTP_PORT") == "" {
		os.Setenv("SMTP_PORT", "587")
	}

	if os.Getenv("SMTP_USERNAME") == "" {
		os.Setenv("SMTP_USERNAME", "test@example.com")
	}

	if os.Getenv("SMTP_PASSWORD") == "" {
		os.Setenv("SMTP_PASSWORD", "test_password")
	}
}
