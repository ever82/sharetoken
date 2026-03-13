import { useState, useEffect, useCallback } from 'react'
import { authService } from '@/services/auth'
import type { AuthState, UserProfile } from '@/types'

export function useAuth() {
  const [state, setState] = useState<AuthState>(authService.getState())

  useEffect(() => {
    const unsubscribe = authService.subscribe((newState) => {
      setState(newState)
    })

    return () => {
      unsubscribe()
    }
  }, [])

  const login = useCallback(async () => {
    return authService.login()
  }, [])

  const logout = useCallback(() => {
    authService.logout()
  }, [])

  const refreshSession = useCallback(async () => {
    return authService.refreshSession()
  }, [])

  const updateProfile = useCallback(async (updates: Partial<UserProfile>) => {
    return authService.updateProfile(updates)
  }, [])

  return {
    state,
    login,
    logout,
    refreshSession,
    updateProfile,
  }
}
