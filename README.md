# **Fitness Locator Applicatione**

## Overview

Welcome to the Fitness Locator application! This application helps users find local gyms, fitness centers, sports clubs, yoga studios, and more based on their postcode and radius. The application also allows users to view detailed information about fitness-related activities, including contact details, facilities, and opening hours. The platform promotes a community-focused experience with a clean, modern interface.

## Features

Find Fitness Locations: Search for gyms, yoga studios, sports clubs, and other fitness-related locations by postcode and radius.
Activity Details: View detailed information about each fitness-related activity, including the name, description, facilities, and contact details.
User Authentication: Users can log in to manage their activities, and admins can manage all activities.
Map Integration: A map displays the fitness locations with different icons for different types of fitness centers.
Community Feel: The application encourages interaction and provides users with workout ideas and guidance.

## Technologies Used

Backend: Go, Gin, GORM, PostgreSQL
Frontend: React, JSX, Tailwind CSS
Authentication: JWT (JSON Web Tokens) for secure authentication
Database: PostgreSQL for storing user and activity data

## Project Structure

```bash
fitness-locator/
│
├── backend/                        # Contains the backend logic
│   ├── controllers/                # Controller files for handling routes
│   ├── models/                     # Database models for users and activities
│   ├── routes/                     # Route definitions
│   ├── middleware/                 # Middleware for authentication
│   └── main.go                     # Main entry point for the backend server
│
├── frontend/                       # Contains the React frontend
│   ├── src/
│   │   ├── components/             # React components
│   │   ├── pages/                  # React pages (for different routes)
│   │   ├── controllers/            # Functions for API calls and logic
│   │   ├── App.js                  # Main entry point for React
│   │   └── assets/                 # TailwindCSS styling & Images
│   └── package.json                # Frontend dependencies and scripts
│
├── database/                       # Database migration files
│   └── migration.sql               # SQL migrations for PostgreSQL
│
├── .env                            # Environment variables
├── .gitignore                      # Git ignore file
└── README.md                       # Project documentation (this file)
```

## Installation

To run the Fitness Locator application locally, follow these steps:

## 1. Backend Setup

a. Change your directory to where you wish to run this script and store the cloned repository:

```bash
cd <filename>
```

b. Clone the repository:

```bash
git clone https://github.com/laurawarren88/go_spa_backend_fitness_locator.git
cd go_spa_backend_fitness_locator
```

c. Install the required dependencies:

```bash
go mod tidy
```

d. Set up the PostgreSQL database:
Create a database for the application.
Configure the database connection in the .env file as below:

```text
ENV=development
DEV_DOMAIN=http://localhost:<port number for frontend>
DEV_SECURE_COOKIE=false
DEV_HTTP_ONLY_COOKIE=false

PROD_DOMAIN=http://<production IP address>
PROD_SECURE_COOKIE=true
PROD_HTTP_ONLY_COOKIE=true

ACCESS_SECRET_KEY=<Set a secret key for JWT access token> 
REFRESH_SECRET_KEY=<Set a secret key for JWT refresh token>

PORT=<port number for backend>

DB_HOST=localhost
DB_USER=<DB variables>
DB_PASSWORD=<<DB variables>
DB_NAME=<DB variables>
DB_PORT=<DB variables>

ADMIN_PASSWORD=<password to set the admin user>
```

e. Run the backend server:

```bash
go run main.go
```

## 2. Frontend Setup

a. Change your directory to where you wish to run this script and store the cloned repository:

```bash
git clone https://github.com/laurawarren88/react_frontend_fitness_locator.git
cd <filename>
```

b. Install the required dependencies:

```bash
npm install
npm start
```

c. Open your web browser to access the Fitness Locator application and navigate to:

```text
http://localhost:<port number> 
```

## Usage

1. Search for Fitness Locations: Enter a postcode and select a radius to search for nearby fitness-related locations.
2. View Activity Details: Click on a specific fitness location to view its detailed information, including description, facilities, contact details, and opening hours.
3. Log in to Manage Activities: Users can log in to create, update, or delete activities they own. Admin users can manage all activities.
