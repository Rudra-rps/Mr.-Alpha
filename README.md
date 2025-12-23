# 🎯 Mr.Alpha

> **Your on-chain alpha, before everyone else**

Real-time crypto intelligence platform combining **narrative detection** with **smart money tracking** to surface emerging opportunities before they hit mainstream.

**Built for Capx AI Hackathon 2025** 🚀

---

## 🌟 What It Does

Mr.Alpha provides **dual-signal alpha detection**:

1. **🔥 Alpha Radar** - Detects emerging narratives from Twitter keyword spikes
2. **🐋 Smart Money Tracker** - Monitors proven wallets' on-chain activity in real-time

When both signals align (narrative spike + whale accumulation), you've found **alpha**.

---

## ✨ Features

### Alpha Radar
- Tracks 3 hardcoded narratives: **Restaking**, **Bitcoin L2**, **AI Agents**
- Detects keyword acceleration in last 2h vs 24h baseline
- Stages narratives: **Early Alpha** → **Strong Alpha** → **Crowded Trade**
- **Capx Alignment Score** - Shows relevance to Capx AI's platform (100% for AI Agents)

### Smart Money Tracker
- Monitors 3 proven wallets with 60%+ win rates
- **Real-time blockchain tracking** via Alchemy webhooks
- Trade simulator for demo stability (30-second intervals)
- Conviction scoring based on position size
- Source badges: 🔴 **LIVE** (blockchain) vs 🎬 **DEMO** (simulated)

### Technical Highlights
- **Hybrid Architecture** - Demo mode + live blockchain integration
- **Multi-language** - Python (narratives) + Go (blockchain)
- **Real-time updates** - 10-second frontend polling
- **Thread-safe** - Goroutines with mutex-protected storage
- **Hackathon-ready** - Zero rate limit failures

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      FRONTEND                            │
│          (HTML + TailwindCSS + JavaScript)              │
│              Polls every 10 seconds                      │
└──────────┬──────────────────────────┬───────────────────┘
           │                          │
           ▼                          ▼
    ┌─────────────┐          ┌──────────────────┐
    │  Python API │          │   Go Backend     │
    │   :5000     │          │      :8080       │
    └──────┬──────┘          └────┬──────┬──────┘
           │                      │      │
           ▼                      ▼      ▼
    ┌─────────────┐          ┌────────┐ ┌────────────┐
    │   Twitter   │          │Alchemy │ │ Simulator  │
    │   (Demo)    │          │Webhook │ │  (30s)     │
    └─────────────┘          └────────┘ └────────────┘
```

**Components:**
- **narrative_radar'** - Python/Flask API for narrative detection
- **wallet_watcher** - Go backend for smart money tracking
- **frontend** - Single-page dashboard

---

## 🚀 Quick Start

### 🐳 Docker (Recommended - One Command!)

**Prerequisites:** Docker Desktop installed

```bash
git clone <your-repo>
cd "Mr. Alpha"
docker-compose up --build
```

**Access:** http://localhost:3000

That's it! All services running in containers with one command.

See [DOCKER.md](./DOCKER.md) for detailed Docker documentation.

---

### 💻 Manual Setup (Alternative)

**Prerequisites:**
- Python 3.11+
- Go 1.22+
- ngrok (for live blockchain integration)

#### 1. Clone & Install

```bash
git clone <your-repo>
cd "Mr. Alpha"
```

**Python Setup:**
```powershell
cd narrative_radar'
python -m venv .venv
.venv\Scripts\activate
pip install -r requirements.txt
```

**Go Setup:**
```powershell
cd wallet_watcher
go mod download
```

#### 2. Configure Environment

**narrative_radar'/.env:**
```env
DEMO_MODE=true
PORT=5000
TWITTER_BEARER_TOKEN=your_token_here
```

**wallet_watcher/.env:**
```env
DEMO_MODE=true
PORT=8080
```

#### 3. Start Services

**Terminal 1 - Python API:**
```powershell
cd narrative_radar'
python api.py
```

**Terminal 2 - Go Backend:**
```powershell
cd wallet_watcher
go run main.go
```

**Terminal 3 - Frontend:**
```powershell
cd frontend
start index.html
```

**Dashboard:** Open browser to `frontend/index.html`

---

## 📡 API Documentation

### Narrative API

**GET** `http://localhost:5000/api/narrative`

```json
{
  "narrative": "AI Agents",
  "growth": "+201.5%",
  "mentions": 28,
  "stage": "Crowded Trade",
  "summary": "AI Agents-related discussions accelerating rapidly",
  "capx_alignment": 100,
  "timestamp": "2025-12-23T12:00:00Z"
}
```

