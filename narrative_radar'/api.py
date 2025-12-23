"""
Mr.Alpha - Alpha Radar API
Serves narrative detection results
"""

from flask import Flask, jsonify
from flask_cors import CORS
from detect_narrative import detect_narrative, DEMO_MODE
import json
import os
from datetime import datetime, timezone

app = Flask(__name__)
CORS(app)  # Enable CORS for frontend

# Cache for narrative data
CACHE_FILE = 'narrative_detected.json'
last_result = None

@app.route('/health', methods=['GET'])
def health():
    """Health check endpoint"""
    return jsonify({
        "status": "ok",
        "service": "alpha_radar",
        "demo_mode": DEMO_MODE
    })

@app.route('/api/narrative', methods=['GET'])
def get_narrative():
    """Get current trending narrative"""
    global last_result
    
    try:
        # Try to load from cache first
        if os.path.exists(CACHE_FILE):
            with open(CACHE_FILE, 'r') as f:
                last_result = json.load(f)
        
        # Run fresh detection
        result = detect_narrative()
        
        if result:
            last_result = result
            # Save to cache
            with open(CACHE_FILE, 'w') as f:
                json.dump(result, f, indent=2)
            
            return jsonify(result)
        elif last_result:
            # Return cached result if available
            return jsonify(last_result)
        else:
            return jsonify({"error": "No narrative data available"}), 404
    
    except Exception as e:
        return jsonify({"error": str(e)}), 500

@app.route('/api/narratives/all', methods=['GET'])
def get_all_narratives():
    """Get all narratives with their stats"""
    try:
        from detect_narrative import get_demo_data, get_twitter_client, search_keywords, NARRATIVES, calculate_growth, determine_stage
        
        if DEMO_MODE:
            results = get_demo_data()
        else:
            client = get_twitter_client()
            results = []
            
            for narrative_name, keywords in NARRATIVES.items():
                mentions_24h = search_keywords(client, keywords, hours_ago=24)
                if mentions_24h == -1:
                    results = get_demo_data()
                    break
                    
                mentions_2h = search_keywords(client, keywords, hours_ago=2)
                if mentions_2h == -1:
                    results = get_demo_data()
                    break
                
                growth = calculate_growth(mentions_24h, mentions_2h)
                
                results.append({
                    'narrative': narrative_name,
                    'mentions_24h': mentions_24h,
                    'mentions_2h': mentions_2h,
                    'growth': growth
                })
        
        # Add stage to each
        for r in results:
            r['stage'] = determine_stage(r['growth'])
            r['growth_str'] = f"+{r['growth']}%" if r['growth'] > 0 else f"{r['growth']}%"
        
        return jsonify(results)
    
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    print("🚀 Mr.Alpha - Alpha Radar API")
    print(f"Demo Mode: {DEMO_MODE}")
    print("Starting on http://localhost:5000\n")
    app.run(debug=True, host='0.0.0.0', port=5000)
