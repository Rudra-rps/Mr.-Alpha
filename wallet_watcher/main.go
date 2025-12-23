package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Wallet represents a tracked smart money wallet
type Wallet struct {
	Address   string  `json:"address"`
	Name      string  `json:"name"`
	Style     string  `json:"style"`
	WinRate   float64 `json:"win_rate"`
	AvgReturn float64 `json:"avg_return"`
}

// Trade represents a wallet transaction
type Trade struct {
	ID          string    `json:"id"`
	WalletAddr  string    `json:"wallet_address"`
	WalletName  string    `json:"wallet_name"`
	Token       string    `json:"token"`
	TokenAddr   string    `json:"token_address"`
	ValueUSD    float64   `json:"value_usd"`
	PositionPct float64   `json:"position_pct"`
	Conviction  string    `json:"conviction"`
	Narrative   string    `json:"narrative"`
	TxHash      string    `json:"tx_hash"`
	Timestamp   time.Time `json:"timestamp"`
}

// Hardcoded smart wallets for MVP
var smartWallets = []Wallet{
	{
		Address:   "0x7f3a152F09324f2aee916CE069D3908603449173",
		Name:      "Whale_0x7f3a",
		Style:     "Early Entry",
		WinRate:   68.5,
		AvgReturn: 245.0,
	},
	{
		Address:   "0x28C6c06298d514Db089934071355E5743bf21d60",
		Name:      "Binance14",
		Style:     "Swing Trader",
		WinRate:   72.3,
		AvgReturn: 180.5,
	},
	{
		Address:   "0x220866B1A2219f40e72f5c628B65D54268cA3A9D",
		Name:      "Whale_0x2208",
		Style:     "LP Provider",
		WinRate:   61.2,
		AvgReturn: 95.3,
	},
}

// In-memory storage for demo (replace with DB in production)
var recentTrades = []Trade{}

func main() {
	// Load env vars
	godotenv.Load()

	// Setup Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(corsMiddleware())

	// Routes
	r.GET("/health", healthCheck)
	r.GET("/api/wallets", getWallets)
	r.GET("/api/alerts", getAlerts)
	r.POST("/webhook/alchemy", alchemyWebhook)

	// Initialize with demo data
	initDemoData()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Mr.Alpha - Smart Money Tracker\n")
	fmt.Printf("Starting on http://localhost:%s\n\n", port)
	r.Run(":" + port)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "smart_money_tracker",
	})
}

func getWallets(c *gin.Context) {
	c.JSON(http.StatusOK, smartWallets)
}

func getAlerts(c *gin.Context) {
	// Return recent trades (last 10)
	limit := 10
	if len(recentTrades) < limit {
		limit = len(recentTrades)
	}

	c.JSON(http.StatusOK, recentTrades[:limit])
}

func alchemyWebhook(c *gin.Context) {
	// Parse Alchemy webhook payload
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process transaction (simplified for MVP)
	log.Printf("Received webhook: %v", payload)

	// TODO: Parse transaction and create Trade
	// For now, just acknowledge
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func initDemoData() {
	// Create realistic demo trades
	recentTrades = []Trade{
		{
			ID:          "1",
			WalletAddr:  smartWallets[0].Address,
			WalletName:  smartWallets[0].Name,
			Token:       "EIGEN",
			TokenAddr:   "0xec53bf9167f50cdeb3ae105f56099aaab9061f83",
			ValueUSD:    12400.0,
			PositionPct: 15.2,
			Conviction:  "High",
			Narrative:   "Restaking",
			TxHash:      "0xabc123...",
			Timestamp:   time.Now().Add(-10 * time.Minute),
		},
		{
			ID:          "2",
			WalletAddr:  smartWallets[1].Address,
			WalletName:  smartWallets[1].Name,
			Token:       "OLAS",
			TokenAddr:   "0x0001a500a6b18995b03f44bb040a5ffc28e45cb0",
			ValueUSD:    8750.0,
			PositionPct: 8.5,
			Conviction:  "Medium",
			Narrative:   "AI Agents",
			TxHash:      "0xdef456...",
			Timestamp:   time.Now().Add(-25 * time.Minute),
		},
		{
			ID:          "3",
			WalletAddr:  smartWallets[2].Address,
			WalletName:  smartWallets[2].Name,
			Token:       "RUNE",
			TokenAddr:   "0x3155ba85d5f96b2d030a4966af206230e46849cb",
			ValueUSD:    5200.0,
			PositionPct: 4.8,
			Conviction:  "Low",
			Narrative:   "Bitcoin L2",
			TxHash:      "0xghi789...",
			Timestamp:   time.Now().Add(-45 * time.Minute),
		},
	}

	log.Printf("✅ Initialized with %d demo trades\n", len(recentTrades))
}