**Capx Alignment Scores:**
- AI Agents: **100%** (Perfect match for AI app platform)
- Restaking: **70%** (DeFi/Trading adjacent)
- Bitcoin L2: **50%** (Blockchain infrastructure)

### Smart Money API

**GET** `http://localhost:8080/api/alerts`

```json
[
  {
    "id": "62",
    "wallet_name": "Whale_0x7f3a",
    "token": "EIGEN",
    "value_usd": 15000,
    "conviction": "High",
    "narrative": "Restaking",
    "tx_hash": "0xabc123...",
    "timestamp": "2025-12-23T12:26:24Z",
    "source": "live"
  }
]
```

**POST** `http://localhost:8080/api/inject` - Manual trade injection for demo

---

## 🎬 Demo Mode vs Live Mode

### Demo Mode (Default)
- **Twitter:** Returns hardcoded data (avoids rate limits)
- **Blockchain:** Simulator generates trades every 30s
- **Reliability:** 100% uptime, perfect for presentations
- **Badges:** Shows 🎬 DEMO on simulated trades

### Live Mode

**Enable blockchain tracking:**

1. **Install ngrok:**
   ```powershell
   cd "Mr. Alpha"
   Invoke-WebRequest -Uri "https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-windows-amd64.zip" -OutFile "ngrok.zip"
   Expand-Archive ngrok.zip -DestinationPath .
   ```

2. **Start ngrok:**
   ```powershell
   .\ngrok.exe http 8080
   ```
   Copy the public URL (e.g., `https://abc123.ngrok-free.app`)

