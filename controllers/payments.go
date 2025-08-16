package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/88warren/lmw-fitness-backend/workers"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/client"
	"github.com/stripe/stripe-go/v82/webhook"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	DiscountAmount = 1000
)

type CheckoutItem struct {
	PriceID  string `json:"priceId" binding:"required"`
	Quantity int64  `json:"quantity" binding:"min=1"`
}

type CreateCheckoutSessionRequest struct {
	Items             []CheckoutItem `json:"items" binding:"required"`
	IsDiscountApplied bool           `json:"isDiscountApplied"`
	CustomerEmail     string         `json:"customerEmail"`
}

type PaymentController struct {
	StripeClient                  *client.API
	UltimateMindsetPackagePriceID string
	TailoredCoachingPriceID       string
	BeginnerProgramPriceID        string
	AdvancedProgramPriceID        string

	FrontendURL         string
	StripeWebhookSecret string
	BrevoAPIKey         string

	BrevoNewsletterListID       int64
	BrevoMindsetListID          int64
	BrevoBeginnerListID         int64
	BrevoAdvancedListID         int64
	BrevoTailoredCoachingListID int64

	BrevoMindsetPackageTemplateID    int64
	BrevoBeginnerProgramTemplateID   int64
	BrevoAdvancedProgramTemplateID   int64
	BrevoTailoredCoachingTemplateID  int64
	BrevoOrderConfirmationTemplateID int64

	DB *gorm.DB
}

func NewPaymentController(db *gorm.DB) *PaymentController {
	stripeSecretKey := getEnvVar("STRIPE_SECRET_KEY")
	if stripeSecretKey == "" {
		log.Fatalf("FATAL: STRIPE_SECRET_KEY environment variable not set.")
	}

	log.Printf("Backend Stripe secret key: %s", stripeSecretKey)

	stripeWebhookSecret := getEnvVar("STRIPE_WEBHOOK_SECRET")
	if stripeWebhookSecret == "" {
		log.Fatalf("FATAL: STRIPE_WEBHOOK_SECRET environment variable not set.")
	}

	sc := &client.API{}
	sc.Init(stripeSecretKey, nil)

	ultimateMindsetPriceID := getEnvVar("ULTIMATE_MINDSET_PACKAGE_PRICE_ID")
	if ultimateMindsetPriceID == "" {
		log.Fatalf("FATAL: ULTIMATE_MINDSET_PACKAGE_PRICE_ID environment variable not set.")
	}

	tailoredCoachingPriceID := getEnvVar("TAILORED_COACHING_PRICE_ID")
	if tailoredCoachingPriceID == "" {
		log.Fatalf("FATAL: TAILORED_COACHING_PRICE_ID environment variable not set.")
	}

	beginnerPriceID := getEnvVar("BEGINNER_PRICE_ID")
	if beginnerPriceID == "" {
		log.Fatalf("FATAL: BEGINNER_PRICE_ID environment variable not set.")
	}

	advancedPriceID := getEnvVar("ADVANCED_PRICE_ID")
	if advancedPriceID == "" {
		log.Fatalf("FATAL: ADVANCED_PRICE_ID environment variable not set.")
	}

	frontendURL := getEnvVar("FRONTEND_URL")
	if frontendURL == "" {
		log.Fatalf("FATAL: FRONTEND_URL environment variable not set.")
	}

	brevoAPIKey := getEnvVar("BREVO_API_KEY")
	if brevoAPIKey == "" {
		log.Fatalf("FATAL: BREVO_API_KEY environment variable not set.")
	}

	brevoNewsletterListID, _ := ParseInt64Env("BREVO_NEWSLETTER_LIST_ID")
	brevoMindsetListID, _ := ParseInt64Env("BREVO_MINDSET_LIST_ID")
	brevoBeginnerListID, _ := ParseInt64Env("BREVO_BEGINNER_LIST_ID")
	brevoAdvancedListID, _ := ParseInt64Env("BREVO_ADVANCED_LIST_ID")
	brevoTailoredCoachingListID, _ := ParseInt64Env("BREVO_TAILORED_COACHING_LIST_ID")

	brevoMindsetPackageTemplateID, _ := ParseInt64Env("BREVO_MINDSET_PACKAGE_TEMPLATE_ID")
	brevoBeginnerProgramTemplateID, _ := ParseInt64Env("BREVO_BEGINNER_PROGRAM_TEMPLATE_ID")
	brevoAdvancedProgramTemplateID, _ := ParseInt64Env("BREVO_ADVANCED_PROGRAM_TEMPLATE_ID")
	brevoTailoredCoachingTemplateID, _ := ParseInt64Env("BREVO_TAILORED_COACHING_TEMPLATE_ID")
	brevoOrderConfirmationTemplateID, _ := ParseInt64Env("BREVO_ORDER_CONFIRMATION_TEMPLATE_ID")

	return &PaymentController{
		StripeClient:                     sc,
		UltimateMindsetPackagePriceID:    ultimateMindsetPriceID,
		TailoredCoachingPriceID:          tailoredCoachingPriceID,
		BeginnerProgramPriceID:           beginnerPriceID,
		AdvancedProgramPriceID:           advancedPriceID,
		FrontendURL:                      frontendURL,
		StripeWebhookSecret:              stripeWebhookSecret,
		BrevoAPIKey:                      brevoAPIKey,
		BrevoNewsletterListID:            brevoNewsletterListID,
		BrevoMindsetListID:               brevoMindsetListID,
		BrevoBeginnerListID:              brevoBeginnerListID,
		BrevoAdvancedListID:              brevoAdvancedListID,
		BrevoTailoredCoachingListID:      brevoTailoredCoachingListID,
		BrevoMindsetPackageTemplateID:    brevoMindsetPackageTemplateID,
		BrevoBeginnerProgramTemplateID:   brevoBeginnerProgramTemplateID,
		BrevoAdvancedProgramTemplateID:   brevoAdvancedProgramTemplateID,
		BrevoTailoredCoachingTemplateID:  brevoTailoredCoachingTemplateID,
		BrevoOrderConfirmationTemplateID: brevoOrderConfirmationTemplateID,
		DB:                               db,
	}
}

