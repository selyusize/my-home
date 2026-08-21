import { isRouteErrorResponse, useRouteError } from 'react-router-dom'

export function RouteError() {
  const error = useRouteError()
  const message = getRouteErrorMessage(error)

  return (
    <div className="flex h-full items-center justify-center bg-background p-6 text-foreground">
      <pre className="max-w-xl whitespace-pre-wrap text-sm text-red-300">{message}</pre>
    </div>
  )
}

function getRouteErrorMessage(error: unknown) {
  if (isRouteErrorResponse(error)) {
    return `${error.status} ${error.statusText}`
  }
  if (error instanceof Error) {
    return error.message
  }
  return 'Unknown error'
}
