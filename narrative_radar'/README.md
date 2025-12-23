# Narrative Radar - Fast Mode

Detects emerging crypto narratives by tracking keyword spikes on Twitter.

## Setup

1. Install dependencies:
```bash
pip install -r requirements.txt
```

2. Set up Twitter API credentials:
   - Copy `.env.example` to `.env`
   - Add your Twitter Bearer Token

3. Run the detector:
```bash
python detect_narrative.py
```

## How It Works

1. **Hardcoded Keywords**: Monitors 3 key narratives with predefined keywords
2. **Twitter Scan**: Compares mentions in last 24h vs last 2h
3. **Growth Calculation**: Identifies which narrative is heating up
4. **JSON Output**: Returns the top trending narrative

## Output Format

```json
{
  "narrative": "Restaking",
  "growth": "+180%",
  "mentions": 124,
  "stage": "HEATING_UP",
  "summary": "Restaking-related discussions accelerating rapidly",
  "timestamp": "2025-12-23T10:30:00Z"
}
```

## Stages

- **HEATING_UP**: Growth > 150%
- **EMERGING**: Growth > 50%
- **GROWING**: Growth > 0%
- **COOLING**: Declining mentions
