import { get, post } from '@/lib/request'
import type { UserJSON, UserRole } from '@/types/api'

/**
 * Login request payload
 */
interface LoginRequest {
  email: string
  password: string
}

/**
 * Register request payload
 */
interface RegisterRequest {
  username: string
  email: string
  password: string
  role: UserRole
}

/**
 * Authentication API service
 * Handles login, registration, and token refresh
 */
export const authApi = {
  /**
   * User login
   * POST /auth/login
   */
  async login(credentials: LoginRequest): Promise<{
    access_token: string
    refresh_token: string
    user: UserJSON
  }> {
    return post<
      LoginRequest,
      { access_token: string; refresh_token: string; user: UserJSON }
    >('/auth/login', credentials)
  },

  /**
   * User registration
   * POST /auth/register
   */
  async register(userData: RegisterRequest): Promise<UserJSON> {
    return post<RegisterRequest, UserJSON>('/auth/register', userData)
  },

  /**
   * Refresh access token
   * POST /auth/refresh
   */
  async refresh(refreshToken: string): Promise<{
    access_token: string
    refresh_token: string
  }> {
    return post<
      { refresh_token: string },
      { access_token: string; refresh_token: string }
    >('/auth/refresh', { refresh_token: refreshToken })
  },
}
