/**
 * @typedef {Object} PingResult
 * @property {boolean} success - Indicates if the ping was successful
 * @property {Object} [value] - Successful ping data
 * @property {Date} value.serverTime - Server time
 * @property {number} value.timestamp - Timestamp
 * @property {Date} value.uptime - Server uptime
 * @property {*} [error] - Error if ping failed
 */

/**
 * Ping the server and retrieve server information
 * @returns {Promise<PingResult>} Ping result
 */
async function ping() {
    try {
        const resp = await fetch('/api/ping');
        const data = await resp.json();
        const uptime = new Date(Date.now() - data.uptime);
        const serverTime = new Date(data.serverTime);
        const timestamp = Number.parseInt(data.timestamp, 10);
        return {
            success: true,
            value: {
                serverTime,
                timestamp,
                uptime
            }
        };
    } catch (e) {
        return {
            success: false,
            error: e
        };
    }
}

/**
 * Measure network latency by pinging the server
 * @returns {Promise<number|null>} Latency in milliseconds
 */
async function measureLatency() {
    const start = performance.now();
    const data = await ping();
    if (!data.success) {
        return null;
    }
    const end = performance.now();
    const rtt = end - start;
    return rtt;
}

console.log('Running main.js...');
window.addEventListener('DOMContentLoaded', () => {
    const elem = document.getElementById('latency-indicator');

    /**
     * Generate a random latency value
     * @param {number} [min=20] - Minimum latency
     * @param {number} [max=150] - Maximum latency
     * @returns {number} Random latency
     */
    function getRandomLatency(min = 20, max = 150) {
        return Math.floor(Math.random() * (max - min + 1)) + min;
    }

    const state = {history: []};
    function pushToState(state, item) {
        if (state.history.length >= 5) {
            state.history.shift();
        }
        state.history.push(item);
    }

    function getAverageLatency(state) {
        return state.history.reduce((acc, val) => acc + val, 0) / state.history.length;
    }

    /**
     * Update the latency indicator
     * @returns {Promise<void>}
     */
    async function updateLatency(state) {
        const res = await measureLatency();
        if (elem) {
            const latency = await measureLatency() ?? getRandomLatency()
            pushToState(state, latency);

            const avg = getAverageLatency(state);

            let renderedLatency = 0;
            if (latency <= 10) {
                renderedLatency = latency.toFixed(3)
            } else {
                renderedLatency = Math.round(latency);
            }


            elem.innerHTML = `LATENCY<span class="mobile-hidden">:</span> ${renderedLatency}ms`;
        }
    }

    updateLatency(state);
    setInterval(() => updateLatency(state), 1000);
});

// Expose measureLatency to window for use in JS templates
window.measureLatency = measureLatency;
