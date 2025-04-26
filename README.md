# TableToppers

A web application designed for efficient restaurant table reservations and management.

---

## Table of Contents

1.  [Project Goals](#project-goals)
2.  [Key Features](#key-features)
3.  [Technology Stack](#technology-stack)
4.  [Key Dependencies](#key-dependencies)
5.  [Getting Started](#getting-started)
    * [Prerequisites](#prerequisites)
    * [Setup](#setup)
6.  [Running the Application](#running-the-application)
7.  [Testing](#testing)
8.  [Docker (Frontend)](#docker-frontend)
9.  [API Documentation](#api-documentation)
10. [Architecture Overview](#architecture-overview)
11. [Configuration](#configuration)
12. [Sprint Progress Summary](#sprint-progress-summary)
    * [Sprint 1: Foundation & Initial User Features](#sprint-1-foundation--initial-user-features)
    * [Sprint 2: User Authentication & Core Models](#sprint-2-user-authentication--core-models)
    * [Sprint 3: Core Management & Reservation Features](#sprint-3-core-management--reservation-features)
    * [Sprint 4: Waitlist, Business Portal & Testing](#sprint-4-waitlist-business-portal--testing)
13. [Team & Contributions](#team--contributions)
14. [Future Work](#future-work)
15. [Contributing](#contributing)
16. [License](#license)

---

## Project Goals

TableToppers aims to streamline the often frustrating process of restaurant discovery and reservation. We address common pain points for both customers (difficulty finding available tables, opaque waitlist processes) and restaurants (inefficient table management, manual reservation tracking). Our platform provides a seamless interface for Browse, booking, waitlisting, and managing restaurant seating, enhancing convenience and efficiency for everyone involved.

---

## Key Features

**Restaurant Management:**

* Create, view, and modify the layout of available tables, defining details like guest capacity.
* View, edit, create, and assign restaurant reservations to specific tables.
* Manage the restaurant menu (create, view, edit) for customer visibility. *(Future Enhancement)*
* Adjust reservation availability based on holidays, private events, or operational changes.

**Restaurant Employees:**

* View a comprehensive table layout indicating occupied, available, and reserved tables.
* Edit the layout to accommodate walk-in customers and call-in reservations.
* Estimate customer wait times based on current occupancy and table status.
* Enroll walk-in customers in a digital waiting list with notifications when their table is nearly ready.
* Manage waitlist entries (add, view, assign table, remove).

**Customers:**

* Select a restaurant to view its information and menu.
* Easily create, view, and modify their own reservations.
* Secure user authentication and login system.
* Join a waitlist if the restaurant is fully booked.
* Receive restaurant recommendations based on location or past reservations. *(Future Enhancement)*

---

## Technology Stack

* **Backend:** Go, Supabase (Authentication, Database)
* **Frontend:** Angular (v19+), Angular Material, TypeScript
* **Testing:** Jest (Unit), Cypress (End-to-End)
* **Containerization:** Docker (for Frontend)
* **Development Tools (Frontend):** Leverages features from the [wlucha/angular-starter](https://github.com/wlucha/angular-starter) template, including:
    * Transloco (Internationalization)
    * Compodoc (Auto Documentation)
    * Storybook (Component Examples)
    * Source Map Explorer (Bundle Analysis)
    * ESLint, Prettier, Commit Linting (Code Quality)
    * AuditJS (Security Audits)
    * Auto-Changelog

---

## Key Dependencies

This project relies on several key libraries and frameworks managed by Go modules and npm.

### Backend (Go)

Dependencies are managed using Go Modules (`go.mod` and `go.sum`). Key dependencies include:

* **Go Standard Library:** Utilized extensively for core functionalities (HTTP handling, JSON processing, etc.).
* **Supabase Go Client (`supabase-go`):** Used for interacting with the Supabase backend (authentication, database operations). *(Verify exact library name if different)*
* **HTTP Router:** A router library (e.g., `net/http`'s `ServeMux`, `gorilla/mux`, `chi`, or `gin-gonic`) is used to handle API request routing. *(Specify which one if applicable)*
* **Testing Libraries:** Standard Go testing package, potentially supplemented by assertion libraries (e.g., `stretchr/testify`).

*For a complete list of backend dependencies and their versions, please refer to the `go.mod` file in the backend directory.*

### Frontend (Angular / npm)

Dependencies are managed using npm (`package.json` and `package-lock.json`). Key dependencies include:

* **Angular Framework (`@angular/core`, `@angular/common`, `@angular/router`, etc.):** The core framework for building the SPA.
* **Angular Material (`@angular/material`):** Component library for UI elements following Material Design principles.
* **RxJS:** Library for reactive programming using Observables, heavily used within Angular.
* **TypeScript:** The primary language used for frontend development.
* **Transloco (`@ngneat/transloco`):** Library for internationalization (i18n).
* **Jest (`jest`):** Framework for running unit tests.
* **Cypress (`cypress`):** Framework for running end-to-end tests.
* **Storybook (`@storybook/angular`):** Tool for developing and showcasing UI components in isolation.

*For a complete list of frontend dependencies (including development dependencies) and their versions, please refer to the `package.json` file in the frontend directory.*

---

## Getting Started

Follow these instructions to set up and run the project locally.

### Prerequisites

Before you begin, ensure you have the following installed:

* [Go (v1.19+)](https://golang.org/dl/)
* [Node.js and npm](https://nodejs.org/) (LTS version recommended)
* Git
* A Supabase Project (Sign up at [supabase.com](https://supabase.com/))
* Docker (Optional, for running the frontend in a container)

### Setup

1.  **Clone the Repository:**
    ```bash
    # !!! Replace <your-tabletoppers-repository-url> with the actual URL of your repo !!!
    git clone <your-tabletoppers-repository-url>
    cd <your-tabletoppers-directory> # Navigate into the cloned project directory
    ```

2.  **Backend Setup:**
    * Navigate to the backend project directory (e.g., `cd backend` or stay in root if it's there. **Adjust path as needed.**).
    * Create a `.env` file in the backend root directory:
        ```env
        SUPABASE_URL=your_supabase_url
        SUPABASE_ANON_KEY=your_supabase_anon_key
        ```
        *(Replace `your_supabase_url` and `your_supabase_anon_key` with your actual Supabase project credentials).*
    * Initialize Go modules (if `go.mod` doesn't exist) and install dependencies:
        ```bash
        # Run only if go.mod does not exist:
        # go mod init <your-backend-module-path> # e.g., [github.com/your-username/tabletoppers/backend](https://github.com/your-username/tabletoppers/backend)

        # Install/update dependencies
        go mod tidy
        ```

3.  **Frontend Setup:**
    * Navigate to the frontend project directory (e.g., `cd ../frontend` or `cd frontend`. **Adjust path as needed.**).
    * Install Node.js dependencies:
        ```bash
        npm install
        ```

---

## Running the Application

You'll typically need to run both the backend and frontend services simultaneously.

1.  **Run the Backend Service (Go + Supabase):**
    * Ensure you are in the backend directory.
    * Execute the main Go program:
        ```bash
        go run main.go
        ```
    * The backend API server should start, listening on a configured port (check backend code/config for the specific port, often `8080`).

2.  **Run the Frontend Application (Angular):**
    * Ensure you are in the frontend directory.
    * Start the Angular development server:
        ```bash
        npm run start
        ```
    * Open your web browser and navigate to `http://localhost:4200` (or the address provided in the console output).

---

## Testing

The project includes configurations for unit and end-to-end testing:

* **Backend Unit Tests (Go):**
    * Run tests using the standard Go testing tools from the backend directory:
        ```bash
        go test ./...
        ```
* **Frontend Unit Tests (Jest):**
    ```bash
    # Run from the frontend directory
    npm run test
    ```
* **Frontend End-to-End Tests (Cypress):**
    ```bash
    # Run from the frontend directory (check package.json for exact command if different)
    npm run e2e
    # Or to open the Cypress test runner UI:
    # npm run cy:open
    ```

---

## Docker (Frontend)

You can build and run the frontend application using Docker.

```bash
# Navigate to the frontend directory

# Build the Docker image (tagging it as 'tabletoppers-frontend')
docker build . -t tabletoppers-frontend

# Run the Docker container, mapping host port 3000 to container port 80
docker run -p 3000:80 tabletoppers-frontend
