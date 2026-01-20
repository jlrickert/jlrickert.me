/**
 * @typedef {Object} PingResult
 * @property {boolean} success - Indicates if the ping was successful
 * @property {Object} [value] - Successful ping data
 * @property {Date} value.serverTime - Server time
 * @property {number} value.timestamp - Timestamp
 * @property {Date} value.uptime - Server uptime
 * @property {*} [error] - Error if ping failed
 */

class State {
    history = []

    constructor() {
        this.history = [];
    }

    pushToState(item) {
        if (this.history.length >= 5) {
            this.history.shift();
        }
        this.history.push(item);
    }

    getAverageLatency() {
        return this.history.reduce((acc, val) => acc + val, 0) / this.history.length;
    }
}

// ============================================================================
// NEW WEBSOCKET-BASED PING LOGIC
// ============================================================================

let ws = null;
let pingState = new State();

/**
 * Initialize WebSocket connection to server
 * @returns {Promise<void>}
 */
async function initWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/ws`;

    try {
        ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            console.log('WebSocket connected to', wsUrl);
            // Start sending ping requests
            sendPingRequest();
        };

        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            handlePingResponse(data);
        };

        ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        ws.onclose = () => {
            console.log('WebSocket disconnected');
            // Attempt to reconnect after 3 seconds
            setTimeout(initWebSocket, 3000);
        };
    } catch (error) {
        console.error('Failed to initialize WebSocket:', error);
        setTimeout(initWebSocket, 3000);
    }
}

/**
 * Send a ping request to the server
 */
function sendPingRequest() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        const start = performance.now();
        ws.pingStart = start;
        ws.send(JSON.stringify({ type: 'ping' }));
    }
}

/**
 * Handle ping response from server
 * @param {Object} data - Server ping response
 */
function handlePingResponse(data) {
    const end = performance.now();
    const latency = end - (ws.pingStart || end);

    pingState.pushToState(latency);

    const elem = document.getElementById('latency-indicator');
    if (elem) {
        const avg = pingState.getAverageLatency();
        let renderedLatency = avg <= 10 ? avg.toFixed(2) : Math.round(avg);
        renderedLatency = renderedLatency
            .toString()
            .padStart(4, ' ')
            .replace(/ /g, '&nbsp;');

        elem.innerHTML = `LATENCY<span class="mobile-hidden">:</span> ${renderedLatency}ms`;
    }
}

/**
 * Measure network latency via WebSocket
 * @returns {Promise<number|null>} Latency in milliseconds
 */
async function measureLatency() {
    return new Promise((resolve) => {
        if (!ws || ws.readyState !== WebSocket.OPEN) {
            resolve(null);
            return;
        }

        const start = performance.now();
        const messageHandler = () => {
            const end = performance.now();
            ws.removeEventListener('message', messageHandler);
            resolve(end - start);
        };

        ws.addEventListener('message', messageHandler);
        ws.send(JSON.stringify({ type: 'ping' }));
    });
}
// ============================================================================
// END NEW WEBSOCKET-BASED PING LOGIC
// ============================================================================

console.log('Running main.js...');
window.addEventListener('DOMContentLoaded', () => {
    // Initialize WebSocket connection for real-time ping data
    initWebSocket();

    // Send ping requests every 1 second
    setInterval(() => {
        sendPingRequest();
    }, 1000);
});

// Expose measureLatency to window for use in JS templates
window.measureLatency = measureLatency;
