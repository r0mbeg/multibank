import {createFileRoute, redirect, Outlet} from '@tanstack/react-router'
import {useAuthStore} from "../stores/authStore.ts";
import Header from "../components/Header.tsx";
import Sidebar from "../components/Sidebar.tsx";

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
    component: () => (
        <>
            <Header/>
            <div className={'flex flex-1 gap-4'}>
                <Sidebar/>
                <main className={'bg-white w-full shadow-md rounded-xl p-4'}>
                    <Outlet/>
                </main>
            </div>
        </>
    ),
})
