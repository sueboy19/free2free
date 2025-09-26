# Project Summary

## Overall Goal
The goal is to develop a "buy one get one free" matching website with user authentication via Facebook/Instagram, admin panel for managing activities/locations, user matching functionality, and review system with Swagger API documentation.

## Key Knowledge
- **Technology Stack**: Go 1.25 + Gin framework + GORM + MariaDB + Goth OAuth library
- **Development Tools**: Air for hot reloading, Swagger for API documentation
- **Database**: MariaDB via Docker Compose with automatic schema migration using GORM
- **Authentication**: OAuth 2.0 with Facebook and Instagram providers
- **Environment Configuration**: Uses .env files with variables for DB connection, session keys, and OAuth credentials
- **Project Structure**: Modular with separate files for admin, user, organizer, review, and review-like functionality
- **Windows Compatibility**: Avoids Makefile in favor of batch scripts, uses air instead of make for development
- **API Documentation**: Comprehensive Swagger/OpenAPI documentation covering all endpoints and data models
- The project had redundant JWT-related functions (`generateJWT`/`validateJWT`) and structs (`JwtClaims`) that were removed in favor of `generateJWTToken`/`validateJWTToken` and `Claims` struct to eliminate code duplication

## Recent Actions
- Successfully migrated database from MySQL to MariaDB with Docker Compose setup
- Implemented comprehensive GORM-based data models and relationships
- Added Swagger/OpenAPI documentation annotations to all API endpoints
- Configured Air hot-reloading development environment
- Created Windows-compatible batch scripts for building and running the application
- Updated environment variable handling to include DB_HOST configuration
- Fixed numerous compilation issues and dependency conflicts
- Identified and removed redundant JWT functions (`generateJWT`, `validateJWT`) and struct (`JwtClaims`)
- Standardized on single JWT implementation using `generateJWTToken`, `validateJWTToken`, and `Claims` struct with IsAdmin field
- Updated JWT token generation to include IsAdmin field in claims

## Current Plan
1.  [DONE] Set up MariaDB database with Docker Compose
2.  [DONE] Implement GORM models and auto-migration
3.  [DONE] Add Swagger API documentation to all endpoints
4.  [DONE] Configure Air hot-reloading development environment
5.  [DONE] Create Windows-compatible development scripts
6.  [DONE] Analyze code logic and remove redundant code
7.  [DONE] Run tests to ensure code integrity after removals
8.  [TODO] Verify all API endpoints function correctly with Swagger UI
9.  [TODO] Implement comprehensive error handling and validation
10. [TODO] Add unit and integration tests
11. [TODO] Deploy and test in staging environment

关于用户在Swagger中使用Facebook登录的问题，这需要实现一个特殊的认证机制，因为Swagger UI本身不能直接处理OAuth重定向。通常的做法是：
1. 在Swagger中添加一个API密钥认证选项
2. 用户先通过网站前端完成Facebook登录，获取JWT token或session
3. 将token/session ID手动输入到Swagger UI的认证字段中
4. Swagger会在后续请求中将该token作为Authorization header发送

这需要在后端实现相应的JWT token生成和验证机制，或者允许Swagger直接使用session ID进行认证。

---

## Summary Metadata
**Update time**: 2025-09-26T16:42:11.831Z 
