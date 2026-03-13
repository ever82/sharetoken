import * as React from 'react'
import { Wallet, LogOut, ChevronDown, AlertCircle } from 'lucide-react'
import { cn } from '@/utils'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { walletService, WalletService } from '@/services/wallet'
import type { WalletState } from '@/types'

interface WalletButtonProps {
  className?: string
}

export function WalletButton({ className }: WalletButtonProps) {
  const [state, setState] = React.useState<WalletState>({
    address: null,
    balance: '0',
    isConnected: false,
    isConnecting: false,
    chainId: 'sharetokens-1',
  })
  const [showInstallDialog, setShowInstallDialog] = React.useState(false)
  const [error, setError] = React.useState<string | null>(null)
  const [dropdownOpen, setDropdownOpen] = React.useState(false)

  React.useEffect(() => {
    // Subscribe to wallet state changes
    const unsubscribe = walletService.subscribe((newState) => {
      setState(newState)
    })

    // Initial state
    setState(walletService.getState())

    return () => {
      unsubscribe()
    }
  }, [])

  const handleConnect = async () => {
    setError(null)

    try {
      // Check if Keplr is installed
      const isInstalled = await walletService.isKeplrInstalled()
      if (!isInstalled) {
        setShowInstallDialog(true)
        return
      }

      await walletService.connect()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to connect wallet')
    }
  }

  const handleDisconnect = () => {
    walletService.disconnect()
    setDropdownOpen(false)
  }

  const handleRefreshBalance = async () => {
    await walletService.refreshBalance()
  }

  const formatAddress = (address: string): string => {
    return `${address.slice(0, 8)}...${address.slice(-8)}`
  }

  return (
    <>
      {state.isConnected ? (
        <div className="relative">
          <Button
            variant="outline"
            className={cn('gap-2', className)}
            onClick={() => setDropdownOpen(!dropdownOpen)}
          >
            <Wallet className="h-4 w-4" />
            <span className="hidden sm:inline">{formatAddress(state.address!)}</span>
            <span className="sm:hidden">{formatAddress(state.address!).slice(0, 6)}</span>
            <ChevronDown className="h-3 w-3" />
          </Button>

          {dropdownOpen && (
            <div className="absolute right-0 top-full mt-2 w-64 rounded-lg border bg-background shadow-lg z-50">
              <div className="p-4 space-y-3">
                <div className="space-y-1">
                  <p className="text-xs text-muted-foreground">Address</p>
                  <p className="text-sm font-mono break-all">{state.address}</p>
                </div>
                <div className="space-y-1">
                  <p className="text-xs text-muted-foreground">Balance</p>
                  <div className="flex items-center justify-between">
                    <p className="text-sm font-medium">{state.balance} STT</p>
                    <button
                      onClick={handleRefreshBalance}
                      className="text-xs text-primary hover:underline"
                    >
                      Refresh
                    </button>
                  </div>
                </div>
                <div className="border-t pt-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    className="w-full justify-start gap-2 text-destructive"
                    onClick={handleDisconnect}
                  >
                    <LogOut className="h-4 w-4" />
                    Disconnect
                  </Button>
                </div>
              </div>
            </div>
          )}

          {/* Close dropdown when clicking outside */}
          {dropdownOpen && (
            <div
              className="fixed inset-0 z-40"
              onClick={() => setDropdownOpen(false)}
            />
          )}
        </div>
      ) : (
        <Button
          onClick={handleConnect}
          disabled={state.isConnecting}
          className={cn('gap-2', className)}
        >
          <Wallet className="h-4 w-4" />
          {state.isConnecting ? 'Connecting...' : 'Connect Wallet'}
        </Button>
      )}

      {error && (
        <div className="absolute right-0 top-full mt-2 p-3 bg-destructive/10 border border-destructive rounded-lg text-sm text-destructive max-w-xs">
          <div className="flex items-start gap-2">
            <AlertCircle className="h-4 w-4 mt-0.5 shrink-0" />
            {error}
          </div>
        </div>
      )}

      <Dialog open={showInstallDialog} onOpenChange={setShowInstallDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Install Keplr Wallet</DialogTitle>
            <DialogDescription>
              To use GenieBot, you need to install the Keplr wallet extension.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-4 pt-4">
            <p className="text-sm text-muted-foreground">
              Keplr is a secure wallet for Cosmos-based blockchains. It allows you to:
            </p>
            <ul className="text-sm space-y-1 list-disc list-inside text-muted-foreground">
              <li>Store and manage your STT tokens</li>
              <li>Sign transactions securely</li>
              <li>Authenticate with the GenieBot network</li>
            </ul>
            <div className="flex gap-2 pt-4">
              <Button
                className="flex-1"
                onClick={() => {
                  window.open('https://www.keplr.app/', '_blank')
                  setShowInstallDialog(false)
                }}
              >
                Install Keplr
              </Button>
              <Button
                variant="outline"
                onClick={() => setShowInstallDialog(false)}
              >
                Cancel
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </>
  )
}
