# Environment Variables

This document describes the environment variables used by the Aitestos frontend application.

## Available Variables

### `VITE_API_BASE_URL`

**Description**: Base URL for the API server. All HTTP requests will be prefixed with this value.

**Default**: `/api/v1`

**Usage**:

```bash
# Development (default)
# Uses relative path which proxies to backend during dev server
yarn dev

# Production
# Build with production API endpoint
VITE_API_BASE_URL=https://api.aitestos.com/api/v1 yarn build
```

**Notes**:

- The value should include the API version path (e.g., `/api/v1`)
- Do not include a trailing slash
- In development, you can use relative paths to leverage Vite's proxy
- For production builds, set this before running `yarn build`

### Example Configurations

#### Local Development

```bash
# .env.local
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

#### Staging Environment

```bash
# .env.staging
VITE_API_BASE_URL=https://staging-api.aitestos.com/api/v1
```

#### Production Environment

```bash
# .env.production
VITE_API_BASE_URL=https://api.aitestos.com/api/v1
```

## Adding New Variables

1. Define the variable in `.env` files (see `.env.example` for reference)
2. Access via `import.meta.env.VITE_*` (Vite automatically exposes variables prefixed with `VITE_`)
3. Update TypeScript types if needed:

```typescript
// vite-env.d.ts
interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  // Add new variables here
}
```

## See Also

- [Vite Environment Variables](https://vitejs.dev/guide/env-and-mode.html)
- [Backend API Documentation](../backend/README.md)
