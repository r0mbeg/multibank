import {createFileRoute} from '@tanstack/react-router'
import {useQuery} from "@tanstack/react-query";
import {Api} from "../api/api.ts";
import {useAuthStore} from "../stores/authStore.ts";

export const Route = createFileRoute('/_auth/')({
    component: RouteComponent,
})

function RouteComponent() {
    const {setUser} = useAuthStore()

    useQuery({
        queryKey: ['me'],
        queryFn: async () => {
            const response = await Api.getMe();
            console.log(response.data)
            setUser(response.data)
            return response.data;
        },
        retry: false,
    })

    return (
        <>
            <div>Авторизовались))</div>
        </>
    )
}
