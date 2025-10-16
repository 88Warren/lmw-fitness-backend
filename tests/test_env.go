package tests

import "os"

// SetupTestEnvironment sets up required environment variables for testing
func SetupTestEnvironment() {
	// Set a dummy Stripe key for testing if not already set
	if os.Getenv("STRIPE_SECRET_KEY") == "" {
		os.Setenv("STRIPE_SECRET_KEY", "sk_test_dummy_key_for_testing")
	}

	// Set Stripe webhook secret for testing
	if os.Getenv("STRIPE_WEBHOOK_SECRET") == "" {
		os.Setenv("STRIPE_WEBHOOK_SECRET", "whsec_test_dummy_webhook_secret_for_testing")
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

	// Set Stripe Price IDs for testing
	if os.Getenv("ULTIMATE_MINDSET_PACKAGE_PRICE_ID") == "" {
		os.Setenv("ULTIMATE_MINDSET_PACKAGE_PRICE_ID", "price_test_ultimate_mindset_package")
	}

	if os.Getenv("TAILORED_COACHING_PRICE_ID") == "" {
		os.Setenv("TAILORED_COACHING_PRICE_ID", "price_test_tailored_coaching")
	}

	if os.Getenv("BEGINNER_PRICE_ID") == "" {
		os.Setenv("BEGINNER_PRICE_ID", "price_test_beginner_package")
	}

	if os.Getenv("ADVANCED_PRICE_ID") == "" {
		os.Setenv("ADVANCED_PRICE_ID", "price_test_advanced_package")
	}
}
