// Terminal widget
const terminalWidget = document.getElementById("terminalWidget");
const terminalHeader = document.getElementById("terminalHeader");
const terminalInput = document.getElementById("terminalInput");
const terminalOutput = document.getElementById("terminalOutput");

terminalHeader.addEventListener("click", () => {
    terminalWidget.classList.toggle("minimized");
});

const commands = {
    help:
        "Available commands: whoami, skills, projects, blog, contact, date, clear",
    whoami: "Software Developer specializing in backend systems and DevOps",
    skills:
        "Backend Development, PHP & Composer, Go Programming, DevOps & CI/CD, Database Design, Network Administration, IoT & Home Automation, Git & Version Control",
    projects:
        "4 projects found: E-Commerce Platform, Home Automation Dashboard, DevOps Automation Suite, API Gateway Service",
    blog: "4 blog posts available. Check the blog section above.",
    contact:
        "Email: dev@example.com | GitHub: github.com/username | LinkedIn: linkedin.com/in/username",
    date: new Date().toString(),
    clear: "CLEAR",
};

terminalInput.addEventListener("keypress", (e) => {
    if (e.key === "Enter") {
        const command = terminalInput.value.trim().toLowerCase();
        const outputLine = document.createElement("div");
        outputLine.className = "output-line";
        outputLine.innerHTML =
            `<span style="color: var(--terminal-green-light);">$ ${terminalInput.value}</span>`;
        terminalOutput.appendChild(outputLine);

        if (command === "clear") {
            terminalOutput.innerHTML = "";
        } else if (commands[command]) {
            const responseLine = document.createElement("div");
            responseLine.className = "output-line";
            responseLine.textContent = commands[command];
            terminalOutput.appendChild(responseLine);
        } else if (command) {
            const errorLine = document.createElement("div");
            errorLine.className = "output-line";
            errorLine.style.color = "#FF5555";
            errorLine.textContent = `Command not found: ${command}`;
            terminalOutput.appendChild(errorLine);
        }

        terminalInput.value = "";
        terminalOutput.scrollTop = terminalOutput.scrollHeight;
    }
});
