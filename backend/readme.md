# Restaurant Backend API (Go + Supabase)

This is a backend service for a restaurant reservation system, written in **Go**, and integrated with **Supabase** for authentication and data storage.

---

## Features

- User Registration & Login (via Supabase Auth)
- Restaurant Upload
- Table Management
- Reservation System

---

## Prerequisites

Make sure you have the following installed:

- [Go (v1.19+)](https://golang.org/dl/)
- Git
- A Supabase Project
- `.env` file with Supabase credentials

---

## Environment Setup

Create a `.env` file in the root directory:

```env
SUPABASE_URL=your_supabase_url
SUPABASE_ANON_KEY=your_supabase_anon_key

## Dependencies

# Initialize the Go module (if not done)
go mod init restaurant-backend

# To automatically install all modules listed in go.mod:

go mod tidy

## Run the Project

# make sure you are in the project root folder:

go run main.go

