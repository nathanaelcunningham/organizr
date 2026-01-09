import React from 'react'
import { Card } from '../common/Card'

interface ConfigSectionProps {
  title: string
  description?: string
  children: React.ReactNode
}

export const ConfigSection: React.FC<ConfigSectionProps> = ({ title, description, children }) => {
  return (
    <Card className="mb-6">
      <div className="space-y-4">
        <div>
          <h3 className="text-lg font-semibold text-gray-900">{title}</h3>
          {description && <p className="mt-1 text-sm text-gray-600">{description}</p>}
        </div>
        <div className="space-y-4">{children}</div>
      </div>
    </Card>
  )
}
