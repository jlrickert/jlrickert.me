// Starfield animation
globalThis.window.addEventListener("DOMContentLoaded", () => {
    const canvas = document.getElementById("starfield");
    const ctx = canvas.getContext("2d");
    canvas.width = globalThis.window.innerWidth;
    canvas.height = globalThis.window.innerHeight;

    const stars = [];
    const numStars = 200;

    for (let i = 0; i < numStars; i++) {
        stars.push({
            x: Math.random() * canvas.width,
            y: Math.random() * canvas.height,
            radius: Math.random() * 1.5,
            speed: Math.random() * 0.5 + 0.1,
        });
    }

    function drawStars() {
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = "#FFFFFF";

        stars.forEach((star) => {
            ctx.beginPath();
            ctx.arc(star.x, star.y, star.radius, 0, Math.PI * 2);
            ctx.fill();

            star.y += star.speed;
            if (star.y > canvas.height) {
                star.y = 0;
                star.x = Math.random() * canvas.width;
            }
        });

        requestAnimationFrame(drawStars);
    }

    drawStars();

    globalThis.window.addEventListener("resize", () => {
        canvas.width = globalThis.window.innerWidth;
        canvas.height = globalThis.window.innerHeight;
    });
});
