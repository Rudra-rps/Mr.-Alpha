package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
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
	Source      string    `json:"source"` // "live" or "simulated"
}

// Alchemy Webhook Structures
type AlchemyWebhook struct {
	WebhookID string       `json:"webhookId"`
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"createdAt"`
	Type      string       `json:"type"`
	Event     AlchemyEvent `json:"event"`
}

type AlchemyEvent struct {
	Network  string            `json:"network"`
	Activity []AlchemyActivity `json:"activity"`
}

type AlchemyActivity struct {
	FromAddress string  `json:"fromAddress"`
	ToAddress   string  `json:"toAddress"`
	BlockNum    string  `json:"blockNum"`
	Hash        string  `json:"hash"`
	Value       float64 `json:"value"`
	Asset       string  `json:"asset"`
	Category    string  `json:"category"`
	RawContract struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
	} `json:"rawContract"`
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
var tradeMutex = sync.RWMutex{}
var tradeCounter = 0

// Demo configuration
var demoMode = false
var simulatorInterval = 30 * time.Second

// Token pool for realistic variety
var tokenPool = []struct {
	Symbol    string
	Address   string
	Narrative string
}{
	{"EIGEN", "0xec53bf9167f50cdeb3ae105f56099aaab9061f83", "Restaking"},
	{"OLAS", "0x0001a500a6b18995b03f44bb040a5ffc28e45cb0", "AI Agents"},
	{"RUNE", "0x3155ba85d5f96b2d030a4966af206230e46849cb", "Bitcoin L2"},
	{"LDO", "0x5a98fcbea516cf06857215779fd812ca3bef1b32", "Restaking"},
	{"FET", "0xaea46a60368a7bd060eec7df8cba43b7ef41ad85", "AI Agents"},
}

func main() {
	// Load env vars
	godotenv.Load()

	// Check demo mode
	demoModeEnv := os.Getenv("DEMO_MODE")
	if demoModeEnv == "true" || demoModeEnv == "1" {
		demoMode = true
	}

	// Setup Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(corsMiddleware())

	// Routes
	r.GET("/health", healthCheck)
	r.GET("/api/wallets", getWallets)
	r.GET("/api/alerts", getAlerts)
	r.POST("/webhook/alchemy", alchemyWebhook)
	r.POST("/api/inject", injectTrade) // Manual injection endpoint

	// Initialize with demo data
	initDemoData()

	// Start trade simulator if in demo mode
	if demoMode {
		log.Printf("🎬 DEMO MODE ENABLED - Trade simulator active\n")
		go tradeSimulator()
	} else {
		log.Printf("📡 LIVE MODE - Waiting for real transactions\n")
	}

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
		"status":    "ok",
		"service":   "smart_money_tracker",
		"demo_mode": demoMode,
		"trades":    len(recentTrades),
	})
}

func getWallets(c *gin.Context) {
	c.JSON(http.StatusOK, smartWallets)
}

func getAlerts(c *gin.Context) {
	tradeMutex.RLock()
	defer tradeMutex.RUnlock()

	// Return recent trades (last 10)
	limit := 10
	if len(recentTrades) < limit {
		limit = len(recentTrades)
	}

	c.JSON(http.StatusOK, recentTrades[:limit])
}

func alchemyWebhook(c *gin.Context) {
	// Try to parse as Alchemy webhook first
	var webhook AlchemyWebhook
	if err := c.BindJSON(&webhook); err != nil {
		// Fallback to generic payload for manual injection
		var payload map[string]interface{}
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		trade := createTradeFromPayload(payload)
		addTrade(trade)
		c.JSON(http.StatusOK, gin.H{"status": "received", "trade": trade})
		return
	}

	log.Printf("🔴 LIVE: Received Alchemy webhook")

	// Parse real blockchain transactions
	for _, activity := range webhook.Event.Activity {
		// Only process token transfers (not ETH transfers)
		if activity.Category == "token" || activity.Category == "erc20" {
			trade := parseAlchemyActivity(activity)
			if trade.WalletAddr != "" {
				addTrade(trade)
				log.Printf("🔴 LIVE TRADE: %s bought %s ($%.0f)", trade.WalletName, trade.Token, trade.ValueUSD)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "processed",
		"activity": len(webhook.Event.Activity),
	})
}

// Manual injection endpoint for testing
func injectTrade(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trade := createTradeFromPayload(payload)
	addTrade(trade)

	log.Printf("💉 Manually injected trade: %s bought %s", trade.WalletName, trade.Token)

	c.JSON(http.StatusOK, gin.H{
		"status": "injected",
		"trade":  trade,
	})
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
			Source:      "simulated",
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
			Source:      "simulated",
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
			Source:      "simulated",
		},
	}

	tradeCounter = 3

	log.Printf("✅ Initialized with %d demo trades\n", len(recentTrades))
}

// Trade Simulator - generates realistic trades every 30 seconds
func tradeSimulator() {
	ticker := time.NewTicker(simulatorInterval)
	defer ticker.Stop()

	log.Printf("🔄 Trade simulator started (every %v)\n", simulatorInterval)

	for range ticker.C {
		trade := generateSimulatedTrade()
		addTrade(trade)
		log.Printf("🤖 Simulated: %s bought $%.0f of %s (%s conviction)\n",
			trade.WalletName, trade.ValueUSD, trade.Token, trade.Conviction)
	}
}

// Generate a realistic simulated trade
func generateSimulatedTrade() Trade {
	// Random wallet
	wallet := smartWallets[rand.Intn(len(smartWallets))]

	// Random token from pool
	token := tokenPool[rand.Intn(len(tokenPool))]

	// Random value between $3k-$20k
	valueUSD := float64(3000 + rand.Intn(17000))

	// Position percentage based on value
	var positionPct float64
	var conviction string
	if valueUSD > 12000 {
		positionPct = 10.0 + float64(rand.Intn(8))
		conviction = "High"
	} else if valueUSD > 7000 {
		positionPct = 5.0 + float64(rand.Intn(5))
		conviction = "Medium"
	} else {
		positionPct = 2.0 + float64(rand.Intn(4))
		conviction = "Low"
	}

	tradeCounter++
	txHash := fmt.Sprintf("0x%x...", rand.Intn(0xffffff))

	return Trade{
		ID:          strconv.Itoa(tradeCounter),
		WalletAddr:  wallet.Address,
		WalletName:  wallet.Name,
		Token:       token.Symbol,
		TokenAddr:   token.Address,
		ValueUSD:    valueUSD,
		PositionPct: positionPct,
		Conviction:  conviction,
		Narrative:   token.Narrative,
		TxHash:      txHash,
		Timestamp:   time.Now(),
		Source:      "simulated",
	}
}

// Add trade to the front of the list (most recent first)
func addTrade(trade Trade) {
	tradeMutex.Lock()
	defer tradeMutex.Unlock()

	// Prepend to slice (newest first)
	recentTrades = append([]Trade{trade}, recentTrades...)

	// Keep only last 20 trades
	if len(recentTrades) > 20 {
		recentTrades = recentTrades[:20]
	}
}

// Create trade from webhook/injection payload
func createTradeFromPayload(payload map[string]interface{}) Trade {
	// Extract values with defaults
	walletName, _ := payload["wallet_name"].(string)
	if walletName == "" {
		walletName = "Unknown Wallet"
	}

	token, _ := payload["token"].(string)
	if token == "" {
		token = "UNKNOWN"
	}

	valueUSD, _ := payload["value_usd"].(float64)
	if valueUSD == 0 {
		valueUSD = 5000
	}

	conviction, _ := payload["conviction"].(string)
	if conviction == "" {
		conviction = "Medium"
	}

	narrative, _ := payload["narrative"].(string)
	if narrative == "" {
		narrative = "General"
	}

	tradeCounter++

	return Trade{
		ID:          strconv.Itoa(tradeCounter),
		WalletAddr:  smartWallets[0].Address,
		WalletName:  walletName,
		Token:       token,
		TokenAddr:   "0x0000000000000000000000000000000000000000",
		ValueUSD:    valueUSD,
		PositionPct: 10.0,
		Conviction:  conviction,
		Narrative:   narrative,
		TxHash:      fmt.Sprintf("0x%x...", rand.Intn(0xffffff)),
		Timestamp:   time.Now(),
		Source:      "manual",
	}
}

// Parse real Alchemy activity into Trade
func parseAlchemyActivity(activity AlchemyActivity) Trade {
	// Find which wallet this is
	walletAddr := activity.FromAddress
	walletName := getWalletName(walletAddr)

	if walletName == "" {
		// Not one of our tracked wallets
		return Trade{}
	}

	// Extract token info
	tokenAddr := activity.RawContract.Address
	tokenSymbol := getTokenSymbol(tokenAddr)

	// Estimate USD value (simplified - in production use price oracle)
	valueUSD := estimateValueUSD(activity.Value, tokenSymbol)

	// Calculate conviction based on value
	var positionPct float64
	var conviction string
	if valueUSD > 12000 {
		positionPct = 12.0
		conviction = "High"
	} else if valueUSD > 7000 {
		positionPct = 8.0
		conviction = "Medium"
	} else {
		positionPct = 4.0
		conviction = "Low"
	}

	// Match to narrative
	narrative := matchTokenToNarrative(tokenSymbol)

	tradeCounter++

	return Trade{
		ID:          strconv.Itoa(tradeCounter),
		WalletAddr:  walletAddr,
		WalletName:  walletName,
		Token:       tokenSymbol,
		TokenAddr:   tokenAddr,
		ValueUSD:    valueUSD,
		PositionPct: positionPct,
		Conviction:  conviction,
		Narrative:   narrative,
		TxHash:      activity.Hash,
		Timestamp:   time.Now(),
		Source:      "live",
	}
}

// Helper: Get wallet name from address
func getWalletName(address string) string {
	for _, wallet := range smartWallets {
		if wallet.Address == address {
			return wallet.Name
		}
	}
	return ""
}

// Helper: Get token symbol from contract address (simplified)
func getTokenSymbol(address string) string {
	// Known token addresses (expand as needed)
	knownTokens := map[string]string{
		"0xec53bf9167f50cdeb3ae105f56099aaab9061f83": "EIGEN",
		"0x0001a500a6b18995b03f44bb040a5ffc28e45cb0": "OLAS",
		"0x3155ba85d5f96b2d030a4966af206230e46849cb": "RUNE",
		"0x5a98fcbea516cf06857215779fd812ca3bef1b32": "LDO",
		"0xaea46a60368a7bd060eec7df8cba43b7ef41ad85": "FET",
	}

	if symbol, exists := knownTokens[address]; exists {
		return symbol
	}
	return "UNKNOWN"
}

// Helper: Estimate USD value (simplified - in production use Chainlink/CoinGecko)
func estimateValueUSD(amount float64, token string) float64 {
	// Mock prices for demo
	prices := map[string]float64{
		"EIGEN": 3.20,
		"OLAS":  1.85,
		"RUNE":  4.50,
		"LDO":   2.10,
		"FET":   0.95,
	}

	price := prices[token]
	if price == 0 {
		price = 1.0 // Default
	}

	return amount * price
}

// Helper: Match token to narrative
func matchTokenToNarrative(token string) string {
	narrativeMap := map[string]string{
		"EIGEN": "Restaking",
		"LDO":   "Restaking",
		"OLAS":  "AI Agents",
		"FET":   "AI Agents",
		"RUNE":  "Bitcoin L2",
	}

	if narrative, exists := narrativeMap[token]; exists {
		return narrative
	}
	return "General"
}
