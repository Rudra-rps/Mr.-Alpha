// Configuration
// Check if running through proxy (Docker) or direct (local)
const isDocker = window.location.port === '3000';
const NARRATIVE_API = isDocker ? '/api/narrative' : 'http://localhost:5000/api/narrative';
const ALERTS_API = isDocker ? '/api/alerts' : 'http://localhost:8080/api/alerts';
const POLLING_INTERVAL = 10000; // 10 seconds

// DOM Elements
const narrativeLoading = document.getElementById('narrative-loading');
const narrativeContent = document.getElementById('narrative-content');
const narrativeName = document.getElementById('narrative-name');
const narrativeGrowth = document.getElementById('narrative-growth');
const narrativeStage = document.getElementById('narrative-stage');
const narrativeMentions = document.getElementById('narrative-mentions');
const narrativeSummary = document.getElementById('narrative-summary');
const capxBadge = document.getElementById('capx-badge');

const alertsGrid = document.getElementById('alerts-grid');
const tradeTemplate = document.getElementById('trade-card-template');

// Format Currency
function formatCurrency(value) {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
        maximumFractionDigits: 0
    }).format(value);
}

// Format Relative Time
function timeAgo(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const seconds = Math.floor((now - date) / 1000);

    let interval = seconds / 31536000;
    if (interval > 1) return Math.floor(interval) + "y ago";
    
    interval = seconds / 2592000;
    if (interval > 1) return Math.floor(interval) + "mo ago";
    
    interval = seconds / 86400;
    if (interval > 1) return Math.floor(interval) + "d ago";
    
    interval = seconds / 3600;
    if (interval > 1) return Math.floor(interval) + "h ago";
    
    interval = seconds / 60;
    if (interval > 1) return Math.floor(interval) + "m ago";
    
    if(seconds < 10) return "Just now";
    
    return Math.floor(seconds) + "s ago";
}

// Fetch Narrative Data
async function fetchNarrative() {
    try {
        const response = await fetch(NARRATIVE_API);
        if (!response.ok) throw new Error('Network response was not ok');
        const data = await response.json();
        updateNarrativeUI(data);
    } catch (error) {
        console.error('Error fetching narrative:', error);
        narrativeContent.classList.add('hidden');
        narrativeLoading.innerHTML = `<div class="text-red-400 p-4 border border-red-900/50 bg-red-900/20 rounded-lg">Connection Error: Retrying...</div>`;
        narrativeLoading.classList.remove('animate-pulse');
        narrativeLoading.classList.remove('hidden');
    }
}

// Update Narrative UI
function updateNarrativeUI(data) {
    // Hide loading, show content
    narrativeLoading.classList.add('hidden');
    narrativeContent.classList.remove('hidden');

    narrativeName.textContent = data.narrative;
    narrativeGrowth.textContent = data.growth;
    narrativeStage.textContent = `Stage: ${data.stage}`;
    narrativeMentions.textContent = data.mentions.toLocaleString();
    narrativeSummary.textContent = data.summary;

    // Subtle Capx Alignment Badge (only show for high alignment)
    capxBadge.className = 'px-3 py-1 rounded-full text-xs font-semibold';
    
    if (data.capx_alignment >= 80) {
        capxBadge.classList.remove('hidden');
        capxBadge.classList.add('bg-purple-500/20', 'text-purple-300', 'border', 'border-purple-500/30');
        capxBadge.textContent = `🎯 Capx Aligned`;
    } else if (data.capx_alignment >= 50) {
        capxBadge.classList.remove('hidden');
        capxBadge.classList.add('bg-gray-700/50', 'text-gray-400', 'border', 'border-gray-600/30');
        capxBadge.textContent = `⚡ Partial Match`;
    } else {
        capxBadge.classList.add('hidden');
    }
}

// Fetch Alerts Data
async function fetchAlerts() {
    try {
        const response = await fetch(ALERTS_API);
        if (!response.ok) throw new Error('Network response was not ok');
        const data = await response.json();
        updateAlertsUI(data);
    } catch (error) {
        console.error('Error fetching alerts:', error);
        // Show error state in alerts grid if empty
        if(alertsGrid.children.length === 0 || alertsGrid.querySelector('.animate-pulse')) {
             alertsGrid.innerHTML = `<div class="col-span-full text-center text-red-400 p-8 border border-red-900/50 bg-red-900/20 rounded-xl">Connection Error: Unable to fetch live trades</div>`;
        }
    }
}

// Update Alerts UI
function updateAlertsUI(trades) {
    alertsGrid.innerHTML = ''; // Clear current grid (including skeletons)

    trades.forEach(trade => {
        const clone = tradeTemplate.content.cloneNode(true);
        
        // Populate data
        const sourceBadge = clone.querySelector('.trade-source');
        if (trade.source === 'live') {
            sourceBadge.textContent = '🔴 LIVE';
            sourceBadge.classList.add('bg-red-500', 'text-white');
        } else {
            // Hide demo badge - keep UI clean
            sourceBadge.classList.add('hidden');
        }

        clone.querySelector('.trade-time').textContent = timeAgo(trade.timestamp);
        clone.querySelector('.trade-wallet').textContent = trade.wallet_name;
        clone.querySelector('.trade-token').textContent = trade.token;
        clone.querySelector('.trade-narrative').textContent = trade.narrative;
        clone.querySelector('.trade-value').textContent = formatCurrency(trade.value_usd);
        
        const convictionBadge = clone.querySelector('.trade-conviction');
        convictionBadge.textContent = trade.conviction;
        
        // Conviction Color Coding
        if (trade.conviction === 'High') {
            convictionBadge.classList.remove('bg-gray-800');
            convictionBadge.classList.add('bg-purple-600');
        } else if (trade.conviction === 'Medium') {
            convictionBadge.classList.remove('bg-gray-800');
            convictionBadge.classList.add('bg-blue-600');
        }

        alertsGrid.appendChild(clone);
    });
}

// Initial Load & Polling
document.addEventListener('DOMContentLoaded', () => {
    fetchNarrative();
    fetchAlerts();

    setInterval(() => {
        fetchNarrative();
        fetchAlerts();
    }, POLLING_INTERVAL);
});