// generateRandomPassword generates a secure, random password.
func GenerateRandomPassword() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b), nil
}

func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Warning: Environment variable %s not set.", key)
	}
	return value
}

func ParseInt64Env(key string) (int64, error) {
	s := getEnvVar(key)
	if s == "" {
		return 0, fmt.Errorf("environment variable %s not set", key)
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Printf("Warning: Could not parse environment variable %s as int64: %v", key, err)
		return 0, err
	}
	return val, nil
}

func (pc *PaymentController) FindOrCreateUser(email string) (uint, error) {
	var user models.User

	// Attempt to find the user. If not found, a new user struct will be initialized.
	tx := pc.DB.Where(models.User{Email: email}).FirstOrCreate(&user)
	if tx.Error != nil {
		return 0, fmt.Errorf("could not find or create user: %w", tx.Error)
	}

	// If the user was just created, generate a password and save it.
	if tx.RowsAffected > 0 {
		randomPassword, err := GenerateRandomPassword()
		if err != nil {
			return 0, fmt.Errorf("could not generate temporary password: %w", err)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
		if err != nil {
			return 0, fmt.Errorf("could not hash temporary password: %w", err)
		}

		user.PasswordHash = string(hashedPassword)
		user.MustChangePassword = true

		if err := pc.DB.Save(&user).Error; err != nil {
			return 0, fmt.Errorf("could not save new user password: %w", err)
		}
	}

	return user.ID, nil
}

func (pc *PaymentController) CreateAuthToken(userID uint, programName string, dayNumber int, sessionID string) (string, error) {
	token := GenerateRandomToken()

	authToken := models.AuthToken{
		UserID:      userID,
		Token:       token,
		ProgramName: programName,
		SessionID:   sessionID,
		DayNumber:   dayNumber,
		IsUsed:      false,
	}

	if err := pc.DB.Create(&authToken).Error; err != nil {
		return "", fmt.Errorf("could not save auth token: %w", err)
	}

	return token, nil
}

func GenerateRandomToken() string {
	b := make([]byte, 32) // 32 characters for a robust token
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (pc *PaymentController) CreateCheckoutSession(ctx *gin.Context) {
	var req CreateCheckoutSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Failed to bind checkout items: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if pc.StripeClient == nil {
		log.Println("Stripe client not initialised.")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	log.Printf("Frontend URL: %s", pc.FrontendURL)

	lineItems := []*stripe.CheckoutSessionLineItemParams{}

	productNames := make(map[string]string)
	productNameForPriceID := map[string]string{
		pc.UltimateMindsetPackagePriceID: "Ultimate Habit & Mindset Package",
		pc.TailoredCoachingPriceID:       "Tailored Coaching",
		pc.BeginnerProgramPriceID:        "Beginner Program",
		pc.AdvancedProgramPriceID:        "Advanced Program",
	}

	mindsetPackageProcessedForDiscount := false

	for _, item := range req.Items {
		if name, ok := productNameForPriceID[item.PriceID]; ok {
			productNames[item.PriceID] = name
		} else {
			productNames[item.PriceID] = fmt.Sprintf("Unknown Product (Price ID: %s)", item.PriceID)
		}
		if req.IsDiscountApplied && item.PriceID == pc.UltimateMindsetPackagePriceID && !mindsetPackageProcessedForDiscount {
			originalMindsetPackagePricePence, err := pc.GetPriceUnitAmount(pc.UltimateMindsetPackagePriceID)
			if err != nil {
				log.Printf("Error fetching original price for mindset package: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve product price."})
				return
			}

			discountedPricePence := originalMindsetPackagePricePence - DiscountAmount
			if discountedPricePence < 0 {
				discountedPricePence = 0
			}

			lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String("gbp"),
					UnitAmount: stripe.Int64(discountedPricePence),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String("Ultimate Habit & Mindset Package (Discounted)"),
						Description: stripe.String("Includes £10 discount when bought with another package."),
					},
				},
				Quantity: stripe.Int64(item.Quantity),
			})
			mindsetPackageProcessedForDiscount = true
			log.Printf("Applied £10 discount to Ultimate Habit & Mindset Package. New price: £%.2f", float64(discountedPricePence)/100)
		} else {
			lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(item.PriceID),
				Quantity: stripe.Int64(item.Quantity),
			})
		}
	}

	if len(lineItems) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No valid items to checkout"})
		return
	}

	mode := stripe.CheckoutSessionModePayment
	for _, item := range req.Items {
		if item.PriceID == pc.TailoredCoachingPriceID {
			mode = stripe.CheckoutSessionModeSubscription
			break
		}
	}

	metadata := make(map[string]string)
	var orderedProductNames []string
	for _, item := range req.Items {
		if name, ok := productNameForPriceID[item.PriceID]; ok {
			orderedProductNames = append(orderedProductNames, name)
		}
	}
	if len(orderedProductNames) > 0 {
		metadata["purchased_products"] = fmt.Sprintf("%v", orderedProductNames)
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  lineItems,
		Mode:       stripe.String(mode),
		SuccessURL: stripe.String(pc.FrontendURL + "/payment-success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(pc.FrontendURL + "/payment-cancelled"),
		Metadata:   metadata,
	}

	if req.CustomerEmail != "" {
		params.CustomerEmail = stripe.String(req.CustomerEmail)
		log.Printf("Setting customer email for checkout session: %s", req.CustomerEmail)
	} else {
		log.Println("No customer email provided for checkout session.")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Customer email is required"})
	}
	log.Printf("Creating checkout session with params: %+v", params)

	customerEmail := "not provided"
	if params.CustomerEmail != nil {
		customerEmail = *params.CustomerEmail
	}
	log.Printf("Creating checkout session with params: %+v, CustomerEmail value: %s", params, customerEmail)

	s, err := pc.StripeClient.CheckoutSessions.New(params)
	if err != nil {
		log.Printf("Error creating checkout session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": s.URL})
}

