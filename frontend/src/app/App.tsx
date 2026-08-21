import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { router } from '@app/router'
import { ThemeProvider } from 'next-themes'
import { useState } from 'react'
import { RouterProvider } from 'react-router-dom'
import { Toaster } from '@/shared/ui/sonner'
import { TooltipProvider } from '@/shared/ui/tooltip'

export function App() {
  const [queryClient] = useState(() => new QueryClient())

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider attribute="class" defaultTheme="dark" forcedTheme="dark">
        <TooltipProvider>
          <RouterProvider router={router} />
          <Toaster />
        </TooltipProvider>
      </ThemeProvider>
    </QueryClientProvider>
  )
}
