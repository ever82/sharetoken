import * as React from 'react'
import { ServiceCard } from './service-card'
import type { ServiceRecommendation } from '@/types'

interface ServiceListProps {
  services: ServiceRecommendation[]
  onInvokeService: (serviceId: string) => void
  isLoading?: boolean
}

export function ServiceList({
  services,
  onInvokeService,
  isLoading = false,
}: ServiceListProps) {
  if (services.length === 0) {
    return null
  }

  return (
    <div className="space-y-4">
      <h3 className="text-sm font-medium text-muted-foreground">
        Recommended Services
      </h3>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {services.map((service) => (
          <ServiceCard
            key={service.id}
            service={service}
            onInvoke={onInvokeService}
            isLoading={isLoading}
          />
        ))}
      </div>
    </div>
  )
}
