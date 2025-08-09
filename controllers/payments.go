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
	"strings"

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
	const MaxBodyBytes = int64(65536)
	payload, err := io.ReadAll(io.LimitReader(ctx.Request.Body, MaxBodyBytes))
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Failed to read request body"})
		return
	}

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

		if checkoutSession.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
			log.Printf("Checkout session %s not paid. Skipping Brevo actions.", checkoutSession.ID)
			ctx.JSON(http.StatusOK, gin.H{"received": true, "message": "Payment not successful"})
			return
		}

		lineItemParams := &stripe.CheckoutSessionListLineItemsParams{
			Session: stripe.String(checkoutSession.ID),
		}
		lineItemParams.Expand = []*string{stripe.String("data.price.product")}
		lineItemIterator := pc.StripeClient.CheckoutSessions.ListLineItems(lineItemParams)

		purchasedPriceIDs := []string{}
		purchasedProductNames := []string{}
		for lineItemIterator.Next() {
			li := lineItemIterator.LineItem()
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

		listIDsToAdd := []int64{}
		emailTemplatesToSend := make(map[int64]bool)

		if pc.BrevoNewsletterListID != 0 {
			listIDsToAdd = append(listIDsToAdd, pc.BrevoNewsletterListID)
		} else {
			log.Println("Warning: Brevo Newsletter List ID not configured. Skipping newsletter subscription.")
		}

		for i, priceID := range purchasedPriceIDs {
			productName := purchasedProductNames[i]

			if priceID == pc.UltimateMindsetPackagePriceID || strings.Contains(strings.ToLower(productName), "mindset") {
				if pc.BrevoMindsetListID != 0 {
					listIDsToAdd = append(listIDsToAdd, pc.BrevoMindsetListID)
				} else {
					log.Println("Warning: Brevo Mindset Package List ID not configured.")
				}
				if pc.BrevoMindsetPackageTemplateID != 0 {
					emailTemplatesToSend[pc.BrevoMindsetPackageTemplateID] = true
				} else {
					log.Println("Warning: Brevo Mindset Package Template ID not configured.")
				}
			}

			switch priceID {
			case pc.BeginnerProgramPriceID:
				if pc.BrevoBeginnerListID != 0 {
					listIDsToAdd = append(listIDsToAdd, pc.BrevoBeginnerListID)
				}
				if pc.BrevoBeginnerProgramTemplateID != 0 {
					emailTemplatesToSend[pc.BrevoBeginnerProgramTemplateID] = true
				}
			case pc.AdvancedProgramPriceID:
				if pc.BrevoAdvancedListID != 0 {
					listIDsToAdd = append(listIDsToAdd, pc.BrevoAdvancedListID)
				}
				if pc.BrevoAdvancedProgramTemplateID != 0 {
					emailTemplatesToSend[pc.BrevoAdvancedProgramTemplateID] = true
				}
			case pc.TailoredCoachingPriceID:
				if pc.BrevoTailoredCoachingListID != 0 {
					listIDsToAdd = append(listIDsToAdd, pc.BrevoTailoredCoachingListID)
				}
				if pc.BrevoTailoredCoachingTemplateID != 0 {
					emailTemplatesToSend[pc.BrevoTailoredCoachingTemplateID] = true
				}
			}
		}

		if len(listIDsToAdd) > 0 {
			err = pc.AddContactToBrevo(customerEmail, listIDsToAdd)
			if err != nil {
				log.Printf("Error adding contact %s to Brevo lists %v: %v", customerEmail, listIDsToAdd, err)
			} else {
				currentListsViaAPI, getErr := pc.GetContactBrevo(customerEmail)
				if getErr != nil {
					log.Printf("DIAGNOSTIC ERROR: Failed to get contact %s details from Brevo after update attempt: %v", customerEmail, getErr)
				} else {
					log.Printf("DIAGNOSTIC SUCCESS: Brevo API reports contact %s is now in lists: %v (After successful PUT/POST operation)", customerEmail, currentListsViaAPI)
				}
			}
		} else {
			log.Println("No Brevo lists identified for this purchase.")
		}

		for templateID := range emailTemplatesToSend {
			emailParams := map[string]interface{}{
				"CUSTOMER_EMAIL":  customerEmail,
				"PURCHASED_ITEMS": purchasedProductNames,
			}
			err = pc.SendBrevoTransactionalEmail(customerEmail, templateID, emailParams)
			if err != nil {
				log.Printf("Error sending transactional email (Template ID: %d) to %s: %v", templateID, customerEmail, err)
			}
		}

		if pc.BrevoOrderConfirmationTemplateID != 0 {
			orderConfirmationParams := map[string]interface{}{
				"ORDER_ID":       checkoutSession.ID,
				"CUSTOMER_EMAIL": customerEmail,
				"TOTAL_AMOUNT":   fmt.Sprintf("%.2f", float64(checkoutSession.AmountTotal)/100.0),
				"CURRENCY":       checkoutSession.Currency,
				"PRODUCT_NAMES":  purchasedProductNames,
				"PURCHASE_DATE":  stripe.String(fmt.Sprintf("%d", checkoutSession.Created)),
			}
			err = pc.SendBrevoTransactionalEmail(customerEmail, pc.BrevoOrderConfirmationTemplateID, orderConfirmationParams)
			if err != nil {
				log.Printf("Error sending general order confirmation email to %s: %v", customerEmail, err)
			}
		} else {
			log.Println("Warning: Brevo Order Confirmation Template ID not configured. Skipping general order confirmation email.")
		}

	case "invoice.payment_succeeded":
		var invoice stripe.Invoice
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			log.Printf("Error parsing invoice JSON: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
			return
		}

		customerEmail := invoice.CustomerEmail
		if customerEmail == "" && invoice.Customer != nil {
			customer, custErr := pc.StripeClient.Customers.Get(invoice.Customer.ID, nil)
			if custErr == nil && customer != nil {
				customerEmail = customer.Email
			}
		}

		if customerEmail != "" {
			log.Printf("Invoice payment succeeded for invoice %s. Customer email: %s", invoice.ID, customerEmail)
		} else {
			log.Printf("Invoice payment succeeded for invoice %s, but customer email could not be determined.", invoice.ID)
		}

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
