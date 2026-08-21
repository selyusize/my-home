import { AppLayout } from '@app/layout'
import {
  BitrixPage,
  ExplorerPage,
  GitHubPage,
  GitLabPage,
  GlobeIndexRedirect,
  GlobeLayout,
} from '@pages/globe'
import { HomePage } from '@pages/home'
import { UserPage } from '@pages/user'
import { createHashRouter } from 'react-router-dom'
import { RouteError } from './RouteError'

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
        element: <GlobeLayout />,
        children: [
          {
            index: true,
            element: <GlobeIndexRedirect />,
          },
          {
            path: 'bitrix',
            element: <BitrixPage />,
          },
          {
            path: 'gitlab',
            element: <GitLabPage />,
          },
          {
            path: 'github',
            element: <GitHubPage />,
          },
          {
            path: 'explorer',
            element: <ExplorerPage />,
          },
        ],
      },
      {
        path: 'profile',
        element: <UserPage />,
      },
    ],
  },
])