func (pc *PaymentController) AddContactToBrevo(email string, listIDs []int64) error {
	log.Printf("Attempting to add contact %s to Brevo lists %v", email, listIDs)

	brevoURL := "https://api.brevo.com/v3/contacts"
	payload := map[string]interface{}{
		"email":            email,
		"listIds":          listIDs,
		"emailBlacklisted": false,
		"smsBlacklisted":   false,
		"attributes":       map[string]string{},
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling Brevo payload: %w", err)
	}

	req, err := http.NewRequest("POST", brevoURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating Brevo request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", pc.BrevoAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending Brevo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent {
		log.Printf("Successfully added/updated contact %s in Brevo (via POST, status %d)", email, resp.StatusCode)
		return nil
	}

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return fmt.Errorf("failed to read Brevo response body, status: %d, original error: %w", resp.StatusCode, err)
	}

	var brevoError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	json.Unmarshal(bodyBytes, &brevoError)

	if resp.StatusCode == http.StatusBadRequest && brevoError.Code == "duplicate_parameter" && brevoError.Message == "Unable to create contact, email is already associated with another Contact" {
		log.Printf("Contact %s already exists in Brevo. Attempting to update list memberships via PUT.", email)

		updateURL := fmt.Sprintf("https://api.brevo.com/v3/contacts/%s", email)
		existingListIDs, getErr := pc.GetContactBrevo(email)
		if getErr != nil {
			log.Printf("Could not retrieve current list memberships for %s: %v", email, getErr)
		}

		finalListIDs := mergeUniqueInt64(existingListIDs, listIDs)
		updatePayload := map[string]interface{}{
			"listIds": finalListIDs,
		}

		jsonUpdatePayload, err := json.Marshal(updatePayload)
		if err != nil {
			return fmt.Errorf("error marshalling Brevo update payload for PUT: %w", err)
		}

		log.Printf("PUT Payload for %s: %s", email, string(jsonUpdatePayload))

		updateReq, err := http.NewRequest("PUT", updateURL, bytes.NewBuffer(jsonUpdatePayload))
		if err != nil {
			return fmt.Errorf("error creating Brevo PUT request: %w", err)
		}
		updateReq.Header.Set("Content-Type", "application/json")
		updateReq.Header.Set("api-key", pc.BrevoAPIKey)

		updateResp, err := client.Do(updateReq)
		if err != nil {
			return fmt.Errorf("error sending Brevo PUT request: %w", err)
		}
		defer updateResp.Body.Close()

		if updateResp.StatusCode == http.StatusNoContent {
			log.Printf("Successfully updated contact %s list memberships in Brevo (via PUT, status %d).", email, updateResp.StatusCode)
			return nil
		}

		updateBodyBytes, _ := io.ReadAll(updateResp.Body)
		return fmt.Errorf("failed to update contact lists in Brevo, status: %d, response: %s", updateResp.StatusCode, string(updateBodyBytes))
	}

	return fmt.Errorf("failed to add contact to Brevo, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
}

