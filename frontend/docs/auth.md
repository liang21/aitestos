# Authentication & Security

## Token Management

### JWT Token Flow

1. **Login**: User receives `access_token` and `refresh_token`
2. **Storage**: Tokens stored in `localStorage`
   - `access_token`: Short-lived JWT for API calls
   - `refresh_token`: Long-lived token for obtaining new access tokens

3. **Request Interceptor** (`src/lib/request.ts`)
   - Automatically adds `Authorization: Bearer <token>` header to all API requests
   - Skips auth endpoints (login, refresh)

### Token Refresh Logic

When an API call returns **401 Unauthorized**:

1. If a refresh is already in progress, request is queued
2. Uses `refresh_token` to obtain new tokens
3. Updates stored tokens
4. Retries original request with new token
5. If refresh fails, user is redirected to login (via `setAuthExpiredHandler`)

### CSRF Protection

**Status**: ✅ **Not Required**

Backend uses JWT authentication (stateless), not session-based CSRF tokens.

- No CSRF token needed for API requests
- JWT tokens provide sufficient protection against CSRF attacks

### Security Best Practices

| Practice         | Implementation                                     |
| ---------------- | -------------------------------------------------- |
| HTTPS Required   | ✅ Enforce in production                           |
| Token Storage    | localStorage (consider httpOnly cookies in future) |
| Token Expiration | Automatic refresh on 401                           |
| Token Revocation | Handled via 401 → refresh flow                     |
| Public Routes    | `/login`, `/register` (no token required)          |

## API Security Headers

All requests include:

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

## Error Handling

| Status Code | Handling                                  |
| ----------- | ----------------------------------------- |
| 401         | Token refresh attempt                     |
| 403         | Permission denied                         |
| 404         | Resource not found                        |
| 409         | Conflict                                  |
| 500         | Server error (show user-friendly message) |
