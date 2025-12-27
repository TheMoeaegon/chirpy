# Chirpy

A lightweight social media API built with Go that allows users to create accounts, post short messages (chirps), and manage their content.

## What This Project Does

Chirpy is a Twitter-like social media backend API that provides:

- **User Management**: Create accounts, login, and update user information
- **Chirps (Posts)**: Create, read, and delete short messages
- **Authentication**: JWT-based authentication system with refresh tokens
- **Admin Features**: Metrics tracking and user management
- **Premium Features**: Integration with Polka webhooks for premium user upgrades

## Why You Should Care

This project demonstrates:

- **Clean API Design**: RESTful endpoints following best practices
- **Modern Go Development**: Uses current Go idioms and popular libraries
- **Database Integration**: PostgreSQL with proper query management using sqlc
- **Security**: Password hashing with Argon2id and JWT authentication
- **Real-world Features**: Includes pagination, filtering, and webhook integrations

Perfect for learning backend development, understanding API design patterns.

## How to Install and Run

### Prerequisites

- Go 1.24.2 or later
- PostgreSQL database
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Moee1149/chirpy.git
   cd chirpy
