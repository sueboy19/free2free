# Project Summary

## Overall Goal
Develop a "buy one get one free" matching website with user authentication via Facebook/Instagram, admin panel for managing activities/locations, user matching functionality, and review system with Swagger API documentation.

## Key Knowledge
- Technology stack: Go 1.25 + Gin framework + GORM + MariaDB + Goth OAuth library
- Development tools: Air for hot reloading, Swagger for API documentation
- Database: MariaDB via Docker Compose with automatic schema migration using GORM
- Authentication: OAuth 2.0 with Facebook and Instagram providers
- Environment configuration: Uses .env files with variables for DB connection, session keys, and OAuth credentials
- Project structure: Modular with separate files for admin, user, organizer, review, and review-like functionality
- Windows-compatible: Avoids Makefile in favor of batch scripts, uses air instead of make for development

## Recent Actions
- Successfully migrated database from MySQL to MariaDB with Docker Compose setup
- Implemented comprehensive GORM-based data models and relationships
- Added Swagger/OpenAPI documentation annotations to all API endpoints
- Configured Air hot-reloading development environment
- Created Windows-compatible batch scripts for building and running the application
- Updated environment variable handling to include DB_HOST configuration
- Fixed numerous compilation issues and dependency conflicts

## Current Plan
1. [DONE] Set up MariaDB database with Docker Compose
2. [DONE] Implement GORM models and auto-migration
3. [DONE] Add Swagger API documentation to all endpoints
4. [DONE] Configure Air hot-reloading development environment
5. [DONE] Create Windows-compatible development scripts
6. [IN PROGRESS] Test OAuth integration with Facebook/Instagram
7. [TODO] Verify all API endpoints function correctly with Swagger UI
8. [TODO] Implement comprehensive error handling and validation
9. [TODO] Add unit and integration tests
10. [TODO] Deploy and test in staging environment

---

## Summary Metadata
**Update time**: 2025-09-21T10:45:25.569Z 
