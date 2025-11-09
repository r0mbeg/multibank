import { Outlet, createRootRouteWithContext } from '@tanstack/react-router'
import NotFoundPage from "../pages/404.tsx";

export const Route = createRootRouteWithContext()({
    component: () => (
        <>
            <Outlet />
        </>
    ),
    notFoundComponent: () => (
        <NotFoundPage />
    )
})