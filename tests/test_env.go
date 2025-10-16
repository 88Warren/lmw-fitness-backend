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

	// Set frontend URL for testing
	if os.Getenv("FRONTEND_URL") == "" {
		os.Setenv("FRONTEND_URL", "http://localhost:3000")
	}

	// Set Brevo API key for testing
	if os.Getenv("BREVO_API_KEY") == "" {
		os.Setenv("BREVO_API_KEY", "test_brevo_api_key")
	}

	// Set Brevo newsletter list ID for testing
	if os.Getenv("BREVO_NEWSLETTER_LIST_ID") == "" {
		os.Setenv("BREVO_NEWSLETTER_LIST_ID", "1")
	}

	// Set additional environment variables for testing
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "test")
	}

	if os.Getenv("ALLOWED_ORIGIN") == "" {
		os.Setenv("ALLOWED_ORIGIN", "http://localhost:3000")
	}

	if os.Getenv("ALLOW_CREDENTIALS") == "" {
		os.Setenv("ALLOW_CREDENTIALS", "true")
	}

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8083")
	}

	if os.Getenv("ADMIN_EMAIL") == "" {
		os.Setenv("ADMIN_EMAIL", "test_admin@example.com")
	}

	if os.Getenv("ADMIN_PASSWORD") == "" {
		os.Setenv("ADMIN_PASSWORD", "TestAdminPassword123!")
	}

	if os.Getenv("SEED_USER_EMAIL") == "" {
		os.Setenv("SEED_USER_EMAIL", "test_user@example.com")
	}

	if os.Getenv("SEED_USER_PASSWORD") == "" {
		os.Setenv("SEED_USER_PASSWORD", "TestUserPassword123!")
	}

	if os.Getenv("SMTP_FROM") == "" {
		os.Setenv("SMTP_FROM", "test@example.com")
	}

	if os.Getenv("SMTP_TO") == "" {
		os.Setenv("SMTP_TO", "test@example.com")
	}

	if os.Getenv("RECAPTCHA_SECRET") == "" {
		os.Setenv("RECAPTCHA_SECRET", "test_recaptcha_secret")
	}

	// Set Brevo additional environment variables
	if os.Getenv("BREVO_API_URL") == "" {
		os.Setenv("BREVO_API_URL", "https://api.brevo.com/v3")
	}

	if os.Getenv("BREVO_NEWSLETTER_DOI_TEMPLATE_ID") == "" {
		os.Setenv("BREVO_NEWSLETTER_DOI_TEMPLATE_ID", "1")
	}

	if os.Getenv("BREVO_DOI_REDIRECT_URL") == "" {
		os.Setenv("BREVO_DOI_REDIRECT_URL", "http://localhost:3000/newsletter/confirmed")
	}

	if os.Getenv("BREVO_BEGINNER_LIST_ID") == "" {
		os.Setenv("BREVO_BEGINNER_LIST_ID", "2")
	}

	if os.Getenv("BREVO_ADVANCED_LIST_ID") == "" {
		os.Setenv("BREVO_ADVANCED_LIST_ID", "3")
	}

	if os.Getenv("BREVO_MINDSET_LIST_ID") == "" {
		os.Setenv("BREVO_MINDSET_LIST_ID", "4")
	}

	if os.Getenv("BREVO_TAILORED_COACHING_LIST_ID") == "" {
		os.Setenv("BREVO_TAILORED_COACHING_LIST_ID", "5")
	}

	if os.Getenv("BREVO_BEGINNER_PROGRAM_TEMPLATE_ID") == "" {
		os.Setenv("BREVO_BEGINNER_PROGRAM_TEMPLATE_ID", "6")
	}

	if os.Getenv("BREVO_ADVANCED_PROGRAM_TEMPLATE_ID") == "" {
		os.Setenv("BREVO_ADVANCED_PROGRAM_TEMPLATE_ID", "7")
	}

	if os.Getenv("BREVO_MINDSET_PACKAGE_TEMPLATE_ID") == "" {
		os.Setenv("BREVO_MINDSET_PACKAGE_TEMPLATE_ID", "8")
	}

	if os.Getenv("BREVO_TAILORED_COACHING_TEMPLATE_ID") == "" {
		os.Setenv("BREVO_TAILORED_COACHING_TEMPLATE_ID", "9")
	}

	if os.Getenv("BREVO_ORDER_CONFIRMATION_TEMPLATE_ID") == "" {
		os.Setenv("BREVO_ORDER_CONFIRMATION_TEMPLATE_ID", "10")
	}
}
