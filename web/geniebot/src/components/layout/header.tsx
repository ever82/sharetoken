import * as React from 'react'
import { Bot, Menu, Settings, User } from 'lucide-react'
import { cn } from '@/utils'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { WalletButton } from '@/components/wallet/wallet-button'
import { AuthButton } from '@/components/auth/auth-button'

interface HeaderProps {
  className?: string
  onMenuClick?: () => void
}

export function Header({ className, onMenuClick }: HeaderProps) {
  return (
    <header
      className={cn(
        'flex h-16 items-center justify-between border-b bg-background px-4',
        className
      )}
    >
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="icon" onClick={onMenuClick}>
          <Menu className="h-5 w-5" />
        </Button>
        <div className="flex items-center gap-2">
          <Bot className="h-6 w-6 text-primary" />
          <span className="text-xl font-bold">GenieBot</span>
        </div>
      </div>

      <div className="flex items-center gap-2">
        <AuthButton className="hidden sm:flex" />
        <WalletButton />
        <Button variant="ghost" size="icon">
          <Settings className="h-5 w-5" />
        </Button>
      </div>
    </header>
  )
}
