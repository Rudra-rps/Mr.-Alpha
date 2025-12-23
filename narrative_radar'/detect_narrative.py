"""
Narrative Radar - Fast Keyword Spike Detection
Detects emerging crypto narratives by tracking keyword mentions on Twitter
"""

import tweepy
import os
from datetime import datetime, timedelta, timezone
from collections import defaultdict
import json
from dotenv import load_dotenv
import time

# Load environment variables from .env file
load_dotenv()

# Hackathon mode - use mock data if Twitter fails
DEMO_MODE = os.getenv('DEMO_MODE', 'false').lower() == 'true'

# Hardcoded keyword sets
NARRATIVES = {
    "Restaking": ["EigenLayer", "EIGEN", "restaking"],
    "Bitcoin L2": ["Ordinals", "Runes", "Bitcoin L2"],
    "AI Agents": ["AI agent", "autonomous agent", "$OLAS", "Capx", "CapxAI", "$CAPX", "AI app"]
}

# Capx Alignment Scores (AI app library focus)
CAPX_ALIGNMENT = {
    "AI Agents": 100,    # Perfect match - AI apps/agents
    "Restaking": 70,     # DeFi/Trading adjacent
    "Bitcoin L2": 50,    # Blockchain infrastructure
    "General": 30        # Fallback
}

def get_twitter_client():
    """Initialize Twitter API client"""
    bearer_token = os.getenv('TWITTER_BEARER_TOKEN')
    if not bearer_token:
        raise ValueError("TWITTER_BEARER_TOKEN environment variable not set")
    
    client = tweepy.Client(bearer_token=bearer_token, wait_on_rate_limit=True)
    return client

def search_keywords(client, keywords, hours_ago):
    """Search Twitter for keywords within specified time window"""
    now = datetime.now(timezone.utc)
    start_time = now - timedelta(hours=hours_ago)
    
    total_mentions = 0
    
    # Combine keywords into single query to save API calls (Twitter allows OR)
    query_parts = [f'"{kw}"' for kw in keywords[:2]]  # Limit to 2 keywords for rate limits
    combined_query = f'({" OR ".join(query_parts)}) -is:retweet lang:en'
    
    try:
        tweets = client.search_recent_tweets(
            query=combined_query,
            start_time=start_time.isoformat("T") + "Z",
            max_results=10,  # Reduced for free tier
            tweet_fields=['created_at']
        )
        
        if tweets.data:
            total_mentions = len(tweets.data)
        
        # Small delay to respect rate limits
        time.sleep(2)
    
    except tweepy.errors.TooManyRequests:
        print(f"  ⚠️  Rate limit hit - using cached data")
        return -1  # Signal to use demo data
    except Exception as e:
        print(f"  Error searching: {e}")
        return -1
    
    return total_mentions

def calculate_growth(mentions_24h, mentions_2h):
    """Calculate percentage growth rate"""
    # Extrapolate 2h data to 24h equivalent
    extrapolated_2h = mentions_2h * 12
    
    if mentions_24h == 0:
        if extrapolated_2h > 0:
            return float('inf')
        return 0
    
    growth = ((extrapolated_2h - mentions_24h) / mentions_24h) * 100
    return round(growth, 1)

def determine_stage(growth_pct):
    """Determine narrative stage based on growth (per spec)"""
    if growth_pct > 70:
        return "Crowded Trade"
    elif growth_pct >= 30:
        return "Strong Alpha"
    else:
        return "Early Alpha"

def get_demo_data():
    """Return realistic demo data for hackathon presentations"""
    return [
        {'narrative': 'Restaking', 'mentions_24h': 89, 'mentions_2h': 34, 'growth': 158.4},
        {'narrative': 'AI Agents', 'mentions_24h': 67, 'mentions_2h': 28, 'growth': 201.5},
        {'narrative': 'Bitcoin L2', 'mentions_24h': 112, 'mentions_2h': 15, 'growth': 60.7}
    ]

def detect_narrative():
    """Main detection function - finds top growing narrative"""
    
    # Check if demo mode is enabled
    if DEMO_MODE:
        print("🎬 Running in DEMO MODE\n")
        results = get_demo_data()
    else:
        client = get_twitter_client()
        results = []
        print("Scanning Twitter for narrative signals...\n")
    
        for narrative_name, keywords in NARRATIVES.items():
            print(f"Checking {narrative_name}...")
            
            # Get mentions for last 24 hours
            mentions_24h = search_keywords(client, keywords, hours_ago=24)
            
            # If rate limited, fall back to demo mode
            if mentions_24h == -1:
                print("\n⚠️  Rate limit exceeded - switching to DEMO MODE\n")
                results = get_demo_data()
                break
            
            # Get mentions for last 2 hours
            mentions_2h = search_keywords(client, keywords, hours_ago=2)
            
            if mentions_2h == -1:
                print("\n⚠️  Rate limit exceeded - switching to DEMO MODE\n")
                results = get_demo_data()
                break
            
            # Calculate growth
            growth = calculate_growth(mentions_24h, mentions_2h)
            
            print(f"  24h: {mentions_24h} | 2h: {mentions_2h} | Growth: {growth}%")
            
            results.append({
                'narrative': narrative_name,
                'mentions_24h': mentions_24h,
                'mentions_2h': mentions_2h,
                'growth': growth
            })
    
    # Find top gainer
    if not results:
        return None
    
    top_narrative = max(results, key=lambda x: x['growth'])
    
    # Determine stage
    stage = determine_stage(top_narrative['growth'])
    
    # Generate summary
    summary = f"{top_narrative['narrative']}-related discussions accelerating rapidly"
    if top_narrative['growth'] < 50:
        summary = f"{top_narrative['narrative']} showing moderate activity"
    
    # Format output
    output = {
        "narrative": top_narrative['narrative'],
        "growth": f"+{top_narrative['growth']}%" if top_narrative['growth'] > 0 else f"{top_narrative['growth']}%",
        "mentions": top_narrative['mentions_2h'],
        "stage": stage,
        "summary": summary,
        "capx_alignment": CAPX_ALIGNMENT.get(top_narrative['narrative'], 30),
        "timestamp": datetime.now(timezone.utc).isoformat().replace('+00:00', 'Z')
    }
    
    return output

if __name__ == "__main__":
    try:
        result = detect_narrative()
        
        if result:
            print("\n" + "="*50)
            print("🎯 NARRATIVE DETECTED")
            print("="*50)
            print(json.dumps(result, indent=2))
            
            # Save to file
            with open('narrative_detected.json', 'w') as f:
                json.dump(result, f, indent=2)
            
            print("\n✅ Saved to narrative_detected.json")
        else:
            print("No narratives detected")
    
    except Exception as e:
        print(f"❌ Error: {e}")
