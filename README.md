# TableToppers

A web application designed for efficient restaurant table reservations and management.

---

## Table of Contents

1.  [Key Features](#key-features)
2.  [Team](#team)
3.  [Technology Stack](#technology-stack)
4.  [Sprint Progress Summary](#sprint-progress-summary)
    * [Sprint 1: Foundation & Initial User Features](#sprint-1-foundation--initial-user-features)
    * [Sprint 2: User Authentication & Core Models](#sprint-2-user-authentication--core-models)
    * [Sprint 3: Core Management & Reservation Features](#sprint-3-core-management--reservation-features)
    * [Sprint 4: Waitlist, Business Portal & Testing](#sprint-4-waitlist-business-portal--testing)
5.  [Getting Started](#getting-started)
    * [Prerequisites](#prerequisites)
    * [Setup](#setup)
6.  [Running the Application](#running-the-application)
7.  [Testing](#testing)
8.  [Docker (Frontend)](#docker-frontend)
9.  [Development Tools & Features (Frontend)](#development-tools--features-frontend)
10. [License](#license)

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

## Team

* **Alex Hu** - Backend Development
* **Adit Potta** - Backend Development
* **Hongyu Chen** - Frontend Development
* **Ian Arnold** - Frontend Development

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

## Sprint Progress Summary

This section outlines the major accomplishments from each development sprint.

### Sprint 1: Foundation & Initial User Features

* **Focus:** Setting up the development environment and implementing basic user registration and login.
* **Completed:**
    * Backend and Frontend environment setup (Node, Go, Angular CLI, Git, Supabase integration basics).
    * Initial Backend database schema design.
    * Frontend landing page development.
    * Basic structure for User Registration and Login features (partially complete due to auth complexities).
* **Key User Stories Addressed:** Environment Setup, Landing Page Development, Initial work on Register/Login.

### Sprint 2: User Authentication & Core Models

* **Focus:** Completing the user authentication flow, establishing core backend data models, enabling basic restaurant Browse, and defining user roles.
* **Completed:**
    * **Backend:**
        * Finalized secure User Registration and Login implementation (password hashing, token generation/validation, session handling).
        * Implemented foundational Role-Based Access Control (RBAC) logic (defining roles like Customer, Staff, Manager).
        * Defined and implemented core Go structs/models in the backend for Restaurants, Tables, and Users based on Sprint 1 schema design.
        * Created the initial API endpoint to retrieve a list of all restaurants.
    * **Frontend:**
        * Completed User Registration and Login forms with proper input validation and error handling.
        * Developed basic placeholder pages for logged-in users (e.g., customer profile/dashboard).
        * Implemented a frontend page to display the list of restaurants fetched from the backend API.
        * Refined landing page styles and navigation flows post-login.
* **Key User Stories Addressed:** Fully addressed Register a New User, Login for Existing User; Started Retrieve a List of Restaurants; Established foundations for role-specific access.

### Sprint 3: Core Management & Reservation Features

* **Focus:** Building out core functionalities for restaurant/table management and the customer reservation process.
* **Completed:**
    * **Backend:**
        * Implemented full CRUD API routes for Restaurants and Tables.
        * Implemented API routes for creating and retrieving Reservations per restaurant.
        * Added Unit Tests for Restaurant and Table routes using mocked data.
        * Resolved backend integration issues.
        * Verified all new endpoints using Postman.
    * **Frontend:**
        * Enhanced UI/Navigation (navbar styling, fixed buttons).
        * Added new pages: "Info" page and "Business Portal".
        * Implemented Manager/Admin features: Edit tables, view table status, view/add/check-in/cancel reservations.
        * Implemented Customer features: Create and view personal reservations.
* **Key User Stories Addressed:** Retrieve Restaurants (filtering), Create Restaurant, Get/Create Tables, Create/Get Reservations (All/Single), various UI/UX improvements for staff/managers.

### Sprint 4: Waitlist, Business Portal & Testing

* **Focus:** Implementing the waitlist feature, enhancing the business-facing portal, and significantly increasing test coverage.
* **Completed:**
    * **Backend:**
        * Implemented full GET, POST, DELETE API routes for the Waitlist feature.
        * Added comprehensive Unit Tests for Waitlist, Reservation, and remaining Table/Restaurant routes.
        * Created a detailed backend-specific README file.
        * Finalized API implementation (excluding planned Menu features), documented in `openapi.yaml`.
    * **Frontend:**
        * Developed the Waitlist management component for staff (add customer, assign cohort, estimate wait time, assign table, cancel).
        * Updated customer-facing pages (new arrivals content, footer links).
        * Enhanced Business Portal (added demo GIF, modals for contact/demo/trial, pricing tiers, FAQ section).
        * Added extensive Cypress (E2E) tests for business portal interactions, footer links, etc.
        * Implemented/Extended Unit Tests (Jest) for Reservation, Restaurant, Table, and Waitlist components.
* **Key User Stories Addressed:** Join Waitlist, Manage Waitlist (implied staff story), Business Portal enhancements, extensive testing implementation.

---

## Getting Started

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
    * The backend API server should start, listening on a configured port (check backend code/config for the specific port).

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
