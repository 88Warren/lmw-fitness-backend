// controllers/payment_controller.go
package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/client"
	"github.com/stripe/stripe-go/v82/webhook"
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
	FrontendURL                   string
	StripeWebhookSecret           string

	BrevoAPIKey                    string
	BrevoNewsletterListID          int64
	BrevoBeginnerListID            int64
	BrevoAdvancedListID            int64
	BrevoMindsetPackageTemplateID  int64
	BrevoBeginnerProgramTemplateID int64
	BrevoAdvancedProgramTemplateID int64
}

func NewPaymentController() *PaymentController {
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
	brevoBeginnerListID, _ := ParseInt64Env("BREVO_BEGINNER_LIST_ID")
	brevoAdvancedListID, _ := ParseInt64Env("BREVO_ADVANCED_LIST_ID")
	brevoMindsetPackageTemplateID, _ := ParseInt64Env("BREVO_MINDSET_PACKAGE_TEMPLATE_ID")
	brevoBeginnerProgramTemplateID, _ := ParseInt64Env("BREVO_BEGINNER_PROGRAM_TEMPLATE_ID")
	brevoAdvancedProgramTemplateID, _ := ParseInt64Env("BREVO_ADVANCED_PROGRAM_TEMPLATE_ID")

	return &PaymentController{
		StripeClient:                   sc,
		UltimateMindsetPackagePriceID:  ultimateMindsetPriceID,
		TailoredCoachingPriceID:        tailoredCoachingPriceID,
		BeginnerProgramPriceID:         beginnerPriceID,
		AdvancedProgramPriceID:         advancedPriceID,
		FrontendURL:                    frontendURL,
		StripeWebhookSecret:            stripeWebhookSecret,
		BrevoAPIKey:                    brevoAPIKey,
		BrevoNewsletterListID:          brevoNewsletterListID,
		BrevoBeginnerListID:            brevoBeginnerListID,
		BrevoAdvancedListID:            brevoAdvancedListID,
		BrevoMindsetPackageTemplateID:  brevoMindsetPackageTemplateID,
		BrevoBeginnerProgramTemplateID: brevoBeginnerProgramTemplateID,
		BrevoAdvancedProgramTemplateID: brevoAdvancedProgramTemplateID,
	}
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

	mindsetPackageProcessedForDiscount := false

	for _, item := range req.Items {
		if req.IsDiscountApplied && item.PriceID == pc.UltimateMindsetPackagePriceID && !mindsetPackageProcessedForDiscount {
			originalMindsetPackagePricePence, err := pc.getPriceUnitAmount(pc.UltimateMindsetPackagePriceID)
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

	params := &stripe.CheckoutSessionParams{
		LineItems:  lineItems,
		Mode:       stripe.String(mode),
		SuccessURL: stripe.String(pc.FrontendURL + "/payment-success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(pc.FrontendURL + "/payment-cancelled"),
	}

	if req.CustomerEmail != "" {
		params.CustomerEmail = stripe.String(req.CustomerEmail)
		log.Printf("Setting customer email for checkout session: %s", req.CustomerEmail)
	} else {
		log.Println("No customer email provided for checkout session.")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Customer email is required"})
	}

	log.Printf("Creating checkout session with params: %+v", params)
	log.Printf("Creating checkout session with params: %+v, CustomerEmail value: %s", params, *params.CustomerEmail)

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
	// This is a placeholder. You'll need to use a Brevo SDK or make an HTTP POST request
	// to Brevo's "Add or Update a Contact" endpoint (https://developers.brevo.com/reference/createcontact).
	// The lists parameter should be the IDs of the lists you want to add the contact to.
	// Example using pseudo-code for HTTP request:
	brevoURL := "https://api.brevo.com/v3/contacts"
	payload := map[string]interface{}{
		"email":            email,
		"listIds":          listIDs,
		"emailBlacklisted": false,
		"smsBlacklisted":   false,
		"attributes":       map[string]string{},
	}
	jsonPayload, _ := json.Marshal(payload)

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

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add contact to Brevo, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	log.Printf("Successfully added/updated contact %s in Brevo", email)
	return nil
}

func (pc *PaymentController) SendBrevoTransactionalEmail(email string, templateID int64, params map[string]interface{}) error {
	log.Printf("Attempting to send transactional email to %s using template %d", email, templateID)
	// Placeholder. You'll use Brevo's "Send a transactional email" endpoint (https://developers.brevo.com/reference/sendtransacemail).
	// The templateId and to parameters are crucial. params can be used for dynamic content.

	brevoURL := "https://api.brevo.com/v3/smtp/email"
	payload := map[string]interface{}{
		"to":         []map[string]string{{"email": email}},
		"templateId": templateID,
		"params":     params,
	}
	jsonPayload, _ := json.Marshal(payload)

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
	const MaxBodyBytes = int64(65536)
	payload, err := io.ReadAll(io.LimitReader(ctx.Request.Body, MaxBodyBytes))
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to read request body"})
		return
	}

	// IMPORTANT: Verify the webhook signature
	// You get this secret from your Stripe Dashboard -> Developers -> Webhooks -> Select your endpoint -> Click to reveal secret
	endpointSecret := pc.StripeWebhookSecret
	event, err := webhook.ConstructEvent(payload, ctx.Request.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook signature"})
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
			return
		}

		customerEmail := checkoutSession.CustomerDetails.Email
		if customerEmail == "" {
			log.Printf("Checkout session completed but no customer email found for session %s", checkoutSession.ID)
			ctx.JSON(http.StatusOK, gin.H{"received": true, "message": "No customer email to process"})
			return
		}

		log.Printf("Checkout session %s completed for email: %s", checkoutSession.ID, customerEmail)

		// 1. Add to Newsletter (Default)
		listIDsToAdd := []int64{pc.BrevoNewsletterListID}
		if pc.BrevoNewsletterListID == 0 {
			log.Printf("Warning: Brevo Newsletter List ID not configured. Skipping newsletter subscription.")
			listIDsToAdd = []int64{}
		}

		// 2. Check purchased items for specific programs and mindset package

		lineItemIterator := pc.StripeClient.CheckoutSessions.ListLineItems(&stripe.CheckoutSessionListLineItemsParams{})
		purchasedProductIDs := []string{}
		for lineItemIterator.Next() {
			li := lineItemIterator.LineItem()
			if li.Price != nil {
				purchasedProductIDs = append(purchasedProductIDs, li.Price.ID)
			}
		}
		log.Printf("Purchased product price IDs: %v", purchasedProductIDs)

		hasMindsetPackage := false
		hasBeginnerProgram := false
		hasAdvancedProgram := false

		for _, priceID := range purchasedProductIDs {
			if priceID == pc.UltimateMindsetPackagePriceID {
				hasMindsetPackage = true
			}
			if priceID == pc.BeginnerProgramPriceID {
				hasBeginnerProgram = true
			}
			if priceID == pc.AdvancedProgramPriceID {
				hasAdvancedProgram = true
			}
		}

		if hasBeginnerProgram && pc.BrevoBeginnerListID != 0 {
			listIDsToAdd = append(listIDsToAdd, pc.BrevoBeginnerListID)
		}
		if hasAdvancedProgram && pc.BrevoAdvancedListID != 0 {
			listIDsToAdd = append(listIDsToAdd, pc.BrevoAdvancedListID)
		}

		err = pc.AddContactToBrevo(customerEmail, listIDsToAdd)
		if err != nil {
			log.Printf("Error adding contact %s to Brevo: %v", customerEmail, err)
			// Decide how to handle this: log, retry, send alert. Don't block Stripe.
		}

		// Send Mindset Package attachments
		if hasMindsetPackage && pc.BrevoMindsetPackageTemplateID != 0 {
			err = pc.SendBrevoTransactionalEmail(customerEmail, pc.BrevoMindsetPackageTemplateID, map[string]interface{}{
				"CUSTOMER_NAME": "Valued Customer", // Customize with actual customer name if you capture it
				// Add any other params your Brevo template expects, e.g., download links
				// "DOWNLOAD_LINK_WORKSHEET": "https://yourdomain.com/downloads/worksheet.pdf",
				// "DOWNLOAD_LINK_GUIDE": "https://yourdomain.com/downloads/guide.pdf",
			})
			if err != nil {
				log.Printf("Error sending mindset package email to %s: %v", customerEmail, err)
			}
		}

		// Optional: Send program-specific welcome emails if not handled by automation workflows
		if hasBeginnerProgram && pc.BrevoBeginnerProgramTemplateID != 0 {
			err = pc.SendBrevoTransactionalEmail(customerEmail, pc.BrevoBeginnerProgramTemplateID, map[string]interface{}{})
			if err != nil {
				log.Printf("Error sending beginner program email to %s: %v", customerEmail, err)
			}
		}
		if hasAdvancedProgram && pc.BrevoAdvancedProgramTemplateID != 0 {
			err = pc.SendBrevoTransactionalEmail(customerEmail, pc.BrevoAdvancedProgramTemplateID, map[string]interface{}{})
			if err != nil {
				log.Printf("Error sending advanced program email to %s: %v", customerEmail, err)
			}
		}

	case "invoice.payment_succeeded":
		// This event is important for subscriptions.
		// You might want to update user records, grant access, etc.
		// For tailored coaching, this confirms a recurring payment.
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			log.Printf("Error parsing invoice JSON: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
			return
		}
		log.Printf("Invoice payment succeeded for invoice %s. Customer email: %s", invoice.ID, invoice.CustomerEmail)

	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	ctx.JSON(http.StatusOK, gin.H{"received": true})
}

func (pc *PaymentController) getPriceUnitAmount(priceID string) (int64, error) {
	price, err := pc.StripeClient.Prices.Get(priceID, nil)
	if err != nil {
		return 0, err
	}
	return price.UnitAmount, nil
}
