import {createFileRoute} from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/')({
    component: RouteComponent,
})

function RouteComponent() {
    return (
        <>
            <h1 className={'text-bold text-2xl'}>Добро пожаловать в multibank app</h1>
        </>
    )
}
