import * as React from 'react'
import { Shield, CheckCircle, XCircle, User } from 'lucide-react'
import { cn } from '@/utils'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { authService } from '@/services/auth'
import { walletService } from '@/services/wallet'
import type { AuthState, UserProfile } from '@/types'

interface AuthButtonProps {
  className?: string
}

export function AuthButton({ className }: AuthButtonProps) {
  const [state, setState] = React.useState<AuthState>({
    isAuthenticated: false,
    isLoading: false,
    did: null,
    user: null,
  })
  const [showAuthDialog, setShowAuthDialog] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)

  React.useEffect(() => {
    // Subscribe to auth state changes
    const unsubscribe = authService.subscribe((newState) => {
      setState(newState)
    })

    // Initial state
    setState(authService.getState())

    return () => {
      unsubscribe()
    }
  }, [])

  const handleLogin = async () => {
    setError(null)

    try {
      // Check if wallet is connected
      if (!walletService.isConnected()) {
        await walletService.connect()
      }

      await authService.login()
      setShowAuthDialog(false)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to authenticate')
    }
  }

  const handleLogout = () => {
    authService.logout()
  }

  const formatDid = (did: string): string => {
    const parts = did.split(':')
    if (parts.length >= 3) {
      const address = parts[parts.length - 1]
      return `did:share:${address.slice(0, 8)}...${address.slice(-8)}`
    }
    return did
  }

  return (
    <>
      {state.isAuthenticated ? (
        <Button
          variant="ghost"
          className={cn('gap-2', className)}
          onClick={() => setShowAuthDialog(true)}
        >
          <Avatar className="h-6 w-6">
            <AvatarImage src={state.user?.avatar} />
            <AvatarFallback>
              <User className="h-3 w-3" />
            </AvatarFallback>
          </Avatar>
          <span className="hidden sm:inline">{formatDid(state.did!)}</span>
          <CheckCircle className="h-4 w-4 text-green-500" />
        </Button>
      ) : (
        <Button
          variant="outline"
          className={cn('gap-2', className)}
          onClick={() => setShowAuthDialog(true)}
        >
          <Shield className="h-4 w-4" />
          Authenticate
        </Button>
      )}

      <Dialog open={showAuthDialog} onOpenChange={setShowAuthDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {state.isAuthenticated ? 'Authenticated' : 'Authenticate'}
            </DialogTitle>
            <DialogDescription>
              {state.isAuthenticated
                ? 'You are authenticated with your DID'
                : 'Sign in with your wallet to authenticate'}
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 pt-4">
            {error && (
              <div className="p-3 bg-destructive/10 text-destructive rounded-md text-sm">
                <div className="flex items-start gap-2">
                  <XCircle className="h-4 w-4 mt-0.5 shrink-0" />
                  {error}
                </div>
              </div>
            )}

            {state.isAuthenticated ? (
              <AuthenticatedView user={state.user!} did={state.did!} />
            ) : (
              <UnauthenticatedView
                onLogin={handleLogin}
                isLoading={state.isLoading}
              />
            )}

            <div className="flex gap-2 pt-4">
              {state.isAuthenticated ? (
                <Button
                  variant="destructive"
                  className="flex-1"
                  onClick={handleLogout}
                >
                  Logout
                </Button>
              ) : (
                <Button
                  className="flex-1"
                  onClick={handleLogin}
                  disabled={state.isLoading}
                >
                  {state.isLoading ? 'Authenticating...' : 'Sign In with Wallet'}
                </Button>
              )}
              <Button
                variant="outline"
                onClick={() => setShowAuthDialog(false)}
              >
                Close
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </>
  )
}

function AuthenticatedView({
  user,
  did,
}: {
  user: UserProfile
  did: string
}) {
  return (
    <div className="space-y-4">
      <div className="flex items-center gap-4 p-4 bg-muted rounded-lg">
        <Avatar className="h-16 w-16">
          <AvatarImage src={user.avatar} />
          <AvatarFallback className="text-2xl">
            <User className="h-8 w-8" />
          </AvatarFallback>
        </Avatar>
        <div>
          <h3 className="font-medium">
            {user.name || 'Anonymous User'}
          </h3>
          <p className="text-sm text-muted-foreground font-mono">
            {did}
          </p>
          <div className="flex gap-2 mt-1">
            <Badge variant="secondary">
              Reputation: {user.reputation.toFixed(2)}
            </Badge>
          </div>
        </div>
      </div>

      <div className="space-y-2">
        <h4 className="font-medium">Wallet Address</h4>
        <code className="block p-2 bg-muted rounded text-sm font-mono break-all">
          {user.address}
        </code>
      </div>

      <div className="space-y-2">
        <h4 className="font-medium">Account Created</h4>
        <p className="text-sm text-muted-foreground">
          {new Date(user.createdAt).toLocaleDateString()}
        </p>
      </div>
    </div>
  )
}

function UnauthenticatedView({
  onLogin,
  isLoading,
}: {
  onLogin: () => void
  isLoading: boolean
}) {
  return (
    <div className="space-y-4">
      <div className="flex flex-col items-center justify-center py-8 text-muted-foreground">
        <Shield className="h-16 w-16 mb-4 opacity-50" />
        <p className="text-center">
          Authenticate with your wallet to access personalized features and
          secure your session.
        </p>
      </div>

      <div className="space-y-2 text-sm text-muted-foreground">
        <div className="flex items-center gap-2">
          <CheckCircle className="h-4 w-4 text-green-500" />
          <span>Sign messages securely</span>
        </div>
        <div className="flex items-center gap-2">
          <CheckCircle className="h-4 w-4 text-green-500" />
          <span>Access your task history</span>
        </div>
        <div className="flex items-center gap-2">
          <CheckCircle className="h-4 w-4 text-green-500" />
          <span>Build reputation on the network</span>
        </div>
      </div>

      <p className="text-xs text-muted-foreground text-center">
        You need to connect your wallet first before authenticating.
      </p>
    </div>
  )
}

export function AuthGuard({
  children,
  fallback,
}: {
  children: React.ReactNode
  fallback?: React.ReactNode
}) {
  const [isAuthenticated, setIsAuthenticated] = React.useState(false)
  const [isLoading, setIsLoading] = React.useState(true)

  React.useEffect(() => {
    const unsubscribe = authService.subscribe((state) => {
      setIsAuthenticated(state.isAuthenticated)
      setIsLoading(false)
    })

    setIsAuthenticated(authService.isAuthenticated())
    setIsLoading(false)

    return () => unsubscribe()
  }, [])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="animate-spin h-8 w-8 border-2 border-primary border-t-transparent rounded-full" />
      </div>
    )
  }

  if (!isAuthenticated) {
    return (
      fallback || (
        <div className="flex flex-col items-center justify-center p-8 space-y-4">
          <Shield className="h-12 w-12 text-muted-foreground" />
          <p className="text-muted-foreground">
            Please authenticate to access this feature
          </p>
          <AuthButton />
        </div>
      )
    )
  }

  return <>{children}</>
}
