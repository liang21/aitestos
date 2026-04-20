import { setupServer } from 'msw/node'
import { HttpResponse, http } from 'msw'

export const server = setupServer(
  // Default handlers for common endpoints
  http.post('/api/v1/auth/login', () =>
    HttpResponse.json({
      access_token: 'mock-access-token',
      refresh_token: 'mock-refresh-token',
      user: {
        id: '1',
        username: 'admin',
        email: 'admin@test.com',
        role: 'admin',
      },
    })
  ),
  http.post('/api/v1/auth/refresh', () =>
    HttpResponse.json({
      access_token: 'new-access-token',
      refresh_token: 'new-refresh-token',
    })
  )
)
