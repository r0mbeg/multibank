import {createFileRoute} from '@tanstack/react-router'
import PageTitle from "../components/PageTitle.tsx";
import {useQuery} from "@tanstack/react-query";
import type {Accounts} from "../types/types.ts";
import {Api} from "../api/api.ts";
import {Skeleton, Stack} from "@mui/material";
import AccountItem from "../components/AccountItem.tsx";

export const Route = createFileRoute('/_auth/accounts')({
    component: RouteComponent,
})

function RouteComponent() {
    const {data: accounts, isLoading, error} = useQuery<Accounts[]>({
        queryKey: ['accounts'],
        queryFn: async (): Promise<Accounts[]> => {
            const response = await Api.getAccounts();
            return response.data;
        },
        retry: false,
    })

    return (
        <>
            <PageTitle>Счета</PageTitle>
            {
                isLoading ? (
                    <Stack className={'mt-4'}>
                        <Skeleton/>
                        <Skeleton/>
                        <Skeleton/>
                    </Stack>
                ) : (
                    accounts && accounts.length > 0 ? (
                        <div className="flex flex-col gap-4">
                            {accounts.map((account) => (
                                    <AccountItem key={account.account_id} account={account}/>
                                )
                            )}
                        </div>
                    ) : (
                        <p>Нет доступных счетов</p>
                    ))
            }

            {
                error && (
                    <p>Непредвиденная ошибка при получении данных о счетах, попробуйте обновить страницу...</p>
                )
            }
        </>
    )
}