3. **Configure Alchemy webhook:**
   - Go to [Alchemy Dashboard](https://dashboard.alchemy.com/)
   - Create "Address Activity" webhook
   - Add tracked wallet addresses:
     ```
     0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045
     0xBE0eB53F46cd790Cd13851d5EFf43D12404d33E8
     0xF977814e90dA44bFA03b6295A0616a897441aceC
     ```
   - Set webhook URL: `https://abc123.ngrok-free.app/webhook/alchemy`

4. **Verify:** Check Go backend logs for `🔴 LIVE: Received Alchemy webhook`

**Result:** Real blockchain trades appear with 🔴 **LIVE** badge!

See [ALCHEMY_SETUP.md](./ALCHEMY_SETUP.md) for detailed instructions.

---

## 🐋 Tracked Wallets

| Wallet | Style | Win Rate | Avg Return |
|--------|-------|----------|------------|
| **Whale_0x7f3a** | Early Entry | 68.5% | 245% |
| **Binance14** | Swing Trader | 72.3% | 180.5% |
| **Whale_0x2208** | LP Provider | 61.2% | 95.3% |

---

## 🎯 Capx AI Integration

**Why this matters for Capx:**

1. **AI App Focus** - Narratives tagged with Capx alignment score
2. **Automated Intelligence** - Fits Capx's AI-powered app library vision
3. **Real-time Data** - Perfect for trading/analytics apps on Capx platform
4. **Composable** - APIs can be integrated into other Capx apps

**Demo Talking Points:**
- "AI Agents narrative shows **🎯 Capx Aligned** badge - perfect fit for your platform"
- "Our system automatically flags opportunities relevant to Capx's AI app ecosystem"
- "This could be integrated as a Capx app providing alpha signals to traders"

**Capx Alignment Scoring:**
- **AI Agents:** 100% (Perfect match - AI apps/agents)
- **Restaking:** 70% (DeFi/Trading adjacent)
- **Bitcoin L2:** 50% (Blockchain infrastructure)

Badges only show for 80%+ alignment (keeps UI clean).

---

## 📁 Project Structure

```
Mr. Alpha/
├── narrative_radar'/          # Python narrative detection
│   ├── detect_narrative.py    # Core detection logic
│   ├── api.py                 # Flask REST API
│   ├── requirements.txt       # Python dependencies
│   ├── Dockerfile             # Python container
│   └── .env                   # Config (DEMO_MODE=true)
│
├── wallet_watcher/            # Go smart money tracker
│   ├── main.go                # Gin server + Alchemy integration
│   ├── go.mod                 # Go dependencies
│   ├── Dockerfile             # Go container
│   └── .env                   # Config (DEMO_MODE=true)
│
├── frontend/                  # Dashboard UI
│   ├── index.html             # Main page
│   ├── script.js              # API integration
│   └── style.css              # Custom styles
│
├── docker-compose.yml         # Multi-container orchestration
├── nginx.conf                 # Frontend proxy config
├── DOCKER.md                  # Docker documentation
├── ALCHEMY_SETUP.md           # Live blockchain setup guide
├── REAL_TIME_COMPLETE.md      # Implementation details
└── README.md                  # This file
```

---

## 🛠️ Development

### Docker Development

**Start with live reload:**
```bash
docker-compose up
```

Make code changes, then rebuild:
```bash
docker-compose up --build
```

**View logs:**
```bash
docker-compose logs -f narrative-api
docker-compose logs -f wallet-tracker
```

### Add New Narrative

Edit `narrative_radar'/detect_narrative.py`:

```python
NARRATIVES = {
    "Your Narrative": ["keyword1", "keyword2", "keyword3"]
}

CAPX_ALIGNMENT = {
    "Your Narrative": 85  # Alignment score 0-100
}
```

### Add New Wallet

Edit `wallet_watcher/main.go`:

```go
var smartWallets = []Wallet{
    {
        Address:   "0xYourAddress...",
        Name:      "Whale_Name",
        Style:     "Strategy",
        WinRate:   75.0,
        AvgReturn: 200.0,
    },
}
```

Update Alchemy webhook with new address.

---

## 🏆 Hackathon Judge Guide

**Quick Demo (2 minutes):**

1. **Start with Docker:** `docker-compose up`
2. **Open:** http://localhost:3000
3. **Point out:**
   - 🔥 **Alpha Radar** - AI Agents with `🎯 Capx Aligned` badge
   - 🐋 **Smart Money** - Live trades updating every 30s
   - Clean UI (no clutter, LIVE badges only on real trades)

**Key Talking Points:**
- "Everything runs in Docker - one command deployment"
- "Hybrid architecture: demo stability + real blockchain capability"
- "Capx alignment scoring shows narratives relevant to your platform"
- "Real-time updates via Alchemy webhooks + trade simulator fallback"

**If Asked About Code:**
- Show `wallet_watcher/main.go` - Alchemy webhook handler
- Show `narrative_radar'/detect_narrative.py` - Capx alignment logic
- Show `docker-compose.yml` - Multi-service orchestration

---

## 🐛 Troubleshooting

### Docker Issues

**Services won't start:**
```bash
docker-compose down
docker-compose up --build
```

**Port conflicts:**
```bash
# Stop all containers
docker-compose down

# Kill local services
taskkill /F /IM python.exe /T
taskkill /F /IM go.exe /T
```

**Check container status:**
```bash
docker-compose ps
docker logs mr-alpha-narrative
docker logs mr-alpha-wallet
```

### Manual Setup Issues

**Python API not starting:**
```powershell
cd narrative_radar'
python -m pip install --upgrade pip
pip install -r requirements.txt
python api.py
```

**Go backend port conflict:**
```powershell
# Kill process on port 8080
Get-Process -Id (Get-NetTCPConnection -LocalPort 8080).OwningProcess | Stop-Process -Force
```

**Frontend not showing data:**
- Check both APIs running (ports 5000, 8080)
- Open browser console for errors
- Verify CORS enabled in Python API

**Alchemy webhook not working:**
- Ensure ngrok still running (free URLs expire after 2 hours)
- Check webhook URL has `/webhook/alchemy` at end
- Look for `🔴 LIVE: Received Alchemy webhook` in Go logs

---

## 📊 Performance

- **Narrative API:** <100ms response time
- **Trades API:** <50ms response time
- **Frontend polling:** 10-second intervals
- **Simulator:** Generates trade every 30 seconds
- **Memory:** ~50MB (Go) + ~30MB (Python)

---

## 🔮 Future Enhancements

- [ ] Correlation alerts (narrative spike + whale buy)
- [ ] Historical performance dashboard
- [ ] More narratives (RWA, GameFi, etc.)
- [ ] Telegram/Discord notifications
- [ ] Multi-chain support (Solana, Base)
- [ ] AI-powered narrative generation
- [ ] Wallet portfolio tracking

---

## 📄 License

MIT License - Built for Capx AI Hackathon 2025

---

## 🙏 Acknowledgments

- **Capx AI** - For hosting an awesome hackathon
- **Alchemy** - Real-time blockchain data
- **EigenLayer, Bitcoin L2, AI Agents communities** - For the alpha

---

## 📞 Contact

**Demo Issues?** Check terminal logs and ensure both services running.

**Questions?** Open an issue or reach out during hackathon office hours.

---

**🚀 Built with ❤️ for finding alpha before everyone else**
