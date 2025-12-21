(function () {
    window.addEventListener("DOMContentLoaded", () => {
        const canvas = document.getElementById("starfield");
        if (!canvas) return;

        const ctx = canvas.getContext("2d");

        function resize() {
            canvas.width = window.innerWidth;
            canvas.height = window.innerHeight;
        }

        resize();
        window.addEventListener("resize", resize);

        const numStars = 400;
        const stars = [];

        // Stable time base so the animation phase is consistent across reloads
        const baseTime = performance.timeOrigin || Date.now();

        function getElapsed() {
            return (Date.now() - baseTime) / 1000; // seconds
        }

        // Deterministic "random"
        function seededRandom(seed) {
            const x = Math.sin(seed) * 10000;
            return x - Math.floor(x);
        }

        // Initialize stars with seed-based positions
        for (let i = 0; i < numStars; i++) {
            const seed = i + 1;
            stars.push({
                seed,
                x0: (seededRandom(seed * 1.23) * 2 - 1), // -1..1
                y0: (seededRandom(seed * 4.56) * 2 - 1), // -1..1
                z0: seededRandom(seed * 7.89),           // 0..1
                speed: 0.6 + seededRandom(seed * 3.21) * 1.2
            });
        }

        function draw() {
            const w = canvas.width;
            const h = canvas.height;
            const cx = w / 2;
            const cy = h / 2;

            ctx.fillStyle = "black";
            ctx.fillRect(0, 0, w, h);

            const elapsed = getElapsed();
            const depthRange = w; // scale depth to width
            const fov = 400;

            ctx.fillStyle = "white";

            for (const star of stars) {
                // Distance traveled along z, wrapped so stars respawn behind camera
                const traveled = (star.z0 * depthRange + elapsed * star.speed * 200) % depthRange;
                const z = depthRange - traveled; // coming toward camera

                const k = fov / z;
                const x = star.x0 * w * k + cx;
                const y = star.y0 * h * k + cy;

                if (x < 0 || x >= w || y < 0 || y >= h) continue;

                const size = (1 - z / depthRange) * 2.6; // closer = bigger
                ctx.beginPath();
                ctx.arc(x, y, size, 0, Math.PI * 2);
                ctx.fill();
            }

            requestAnimationFrame(draw);
        }

        draw();
    });
})();
