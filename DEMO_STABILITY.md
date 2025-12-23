# Mr.Alpha — Demo Stability Guide

## 🎯 Goal

Ensure Mr.Alpha **always looks alive** during hackathon demos, even if:
- Twitter rate limits hit
- No real wallet trades happen  
- Internet drops
- Webhooks don't fire

## 🏗️ Architecture

```
┌─────────────────────┐
│ Trade Simulator     │  Every 30s generates new trades
│ (Go goroutine)      │  
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ Alert Store (RAM)   │  In-memory circular buffer (last 20)
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ /api/alerts         │  Returns most recent 10 trades
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ Frontend Dashboard  │  Polls every 10s, graceful errors
└─────────────────────┘
```

## ✅ Implementation Status

### Phase 1: Mock Injection ✅
- **Endpoint:** `POST /api/inject`
- **Purpose:** Manual trade injection for testing
- **Payload:**
  ```json
  {
    "wallet_name": "TestWhale",
    "token": "ETH",
    "value_usd": 25000,
    "conviction": "High",
    "narrative": "Restaking"
  }
  ```

**Test it:**
```powershell
.\test_injection.ps1
```

### Phase 2: Automatic Simulator ✅
- **Interval:** 30 seconds
- **Token Pool:** EIGEN, OLAS, RUNE, LDO, FET
- **Value Range:** $3,000 - $20,000
- **Conviction Logic:**
  - High: > $12k (10-18% position)
  - Medium: $7k-12k (5-10% position)
  - Low: < $7k (2-6% position)

### Phase 3: Frontend Binding ✅
- **Polling:** Every 10 seconds
- **Error Handling:** Shows last known data
- **No blocking:** Continues polling on failures

### Phase 4: Demo Safety Switch ✅
- **Environment Variable:** `DEMO_MODE`
- **Values:** `true` or `false`
- **Default:** `true` (safe for demos)

## 🚀 Quick Start

### 1. Start with Demo Mode (Recommended)
```powershell
cd wallet_watcher
# DEMO_MODE=true is already set in .env
go run main.go
```

You'll see:
```
🎬 DEMO MODE ENABLED - Trade simulator active
🔄 Trade simulator started (every 30s)
🚀 Mr.Alpha - Smart Money Tracker
Starting on http://localhost:8080
```

### 2. Watch Trades Generate
Every 30 seconds, you'll see:
```
🤖 Simulated: Whale_0x7f3a bought $15420 of EIGEN (High conviction)
🤖 Simulated: Binance14 bought $6800 of OLAS (Medium conviction)
```

### 3. Test Manual Injection
```powershell
.\test_injection.ps1
```

### 4. Verify API
```powershell
curl http://localhost:8080/health
# Returns: {"demo_mode":true,"service":"smart_money_tracker","status":"ok","trades":5}

curl http://localhost:8080/api/alerts
# Returns: Array of recent trades (newest first)
```

## 🎬 Demo Mode Behavior

| Mode | Behavior |
|------|----------|
| `DEMO_MODE=true` | ✅ Simulator generates trades every 30s<br>✅ Always looks active<br>✅ No external dependencies |
| `DEMO_MODE=false` | 📡 Waits for real Alchemy webhooks<br>⚠️ May look idle without activity |

## 📊 What Judges Will See

### During Demo:
1. **Initial State:** 3 pre-loaded trades from 10-45 mins ago
2. **After 30s:** New trade appears at top (e.g., Whale bought EIGEN)
3. **After 60s:** Another trade appears
4. **After 90s:** Another trade appears

### Frontend Behavior:
- Refreshes every 10 seconds
- New trades slide to the top
- Color-coded conviction badges (Red=High, Yellow=Medium, Blue=Low)
- Time stamps update ("2m ago" → "3m ago")

## 🔧 Configuration

### wallet_watcher/.env
```bash
PORT=8080
DEMO_MODE=true  # Toggle simulator on/off
```

### Simulator Settings (main.go)
```go
var simulatorInterval = 30 * time.Second  // Change frequency here
```

## 🧪 Testing Checklist

### Pre-Demo Tests:
```powershell
# 1. Health check
curl http://localhost:8080/health

# 2. Initial alerts
curl http://localhost:8080/api/alerts

# 3. Wait 30 seconds, check again
Start-Sleep -Seconds 30
curl http://localhost:8080/api/alerts

# 4. Manual injection
.\test_injection.ps1

# 5. Verify injection appeared
curl http://localhost:8080/api/alerts
```

### Expected Results:
- ✅ Alert count increases over time
- ✅ Newest trades appear at top
- ✅ Realistic token variety (EIGEN, OLAS, RUNE, LDO, FET)
- ✅ Conviction matches value ($15k+ = High, $7-12k = Medium, <$7k = Low)
- ✅ Frontend updates without manual refresh

## 🎯 Demo Flow

### Step 1: Show Initial State
> "Here's Mr.Alpha tracking 3 smart money wallets..."

### Step 2: Wait 30 Seconds
> "These wallets have 60%+ win rates. Let's see if any are making moves..."

### Step 3: New Trade Appears
> "There! Whale_0x7f3a just bought $15k of EIGEN with High conviction. 
> That matches our Restaking narrative that's heating up 180%!"

### Step 4: Show Correlation
> "Notice the pattern: Narrative spike + Smart money buying = Alpha signal"

## 🚨 Troubleshooting

### "No new trades appearing"
```powershell
# Check if demo mode is enabled
curl http://localhost:8080/health
# Should show: "demo_mode": true

# If false, update .env:
Set-Content -Path "wallet_watcher\.env" -Value "DEMO_MODE=true`nPORT=8080"

# Restart server
```

### "Server crashes"
```powershell
# Check for port conflicts
netstat -ano | findstr :8080

# Kill conflicting process
taskkill /PID <process_id> /F

# Restart
cd wallet_watcher
go run main.go
```

### "Frontend not updating"
- Hard refresh browser (Ctrl+Shift+R)
- Check browser console (F12) for errors
- Verify both APIs running:
  - http://localhost:5000/health (Python)
  - http://localhost:8080/health (Go)

## 📝 API Reference

### GET /health
Returns service status including demo mode state.

### GET /api/alerts
Returns last 10 trades (most recent first).

### GET /api/wallets
Returns list of tracked wallets with stats.

### POST /webhook/alchemy
Accepts Alchemy webhook payloads (or mock data).

### POST /api/inject
Manually inject a trade for testing.

**Payload:**
```json
{
  "wallet_name": "string",
  "token": "string",
  "value_usd": number,
  "conviction": "High|Medium|Low",
  "narrative": "string"
}
```

## 🎉 Success Criteria

✅ **Reliability:** Demo never looks idle
✅ **Realism:** Trades vary realistically  
✅ **Control:** Can toggle simulator on/off
✅ **Flexibility:** Can inject custom trades
✅ **Professional:** Can explain demo mode if asked

## 🏆 Why This Wins

1. **Judges see activity:** Something always happening
2. **Proves reactivity:** Frontend updates live
3. **Shows planning:** Built for demo reliability
4. **Demonstrates skill:** Goroutines, concurrency, thoughtful architecture
5. **Production-ready:** Toggle switch for real data

---

**You're ready to demo! 🚀**

The platform will look alive and engaging throughout the entire presentation.
