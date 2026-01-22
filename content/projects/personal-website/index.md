---
title: Personal Website
slug: personal-website
date: 2026-01-01
description: A modern, full-stack personal portfolio website built with Hugo, Go, and Docker
---

This is the website you're currently on! It's a full-stack project showcasing my work, blog posts, and professional information.

⚠️ **Note**: This website is still under heavy development. Expect changes and improvements as I continue to refine the design, features, and architecture.

## Architecture

The site is built with a hybrid architecture combining static site generation with a dynamic API backend:

### Frontend
- **Hugo**: Static site generation for content, blog posts, and project pages
- **Pagefind**: Client-side full-text search without external dependencies
- **Green Nebula Terminal Theme**: Custom theme with retro terminal aesthetics
- **Custom CSS**: Handcrafted styling with terminal aesthetic

### Backend
- **Go API**: RESTful API built with the Chi router
- **Embedded Assets**: Markdown posts and YAML data embedded directly in the binary
- **Clean Architecture**: Layered architecture with asset, data, post, and server layers
- **Health Checks**: Built-in health endpoints for monitoring

### Deployment
- **Docker**: Containerized web server (Caddy) and API services
- **Docker Compose**: Multi-service orchestration for development and production
- **Traefik**: External reverse proxy for HTTPS, routing, and load balancing
- **Ansible**: Infrastructure provisioning and configuration management
- **Rsync**: Static file synchronization for deployments

## Key Features

- **Blog**: Markdown-based blog posts with YAML frontmatter for metadata
- **Search**: Full-text search powered by Pagefind
- **API**: JSON API endpoint for dynamic data access (`/api/data`)
- **Responsive Design**: Mobile-friendly design that works across devices
- **Dark Mode**: Terminal-inspired dark theme
- **Fast**: Static site generation with minimal API calls

## Development

The project uses:
- **Task**: Taskfile for build, test, and deployment automation
- **Docker Compose**: Development environment with hot reload and Traefik
- **Multi-stage Builds**: Efficient Docker builds with builder pattern

## Technologies

- **Languages**: Hugo (templates), Go, TypeScript
- **Build Tools**: Hugo, Go compiler, Docker, Task
- **Infrastructure**: Docker, Docker Compose, Traefik, Caddy, Ansible
- **Deployment**: Rsync, Ansible playbooks

## GitHub

[View the source code on GitHub](https://github.com/jlrickert/jlrickert.me)
