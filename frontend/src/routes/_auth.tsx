import {createFileRoute, redirect, Outlet} from '@tanstack/react-router'
import {useAuthStore} from "../stores/authStore.ts";

export const Route = createFileRoute('/_auth')({
    beforeLoad: ({location}) => {
        const { user } = useAuthStore.getState();
        if (user === null) {
            throw redirect({
                to: '/login',
                search: {redirect: location.pathname}
            })
        }
    },
    component: () => <Outlet />,
})
