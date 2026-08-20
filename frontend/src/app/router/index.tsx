import { AppLayout } from '@app/layout'
import { GlobePage } from '@pages/globe'
import { HomePage } from '@pages/home'
import {
  createHashRouter,
  isRouteErrorResponse,
  useRouteError,
} from 'react-router-dom'

function RouteError() {
  const error = useRouteError()
  const message = isRouteErrorResponse(error)
    ? `${error.status} ${error.statusText}`
    : error instanceof Error
      ? error.message
      : 'Unknown error'

  return (
    <div className="flex h-full items-center justify-center bg-[#06070f] p-6 text-white">
      <pre className="max-w-xl whitespace-pre-wrap text-sm text-red-300">{message}</pre>
    </div>
  )
}

export const router = createHashRouter([
  {
    path: '/',
    element: <AppLayout />,
    errorElement: <RouteError />,
    children: [
      {
        index: true,
        element: <HomePage />,
      },
      {
        path: 'globe',
        element: <GlobePage />,
      },
    ],
  },
])
