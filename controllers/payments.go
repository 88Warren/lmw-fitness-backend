// controllers/payment_controller.go
package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/client"
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
}

type PaymentController struct {
	StripeClient                  *client.API
	UltimateMindsetPackagePriceID string
	TailoredCoachingPriceID       string
	FrontendURL                   string
}

func NewPaymentController() *PaymentController {
	stripeSecretKey := getEnvVar("STRIPE_SECRET_KEY")
	if stripeSecretKey == "" {
		log.Fatalf("FATAL: STRIPE_SECRET_KEY environment variable not set.")
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

	frontendURL := getEnvVar("FRONTEND_URL")
	if frontendURL == "" {
		log.Fatalf("FATAL: FRONTEND_URL environment variable not set.")
	}

	return &PaymentController{
		StripeClient:                  sc,
		UltimateMindsetPackagePriceID: ultimateMindsetPriceID,
		TailoredCoachingPriceID:       tailoredCoachingPriceID,
		FrontendURL:                   frontendURL,
	}
}

func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Warning: Environment variable %s not set.", key)
	}
	return value
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

	log.Printf("Creating checkout session with params: %+v", params)

	s, err := pc.StripeClient.CheckoutSessions.New(params)
	if err != nil {
		log.Printf("Error creating checkout session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"url": s.URL})
}

func (pc *PaymentController) StripeWebhook(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (pc *PaymentController) getPriceUnitAmount(priceID string) (int64, error) {
	price, err := pc.StripeClient.Prices.Get(priceID, nil)
	if err != nil {
		return 0, err
	}
	return price.UnitAmount, nil
}
