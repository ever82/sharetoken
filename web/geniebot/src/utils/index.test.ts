import { describe, it, expect } from 'vitest'
import { cn, formatRelativeTime, truncate } from '@/utils'

describe('cn utility', () => {
  it('merges class names', () => {
    expect(cn('foo', 'bar')).toBe('foo bar')
  })

  it('handles conditional classes', () => {
    const isActive = true
    expect(cn('base', isActive && 'active')).toBe('base active')
  })

  it('handles Tailwind class conflicts', () => {
    expect(cn('px-2 py-1', 'px-4')).toBe('py-1 px-4')
  })
})

describe('formatRelativeTime', () => {
  it('formats seconds ago', () => {
    const now = Math.floor(Date.now() / 1000)
    expect(formatRelativeTime(now)).toBe('0s ago')
  })

  it('formats minutes ago', () => {
    const fiveMinutesAgo = Math.floor(Date.now() / 1000) - 300
    expect(formatRelativeTime(fiveMinutesAgo)).toBe('5m ago')
  })

  it('formats hours ago', () => {
    const twoHoursAgo = Math.floor(Date.now() / 1000) - 7200
    expect(formatRelativeTime(twoHoursAgo)).toBe('2h ago')
  })

  it('formats days ago', () => {
    const threeDaysAgo = Math.floor(Date.now() / 1000) - 259200
    expect(formatRelativeTime(threeDaysAgo)).toBe('3d ago')
  })
})

describe('truncate', () => {
  it('returns original string if shorter than length', () => {
    expect(truncate('hello', 10)).toBe('hello')
  })

  it('truncates string longer than length', () => {
    expect(truncate('hello world', 5)).toBe('hello...')
  })
})