func (pc *PaymentController) SendBrevoTransactionalEmail(email string, templateID int64, params map[string]interface{}) error {
	log.Printf("Attempting to send transactional email to %s using template %d", email, templateID)

	brevoURL := "https://api.brevo.com/v3/smtp/email"
	payload := map[string]interface{}{
		"to":         []map[string]string{{"email": email}},
		"templateId": templateID,
		"params":     params,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling Brevo email payload: %w", err)
	}

	req, err := http.NewRequest("POST", brevoURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating Brevo email request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", pc.BrevoAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending Brevo email request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send transactional email, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("Successfully sent transactional email to %s", email)
	return nil
}

func (pc *PaymentController) StripeWebhook(ctx *gin.Context) {
	log.Printf("=== WEBHOOK RECEIVED ===")
	log.Printf("Headers: %v", ctx.Request.Header)
	log.Printf("Method: %s", ctx.Request.Method)
	log.Printf("URL: %s", ctx.Request.URL.String())

	const MaxBodyBytes = int64(65536)
	payload, err := io.ReadAll(io.LimitReader(ctx.Request.Body, MaxBodyBytes))
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to read request body"})
		return
	}

	log.Printf("Webhook payload received, length: %d bytes", len(payload))
	log.Printf("Payload preview: %s", string(payload[:min(len(payload), 200)]))

	endpointSecret := pc.StripeWebhookSecret
	log.Printf("Webhook secret configured: %v", len(endpointSecret) > 0)

	event, err := webhook.ConstructEvent(payload, ctx.Request.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook signature"})
		return
	}

	log.Printf("Webhook signature verified successfully")
	log.Printf("Received webhook event of type: %s", event.Type)
	log.Printf("Event ID: %s", event.ID)
	log.Printf("Event data preview: %s", string(event.Data.Raw[:min(len(event.Data.Raw), 500)]))

	switch event.Type {
	case "checkout.session.payment_succeeded":
		log.Println("=== Processing 'checkout.session.payment_succeeded' event ===")
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
			return
		}

		log.Printf("Checkout session details:")
		log.Printf("  - Session ID: %s", checkoutSession.ID)
		log.Printf("  - Payment Status: %s", checkoutSession.PaymentStatus)
		log.Printf("  - Customer Email: %s", checkoutSession.CustomerDetails.Email)

		// Only process if payment is actually succeeded
		if checkoutSession.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
			log.Printf("Checkout session %s payment not completed yet. Status: %s", checkoutSession.ID, checkoutSession.PaymentStatus)
			ctx.JSON(http.StatusOK, gin.H{"received": true, "message": "Payment not completed yet"})
			return
		}

		customerEmail := checkoutSession.CustomerDetails.Email
		if customerEmail == "" {
			log.Printf("Checkout session completed but no customer email found for session %s", checkoutSession.ID)
			ctx.JSON(http.StatusOK, gin.H{"received": true, "message": "No customer email to process"})
			return
		}

		log.Printf("Checkout session %s payment succeeded for email: %s", checkoutSession.ID, customerEmail)

		job := models.Job{
			SessionID:     checkoutSession.ID,
			CustomerEmail: customerEmail,
			Status:        "pending",
			Attempts:      0,
		}

		log.Printf("Creating job in database: %+v", job)
		if result := pc.DB.Create(&job); result.Error != nil {
			log.Printf("Failed to create job for session %s: %v", checkoutSession.ID, result.Error)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
			return
		}

		log.Printf("Successfully created job for session %s", checkoutSession.ID)

		// Trigger immediate processing
		if workers.GetGlobalProcessor() != nil {
			workers.GetGlobalProcessor().TriggerJobProcessing()
			log.Printf("Triggered immediate job processing for session %s", checkoutSession.ID)
		} else {
			log.Printf("Warning: Global processor not available, job will be processed on next fallback cycle")
		}

	case "checkout.session.completed":
		log.Println("=== Processing 'checkout.session.completed' event ===")
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
			return
		}

		log.Printf("Checkout session completed (fallback):")
		log.Printf("  - Session ID: %s", checkoutSession.ID)
		log.Printf("  - Payment Status: %s", checkoutSession.PaymentStatus)
		log.Printf("  - Customer Email: %s", checkoutSession.CustomerDetails.Email)

		// Only process if payment is actually completed
		if checkoutSession.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
			log.Printf("Checkout session %s payment not completed yet. Status: %s", checkoutSession.ID, checkoutSession.PaymentStatus)
			ctx.JSON(http.StatusOK, gin.H{"received": true, "message": "Payment not completed yet"})
			return
		}

		customerEmail := checkoutSession.CustomerDetails.Email
		if customerEmail == "" {
			log.Printf("Checkout session completed but no customer email found for session %s", checkoutSession.ID)
			ctx.JSON(http.StatusOK, gin.H{"received": true, "message": "No customer email to process"})
			return
		}

		log.Printf("Checkout session %s payment completed for email: %s (fallback)", checkoutSession.ID, customerEmail)

		// Check if job already exists
		var existingJob models.Job
		if pc.DB.Where("session_id = ?", checkoutSession.ID).First(&existingJob).Error == nil {
			log.Printf("Job already exists for session %s, skipping duplicate creation", checkoutSession.ID)
			ctx.JSON(http.StatusOK, gin.H{"received": true, "message": "Job already exists"})
			return
		}

		job := models.Job{
			SessionID:     checkoutSession.ID,
			CustomerEmail: customerEmail,
			Status:        "pending",
			Attempts:      0,
		}

		log.Printf("Creating job in database (fallback): %+v", job)
		if result := pc.DB.Create(&job); result.Error != nil {
			log.Printf("Failed to create job for session %s: %v", checkoutSession.ID, result.Error)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
			return
		}

		log.Printf("Successfully created job for session %s (fallback)", checkoutSession.ID)

		// Trigger immediate processing
		if workers.GetGlobalProcessor() != nil {
			workers.GetGlobalProcessor().TriggerJobProcessing()
			log.Printf("Triggered immediate job processing for session %s (fallback)", checkoutSession.ID)
		} else {
			log.Printf("Warning: Global processor not available, job will be processed on next fallback cycle")
		}

	case "payment_intent.succeeded":
		log.Println("=== Processing 'payment_intent.succeeded' event ===")
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			log.Printf("Error parsing payment intent JSON: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
			return
		}
		log.Printf("Payment intent succeeded: %s", paymentIntent.ID)

	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	ctx.JSON(http.StatusOK, gin.H{"received": true, "message": "Webhook processed successfully"})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (pc *PaymentController) ProcessPaymentSuccess(sessionID string, customerEmail string) error {
	log.Printf("[DEBUG] Starting ProcessPaymentSuccess for session %s, email %s", sessionID, customerEmail)
	log.Printf("[DEBUG] Payment controller configuration:")
	log.Printf("  - Frontend URL: %s", pc.FrontendURL)
	log.Printf("  - Brevo API Key configured: %v", pc.BrevoAPIKey != "")
	log.Printf("  - Beginner Program Template ID: %d", pc.BrevoBeginnerProgramTemplateID)
	log.Printf("  - Advanced Program Template ID: %d", pc.BrevoAdvancedProgramTemplateID)
	log.Printf("  - Beginner List ID: %d", pc.BrevoBeginnerListID)
	log.Printf("  - Advanced List ID: %d", pc.BrevoAdvancedListID)
	log.Printf("  - Beginner Program Price ID: %s", pc.BeginnerProgramPriceID)
	log.Printf("  - Advanced Program Price ID: %s", pc.AdvancedProgramPriceID)

	log.Printf("[DEBUG] Fetching checkout session %s from Stripe", sessionID)
	checkoutSession, err := pc.StripeClient.CheckoutSessions.Get(sessionID, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch checkout session: %v", err)
		return fmt.Errorf("error fetching checkout session: %w", err)
	}

	log.Printf("[DEBUG] Retrieved checkout session: ID=%s, PaymentStatus=%s, AmountTotal=%d, Currency=%s",
		checkoutSession.ID, checkoutSession.PaymentStatus, checkoutSession.AmountTotal, checkoutSession.Currency)

	if checkoutSession.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
		log.Printf("[WARN] Session %s is not paid, status: %s", sessionID, checkoutSession.PaymentStatus)
		return fmt.Errorf("session %s is not paid, status: %s", sessionID, checkoutSession.PaymentStatus)
	}

	// if checkoutSession.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
	// 	log.Printf("Checkout session %s not paid. Skipping Brevo actions.", checkoutSession.ID)
	// 	ctx.JSON(http.StatusOK, gin.H{"received": true, "message": "Payment not successful"})
	// 	return
	// }

	log.Println("Payment status is 'paid'. Proceeding with Brevo actions.")

	log.Printf("[DEBUG] Fetching line items for session %s", checkoutSession.ID)
	lineItemParams := &stripe.CheckoutSessionListLineItemsParams{
		Session: stripe.String(checkoutSession.ID),
	}
	lineItemParams.Expand = []*string{stripe.String("data.price.product")}
	lineItemIterator := pc.StripeClient.CheckoutSessions.ListLineItems(lineItemParams)

	purchasedPriceIDs := []string{}
	purchasedProductNames := []string{}
	for lineItemIterator.Next() {
		li := lineItemIterator.LineItem()
		log.Printf("[DEBUG] Processing line item: ID=%s, Description=%s, Quantity=%d",
			li.ID, li.Description, li.Quantity)
		if li.Price != nil {
			purchasedPriceIDs = append(purchasedPriceIDs, li.Price.ID)
			if li.Price.Product != nil && li.Price.Product.Name != "" {
				purchasedProductNames = append(purchasedProductNames, li.Price.Product.Name)
			} else if li.Description != "" {
				purchasedProductNames = append(purchasedProductNames, li.Description)
			}
		}
	}
	log.Printf("Purchased product price IDs: %v", purchasedPriceIDs)
	log.Printf("Purchased product names: %v", purchasedProductNames)

	isBeginnerProgramPurchased := false
	isAdvancedProgramPurchased := false
	var beginnerProgramID uint
	var advancedProgramID uint

	for _, priceID := range purchasedPriceIDs {
		if priceID == pc.BeginnerProgramPriceID {
			isBeginnerProgramPurchased = true
		}
		if priceID == pc.AdvancedProgramPriceID {
			isAdvancedProgramPurchased = true
		}
	}
	log.Printf("Beginner Program purchased: %v", isBeginnerProgramPurchased)
	log.Printf("Advanced Program purchased: %v", isAdvancedProgramPurchased)

	// Fetch beginner program
	var beginnerProgram models.WorkoutProgram
	if err := pc.DB.Where("name = ?", "beginner-program").First(&beginnerProgram).Error; err == nil {
		beginnerProgramID = beginnerProgram.ID
	} else {
		log.Printf("[WARN] Beginner program not found in DB: %v", err)
	}

	// Fetch advanced program
	var advancedProgram models.WorkoutProgram
	if err := pc.DB.Where("name = ?", "advanced-program").First(&advancedProgram).Error; err == nil {
		advancedProgramID = advancedProgram.ID
	} else {
		log.Printf("[WARN] Advanced program not found in DB: %v", err)
	}

	// Handle beginner program purchase
	if isBeginnerProgramPurchased && beginnerProgramID != 0 {
		log.Printf("Beginner Program purchased. Creating user account and auth token for %s", customerEmail)

		// Step 1: Find or create the user account
		userID, err := pc.FindOrCreateUser(customerEmail)
		if err != nil {
			log.Printf("Error finding or creating user: %v", err)
		} else {
			// Step 2: Check if user already has the beginner program
			var existingUserProgram models.UserProgram
			err = pc.DB.Where("user_id = ? AND program_id = ?", userID, beginnerProgramID).
				First(&existingUserProgram).Error

			// If the program doesn't exist for the user, create it
			if errors.Is(err, gorm.ErrRecordNotFound) {
				userProgram := models.UserProgram{
					UserID:    userID,
					ProgramID: beginnerProgramID,
				}
				if err := pc.DB.Create(&userProgram).Error; err != nil {
					log.Printf("Error saving UserProgram for user %d and program %d: %v", userID, beginnerProgramID, err)
				} else {
					log.Printf("Successfully linked user %d to program %d in UserPrograms table.", userID, beginnerProgramID)
				}
			} else if err != nil {
				log.Printf("Error checking for existing UserProgram: %v", err)
			} else {
				log.Printf("User %d already has access to program %d. Skipping creation.", userID, beginnerProgramID)
			}

			// Step 3: Create a new auth token for this session
			const programName = "beginner-program"
			const dayNumber = 1
			token, err := pc.CreateAuthToken(userID, programName, dayNumber, checkoutSession.ID)
			if err != nil {
				log.Printf("Error creating auth token: %v", err)
			} else {
				log.Printf("Generated token: %s for session %s", token, checkoutSession.ID)

				workoutURL := fmt.Sprintf("%s/workout-auth?token=%s", pc.FrontendURL, token)
				log.Printf("[DEBUG] Generated workout URL for beginner program: %s", workoutURL)

				templateParams := map[string]interface{}{
					"FIRSTNAME":    "Client",
					"WORKOUT_LINK": workoutURL,
				}

				// Only send email if template ID is configured
				if pc.BrevoBeginnerProgramTemplateID != 0 {
					if err := pc.SendBrevoTransactionalEmail(customerEmail, pc.BrevoBeginnerProgramTemplateID, templateParams); err != nil {
						log.Printf("Error sending beginner program email: %v", err)
					} else {
						log.Printf("Successfully sent Day 1 email with secure link to %s", customerEmail)
					}
				} else {
					log.Printf("Brevo Beginner Program Template ID not configured. Skipping email for %s", customerEmail)
					log.Printf("Workout URL generated: %s", workoutURL)
				}
			}
		}
	}

	// Handle advanced program purchase
	if isAdvancedProgramPurchased && advancedProgramID != 0 {
		log.Printf("Advanced Program purchased. Creating user account and auth token for %s", customerEmail)

		// Step 1: Find or create the user account
		userID, err := pc.FindOrCreateUser(customerEmail)
		if err != nil {
			log.Printf("Error finding or creating user: %v", err)
		} else {
			// Step 2: Check if user already has the advanced program
			var existingUserProgram models.UserProgram
			err = pc.DB.Where("user_id = ? AND program_id = ?", userID, advancedProgramID).
				First(&existingUserProgram).Error

			// If the program doesn't exist for the user, create it
			if errors.Is(err, gorm.ErrRecordNotFound) {
				userProgram := models.UserProgram{
					UserID:    userID,
					ProgramID: advancedProgramID,
				}
				if err := pc.DB.Create(&userProgram).Error; err != nil {
					log.Printf("Error saving UserProgram for user %d and program %d: %v", userID, advancedProgramID, err)
				} else {
					log.Printf("Successfully linked user %d to program %d in UserPrograms table.", userID, advancedProgramID)
				}
			} else if err != nil {
				log.Printf("Error checking for existing UserProgram: %v", err)
			} else {
				log.Printf("User %d already has access to program %d. Skipping creation.", userID, advancedProgramID)
			}

			// Step 3: Create a new auth token for this session
			const programName = "advanced-program"
			const dayNumber = 1
			token, err := pc.CreateAuthToken(userID, programName, dayNumber, checkoutSession.ID)
			if err != nil {
				log.Printf("Error creating auth token for Advanced Program: %v", err)
			} else {
				log.Printf("Generated token: %s for session %s", token, checkoutSession.ID)

				// Link the new token to the session ID
				var authToken models.AuthToken
				if err := pc.DB.Where("token = ?", token).First(&authToken).Error; err == nil {
					authToken.SessionID = checkoutSession.ID
					pc.DB.Save(&authToken)
					log.Printf("Successfully linked new token to session %s", checkoutSession.ID)
				} else {
					log.Printf("Error finding newly created token to link to session: %v", err)
				}

				// Build workout URL + email params
				workoutURL := fmt.Sprintf("%s/workout-auth?token=%s", pc.FrontendURL, token)
				templateParams := map[string]interface{}{
					"FIRSTNAME":    "Client",
					"WORKOUT_LINK": workoutURL,
				}

				// Only send email if template ID is configured
				if pc.BrevoAdvancedProgramTemplateID != 0 {
					if err := pc.SendBrevoTransactionalEmail(customerEmail, pc.BrevoAdvancedProgramTemplateID, templateParams); err != nil {
						log.Printf("Error sending advanced program email: %v", err)
					} else {
						log.Printf("Successfully sent Day 1 email with secure link for Advanced Program to %s", customerEmail)
					}
				} else {
					log.Printf("Brevo Advanced Program Template ID not configured. Skipping email for %s", customerEmail)
					log.Printf("Workout URL generated: %s", workoutURL)
				}
			}
		}
	}

	listIDsToAdd := []int64{}
	emailTemplatesToSend := make(map[int64]bool)

	if pc.BrevoNewsletterListID != 0 {
		listIDsToAdd = append(listIDsToAdd, pc.BrevoNewsletterListID)
		log.Printf("Added Newsletter List ID: %d", pc.BrevoNewsletterListID)
	} else {
		log.Println("Warning: Brevo Newsletter List ID not configured. Skipping newsletter subscription.")
	}

	// Iterate through purchased items and add relevant list and template IDs
	for i, priceID := range purchasedPriceIDs {
		productName := purchasedProductNames[i]
		log.Printf("Checking product: %s (Price ID: %s)", productName, priceID)

		// Logic for Mindset Package
		if priceID == pc.UltimateMindsetPackagePriceID || strings.Contains(strings.ToLower(productName), "mindset") {
			if pc.BrevoMindsetListID != 0 {
				listIDsToAdd = append(listIDsToAdd, pc.BrevoMindsetListID)
				log.Printf("Matched Mindset Package. Added List ID: %d", pc.BrevoMindsetListID)
			} else {
				log.Println("Warning: Brevo Mindset Package List ID not configured.")
			}
			if pc.BrevoMindsetPackageTemplateID != 0 {
				emailTemplatesToSend[pc.BrevoMindsetPackageTemplateID] = true
				log.Printf("Matched Mindset Package. Added Email Template ID: %d", pc.BrevoMindsetPackageTemplateID)
			} else {
				log.Println("Warning: Brevo Mindset Package Template ID not configured.")
			}
		}

		// Logic for other programs
		switch priceID {
		case pc.BeginnerProgramPriceID:
			if pc.BrevoBeginnerListID != 0 {
				listIDsToAdd = append(listIDsToAdd, pc.BrevoBeginnerListID)
				log.Printf("Matched Beginner Program. Added List ID: %d", pc.BrevoBeginnerListID)
			}
			if pc.BrevoBeginnerProgramTemplateID != 0 {
				emailTemplatesToSend[pc.BrevoBeginnerProgramTemplateID] = true
				log.Printf("Matched Beginner Program. Added Email Template ID: %d", pc.BrevoBeginnerProgramTemplateID)
			}
		case pc.AdvancedProgramPriceID:
			if pc.BrevoAdvancedListID != 0 {
				listIDsToAdd = append(listIDsToAdd, pc.BrevoAdvancedListID)
				log.Printf("Matched Advanced Program. Added List ID: %d", pc.BrevoAdvancedListID)
			}
			if pc.BrevoAdvancedProgramTemplateID != 0 {
				emailTemplatesToSend[pc.BrevoAdvancedProgramTemplateID] = true
				log.Printf("Matched Advanced Program. Added Email Template ID: %d", pc.BrevoAdvancedProgramTemplateID)
			}
		case pc.TailoredCoachingPriceID:
			if pc.BrevoTailoredCoachingListID != 0 {
				listIDsToAdd = append(listIDsToAdd, pc.BrevoTailoredCoachingListID)
				log.Printf("Matched Tailored Coaching. Added List ID: %d", pc.BrevoTailoredCoachingListID)
			}
			if pc.BrevoTailoredCoachingTemplateID != 0 {
				emailTemplatesToSend[pc.BrevoTailoredCoachingTemplateID] = true
				log.Printf("Matched Tailored Coaching. Added Email Template ID: %d", pc.BrevoTailoredCoachingTemplateID)
			}
		}
	}

	log.Printf("Final list of Brevo list IDs to add: %v", listIDsToAdd)
	log.Printf("Final list of Brevo email template IDs to send: %v", emailTemplatesToSend)

	// Add contact to Brevo lists
	if len(listIDsToAdd) > 0 {
		log.Printf("Calling AddContactToBrevo for email: %s with lists: %v", customerEmail, listIDsToAdd)
		err = pc.AddContactToBrevo(customerEmail, listIDsToAdd)
		if err != nil {
			log.Printf("Error adding contact %s to Brevo lists %v: %v", customerEmail, listIDsToAdd, err)
		} else {
			log.Printf("Successfully added contact %s to Brevo lists: %v", customerEmail, listIDsToAdd)
		}
	} else {
		log.Println("No Brevo lists identified for this purchase.")
	}

	// Send all identified transactional emails
	for templateID := range emailTemplatesToSend {
		log.Printf("Attempting to send transactional email with Template ID: %d to %s", templateID, customerEmail)
		// These parameters can be customized for each email template
		emailParams := map[string]interface{}{
			"CUSTOMER_EMAIL":  customerEmail,
			"PURCHASED_ITEMS": purchasedProductNames,
		}
		err = pc.SendBrevoTransactionalEmail(customerEmail, templateID, emailParams)
		if err != nil {
			log.Printf("Error sending transactional email (Template ID: %d) to %s: %v", templateID, customerEmail, err)
		} else {
			log.Printf("Successfully sent transactional email (Template ID: %d) to %s", templateID, customerEmail)
		}
	}

	// Send the general order confirmation email, if configured
	if pc.BrevoOrderConfirmationTemplateID != 0 {
		log.Printf("Attempting to send order confirmation email (Template ID: %d) to %s", pc.BrevoOrderConfirmationTemplateID, customerEmail)
		orderConfirmationParams := map[string]interface{}{
			"ORDER_ID":       checkoutSession.ID,
			"CUSTOMER_EMAIL": customerEmail,
			"TOTAL_AMOUNT":   fmt.Sprintf("%.2f", float64(checkoutSession.AmountTotal)/100.0),
			"CURRENCY":       checkoutSession.Currency,
			"PRODUCT_NAMES":  purchasedProductNames,
			"PURCHASE_DATE":  time.Unix(checkoutSession.Created, 0).Format("2006-01-02"),
		}
		err = pc.SendBrevoTransactionalEmail(customerEmail, pc.BrevoOrderConfirmationTemplateID, orderConfirmationParams)
		if err != nil {
			log.Printf("Error sending general order confirmation email to %s: %v", customerEmail, err)
		} else {
			log.Printf("Successfully sent general order confirmation email to %s", customerEmail)
		}
	} else {
		log.Println("Warning: Brevo Order Confirmation Template ID not configured. Skipping general order confirmation email.")
	}

	// case "invoice.payment_succeeded":
	// 	log.Println("Processing 'invoice.payment_succeeded' event.")
	// 	var invoice stripe.Invoice
	// 	err := json.Unmarshal(event.Data.Raw, &invoice)
	// 	if err != nil {
	// 		log.Printf("Error parsing invoice JSON: %v", err)
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
	// 		return
	// 	}

	// 	customerEmail := invoice.CustomerEmail
	// 	if customerEmail == "" && invoice.Customer != nil {
	// 		customer, custErr := pc.StripeClient.Customers.Get(invoice.Customer.ID, nil)
	// 		if custErr == nil && customer != nil {
	// 			customerEmail = customer.Email
	// 		}
	// 	}

	// 	if customerEmail != "" {
	// 		log.Printf("Invoice payment succeeded for invoice %s. Customer email: %s", invoice.ID, customerEmail)
	// 	} else {
	// 		log.Printf("Invoice payment succeeded for invoice %s, but customer email could not be determined.", invoice.ID)
	// 	}

	// default:
	// 	log.Printf("Unhandled event type: %s", event.Type)
	// }

	log.Printf("Background processing for session %s completed successfully", sessionID)
	return nil
}

