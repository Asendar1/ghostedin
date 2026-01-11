# GhostdIn

GhostdIn is a lightweight, self-hosted web application for tracking job applications. It provides a clean interface to log, organize, and analyze your job search progress, helping you stay organized during the application process.

Built with modern Go tooling, it offers fast performance, minimal dependencies, and easy deployment via Docker or directly on your machine.

## Features

### Core Functionality
- Create, read, update, and delete job applications
- Track key details: company, role, application date, status, notes, job posting URL, and resume version
- Application status pipeline: Applied, Phone Screen, Interview, Offer, Rejected, No Response
- Basic analytics: total applications, status distribution, applications over time
- Goal setting and progress tracking (e.g., weekly/monthly application targets)

### Technical Highlights
- Interactive UI with partial updates (no full page reloads)
- Responsive design suitable for desktop and mobile
- Single-file SQLite database for persistence
- Single-binary deployment (easy to run anywhere)
