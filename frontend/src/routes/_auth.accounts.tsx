import {createFileRoute} from '@tanstack/react-router'
import PageTitle from "../components/PageTitle.tsx";

export const Route = createFileRoute('/_auth/accounts')({
    component: RouteComponent,
})

function RouteComponent() {
    return (
        <PageTitle>Счета</PageTitle>
    )
}