func (pc *PaymentController) GetPriceUnitAmount(priceID string) (int64, error) {
	price, err := pc.StripeClient.Prices.Get(priceID, nil)
	if err != nil {
		return 0, err
	}
	return price.UnitAmount, nil
}

func (pc *PaymentController) GetContactBrevo(email string) ([]int64, error) {
	log.Printf("Attempting to get contact %s details from Brevo", email)

	brevoURL := fmt.Sprintf("https://api.brevo.com/v3/contacts/%s", email)

	req, err := http.NewRequest("GET", brevoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating Brevo GET request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("api-key", pc.BrevoAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending Brevo GET request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read Brevo GET response body, status: %d, error: %w", resp.StatusCode, readErr)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get contact from Brevo, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var brevoContact struct {
		ListIds []int64 `json:"listIds"`
	}
	if err := json.Unmarshal(bodyBytes, &brevoContact); err != nil {
		return nil, fmt.Errorf("error unmarshalling Brevo contact details: %w", err)
	}

	log.Printf("Successfully retrieved contact %s from Brevo via GET. Current lists reported by API: %v", email, brevoContact.ListIds)
	return brevoContact.ListIds, nil
}

func mergeUniqueInt64(existing, additions []int64) []int64 {
	m := map[int64]bool{}
	for _, id := range existing {
		m[id] = true
	}
	for _, id := range additions {
		m[id] = true
	}
	merged := []int64{}
	for id := range m {
		merged = append(merged, id)
	}
	return merged
}

func (pc *PaymentController) GetWorkoutLink(ctx *gin.Context) {
	var req struct {
		SessionID string `json:"sessionId" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request payload for GetWorkoutLink: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	log.Printf("Attempting to retrieve workout link for session: %s", req.SessionID)

	// First, try to find the token in the database
	var authToken models.AuthToken
	dbErr := pc.DB.Where("session_id = ?", req.SessionID).First(&authToken).Error

	if dbErr == nil {
		log.Printf("Auth token found: %+v", authToken)
		workoutURL := fmt.Sprintf("%s/workout-auth?token=%s", pc.FrontendURL, authToken.Token)
		log.Printf("Successfully found token for session %s. URL: %s", req.SessionID, workoutURL)
		ctx.JSON(http.StatusOK, gin.H{"workoutLink": workoutURL})
		return
	}

	// If the token was not found, check the database error
	if dbErr != gorm.ErrRecordNotFound {
		log.Printf("Database error retrieving token for session %s: %v", req.SessionID, dbErr)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workout link from database"})
		return
	}

	// If we're here, the record was not found (gorm.ErrRecordNotFound)
	log.Printf("No auth token found for session %s. Checking Stripe status.", req.SessionID)

	// Check Stripe to see if the session is valid and paid
	session, stripeErr := pc.StripeClient.CheckoutSessions.Get(req.SessionID, nil)
	if stripeErr != nil {
		log.Printf("Stripe API error for session %s: %v", req.SessionID, stripeErr)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid session ID or Stripe API error."})
		return
	}

	// If session is paid, but token not found, webhook is still processing
	if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
		log.Printf("Stripe session %s is paid, but token not found. Webhook is likely still processing.", req.SessionID)

		// Check if there's a pending job for this session
		var job models.Job
		if pc.DB.Where("session_id = ? AND status IN (?)", req.SessionID, []string{"pending", "processing"}).First(&job).Error == nil {
			ctx.JSON(http.StatusAccepted, gin.H{
				"error":   "Workout link is being prepared. Please check your email or try again in a moment.",
				"status":  "processing",
				"message": "Payment confirmed, processing your workout access...",
			})
		} else {
			ctx.JSON(http.StatusAccepted, gin.H{
				"error":   "Workout link is being prepared. Please check your email or try again in a moment.",
				"status":  "processing",
				"message": "Payment confirmed, setting up your account...",
			})
		}
		return
	}

	// If we're here, the session is not paid
	log.Printf("Stripe session %s is not paid. Payment status: %s", req.SessionID, session.PaymentStatus)
	ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not completed for this session."})
}

func (pc *PaymentController) TestWebhook(ctx *gin.Context) {
	log.Printf("=== WEBHOOK TEST ENDPOINT HIT ===")
	log.Printf("Headers: %v", ctx.Request.Header)
	log.Printf("Method: %s", ctx.Request.Method)
	log.Printf("URL: %s", ctx.Request.URL.String())

	// Test database connection
	var jobCount int64
	if err := pc.DB.Model(&models.Job{}).Count(&jobCount).Error; err != nil {
		log.Printf("Database error: %v", err)
	} else {
		log.Printf("Database connection successful. Total jobs: %d", jobCount)
	}

	// Test worker status
	workerStatus := "not available"
	if workers.GetGlobalProcessor() != nil {
		workerStatus = "available"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":                   "Webhook endpoint is reachable",
		"timestamp":                 time.Now().Unix(),
		"webhook_secret_configured": len(pc.StripeWebhookSecret) > 0,
		"database_connected":        true,
		"total_jobs":                jobCount,
		"worker_status":             workerStatus,
	})
}
