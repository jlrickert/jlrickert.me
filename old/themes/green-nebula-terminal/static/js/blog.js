// Blog modal
const blogPosts = [
  {
    title: "Implementing Repository Pattern in PHP",
    date: "November 15, 2025",
    content:
      `The Repository Pattern is a design pattern that mediates between the domain and data mapping layers. In this post, I'll explore how to implement this pattern in modern PHP applications.\n\nKey Benefits:\n• Separation of concerns\n• Testability improvements\n• Database abstraction\n• Cleaner code organization\n\nImplementation example:\n\ninterface UserRepositoryInterface {\n    public function findById($id);\n    public function save($user);\n    public function delete($user);\n}\n\nclass UserRepository implements UserRepositoryInterface {\n    // Implementation details...\n}\n\nThis pattern allows us to swap implementations easily and write cleaner, more maintainable code.`,
  },
  {
    title: "Building Scalable APIs with Go",
    date: "November 10, 2025",
    content:
      `Go's simplicity and performance make it an excellent choice for building scalable APIs. Here are some best practices I've learned.\n\nKey Principles:\n• Use context for timeouts and cancellation\n• Implement proper error handling\n• Structure your code with clean architecture\n• Use middleware for cross-cutting concerns\n\nExample middleware pattern:\n\nfunc Logger(next http.Handler) http.Handler {\n    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {\n        start := time.Now()\n        next.ServeHTTP(w, r)\n        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))\n    })\n}\n\nWith Go's standard library and some well-chosen packages, you can build production-ready APIs quickly.`,
  },
  {
    title: "DevOps Automation: Streamlining Releases",
    date: "November 3, 2025",
    content:
      `Automating release processes saves time and reduces errors. Here's how I streamlined our deployment pipeline.\n\nTools Used:\n• git-cliff for changelog generation\n• goreleaser for Go binary releases\n• GitHub Actions for CI/CD\n• Semantic versioning for consistency\n\nWorkflow Overview:\n1. Developer creates PR\n2. CI runs tests and linting\n3. On merge to main, auto-tag based on commits\n4. goreleaser builds and publishes\n5. git-cliff generates changelog\n6. Deploy to staging automatically\n\nThis setup reduced our release time from hours to minutes and improved reliability significantly.`,
  },
  {
    title: "Smart Home Network Architecture",
    date: "October 28, 2025",
    content:
      `Designing a secure smart home network requires careful planning. Here's my approach to isolating IoT devices.\n\nNetwork Segmentation:\n• VLAN 10: Trusted devices (computers, phones)\n• VLAN 20: IoT devices (cameras, sensors)\n• VLAN 30: Guest network\n\nSecurity Measures:\n• Firewall rules blocking IoT -> Trusted traffic\n• mDNS reflector for cross-VLAN discovery\n• Pi-hole for DNS filtering\n• Regular firmware updates\n\nHome Assistant Integration:\nHome Assistant sits on the trusted network with firewall rules allowing it to communicate with IoT devices. This gives us control without exposing our main network.\n\nResult: A secure, efficient smart home network that doesn't compromise on functionality.`,
  },
];

const blogModal = document.getElementById("blogModal");
const blogModalContent = document.getElementById(
  "blogModalContent",
);
const modalClose = document.getElementById("modalClose");

document.querySelectorAll(".blog-post").forEach((post) => {
  post.addEventListener("click", () => {
    const postIndex = post.getAttribute("data-post");
    const postData = blogPosts[postIndex];

    blogModalContent.innerHTML = `
                    <div class="section-header">$ cat ~/blog/${
      postData.title.toLowerCase().replace(/\s+/g, "_")
    }.md</div>
                    <div class="divider">--------------------------------------------------------------------------------</div>
                    <div style="color: var(--terminal-green-dark); margin-bottom: 20px;">${postData.date}</div>
                    <h2 style="color: var(--terminal-green-light); margin-bottom: 20px;">${postData.title}</h2>
                    <div style="white-space: pre-line; line-height: 1.8;">${postData.content}</div>
                `;

    blogModal.classList.add("active");
  });
});

modalClose.addEventListener("click", () => {
  blogModal.classList.remove("active");
});

blogModal.addEventListener("click", (e) => {
  if (e.target === blogModal) {
    blogModal.classList.remove("active");
  }
});
