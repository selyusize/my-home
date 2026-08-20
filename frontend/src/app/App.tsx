import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { router } from '@app/router'
import { useState } from 'react'
import { RouterProvider } from 'react-router-dom'
import { TooltipProvider } from '@/shared/ui'

export function App() {
  const [queryClient] = useState(() => new QueryClient())

  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <RouterProvider router={router} />
      </TooltipProvider>
    </QueryClientProvider>
  )
}
