# Phase 1: Authentication & Core Infrastructure

## Goal
Establish a secure authentication system and the foundational registry for the POS system.

## Key Features
- [x] **JWT Authentication**: Full implementation of access tokens and secure HttpOnly refresh tokens.
- [x] **User Management**: Secure password hashing (Bcrypt) and standardized profile (Me) endpoint.
- [x] **Base Registry (Unit & Category)**:
    - [x] **Units (UOM)**: Support for multiple units per material with decimal conversion logic.
    - [x] **Categories**: Tax-inclusive grouping for materials and products.
- [x] **Global Response System**: Standardized JSON response format (`success`, `message`, `data`, `pagination`).
- [x] **Media System**: Centralized local storage provider for images and assets.

## Technical Standards
- Clean Architecture (Domain, Usecase, Repository, Delivery layers).
- Global Middleware (JWT, Logger, Recovery, RequestID).
- Multi-environment config via YAML (`.env` replacement).
