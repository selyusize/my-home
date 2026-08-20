import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AppLayout } from '@app/layout'
import { UserPage } from '@pages/user'
import { useState } from 'react'
import { TooltipProvider } from '@/shared/ui'

export function App() {
  const [queryClient] = useState(() => new QueryClient())

  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <AppLayout>
          <UserPage />
        </AppLayout>
      </TooltipProvider>
    </QueryClientProvider>
  )
}
