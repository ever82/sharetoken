import { useState, useEffect, useCallback } from 'react'
import { walletService } from '@/services/wallet'
import type { WalletState } from '@/types'

export function useWallet() {
  const [state, setState] = useState<WalletState>(walletService.getState())

  useEffect(() => {
    const unsubscribe = walletService.subscribe((newState) => {
      setState(newState)
    })

    return () => {
      unsubscribe()
    }
  }, [])

  const connect = useCallback(async () => {
    return walletService.connect()
  }, [])

  const disconnect = useCallback(() => {
    walletService.disconnect()
  }, [])

  const refreshBalance = useCallback(async () => {
    return walletService.refreshBalance()
  }, [])

  const sendTokens = useCallback(
    async (recipient: string, amount: string, memo?: string) => {
      return walletService.sendTokens(recipient, amount, memo)
    },
    []
  )

  return {
    state,
    connect,
    disconnect,
    refreshBalance,
    sendTokens,
  }
}
