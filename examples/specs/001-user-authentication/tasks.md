# Tasks

## Sprint 1 - Core Auth

- [x] Set up PostgreSQL schema and migrations
- [x] Implement User repository with CRUD operations
- [x] Build password hashing service (bcrypt, cost=12)
- [x] Implement JWT token generation (RS256)
- [x] Create `POST /auth/register` endpoint
- [x] Create `POST /auth/login` endpoint
- [x] Set up Redis for refresh token storage
- [x] Implement token refresh with rotation
- [ ] Add email verification flow
- [ ] Implement `POST /auth/logout`
- [ ] Write unit tests for auth service (target: 90% coverage)

## Sprint 2 - OAuth & Security

- [ ] Register Google OAuth application
- [ ] Register GitHub OAuth application
- [ ] Implement OAuth authorization code flow
- [ ] Build account linking by verified email
- [ ] Add rate limiting middleware
- [ ] Implement brute force protection (exponential backoff)
- [ ] Add HIBP password check integration
- [ ] Set up concurrent session limits

## Sprint 3 - RBAC

- [ ] Create roles and permissions tables
- [ ] Implement RBAC middleware
- [ ] Build admin API for role management
- [ ] Add permission checking to all protected routes
- [ ] Write integration tests for full auth flows
- [ ] Security audit and penetration testing
- [ ] Documentation and API reference
