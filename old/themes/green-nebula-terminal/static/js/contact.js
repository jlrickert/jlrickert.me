// Contact form
document.getElementById("contactForm").addEventListener(
    "submit",
    (e) => {
        e.preventDefault();
        const formData = new FormData(e.target);
        const name = formData.get("name");
        const email = formData.get("email");
        const message = formData.get("message");

        // Simulate form submission
        alert(
            `Message sent!\n\n> From: ${name}\n> Email: ${email}\n> Message: ${message}\n\n[This is a demo - no actual message was sent]`,
        );
        e.target.reset();
    },
);
