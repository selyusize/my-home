import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { HomePage } from '@pages/home'
import { useState } from 'react'

export function App() {
  const [queryClient] = useState(() => new QueryClient())

  return (
    <QueryClientProvider client={queryClient}>
      <HomePage />
    </QueryClientProvider>
  )
}
