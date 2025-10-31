import { Outlet, createRootRouteWithContext } from '@tanstack/react-router'
import NotFoundPage from "../pages/404.tsx";
// import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'

export const Route = createRootRouteWithContext()({
    component: () => (
        <>
            <Outlet />
            {/*<TanStackRouterDevtools position="bottom-right" initialIsOpen={false} />*/}
        </>
    ),
    notFoundComponent: () => (
        <NotFoundPage />
    )
})