import {createFileRoute, useNavigate} from '@tanstack/react-router'
import {Alert, Button, Skeleton, Snackbar} from "@mui/material";
import {Api} from '../api/api';
import {useEffect, useState} from "react";
import {useQuery} from "@tanstack/react-query";
import type {Consents} from "../types/types.ts";
import PageTitle from "../components/PageTitle.tsx";
import {useDeleteConsent} from "../hooks/useDeleteConsent.ts";

export const Route = createFileRoute('/_auth/consents')({
    component: RouteComponent,
})

function RouteComponent() {
    const [snackbarOpen, setSnackbarOpen] = useState(false);
    const [snackbarMessage, setSnackbarMessage] = useState('');
    const navigate = useNavigate();

    const {data: consents, isLoading, error} = useQuery<Consents[]>({
        queryKey: ['consents'],
        queryFn: async (): Promise<Consents[]> => {
            const response = await Api.getConsents();
            return response.data;
        },
        retry: false,
    })

    const {mutate: deleteConsent, isPending, isError} = useDeleteConsent(setSnackbarMessage, setSnackbarOpen);


    useEffect(() => {
        if (error) {
            console.log(error)
            setSnackbarMessage(error.message || 'Произошла ошибка');
            setSnackbarOpen(true);
        }
    }, [error])


    const handleSnackbarClose = () => {
        setSnackbarOpen(false); // Закрываем Snackbar
    };

    const handleAddAccount = () => {
        navigate({
            to: '/banks'
        })
    }

    const handleDeleteConsent = (bank_code: string, consent_id: string) => {
        deleteConsent({bank_code, consent_id})
    }

    return (
        <>
            <div className={'flex justify-between items-center'}>
                <PageTitle>Согласия</PageTitle>
                <Button variant={'contained'} onClick={handleAddAccount}>Добавить новый аккаунт</Button>
            </div>

            {isLoading ?
                (
                    <Skeleton/>
                ) : (
                    consents && consents.length > 0 &&
                    (consents.map((item, idx) => (
                        <div key={idx} className={'flex justify-between border mb-4 shadow-md rounded-md p-4'}>
                            <p className={'text-xl'}>{item.bankName}</p>

                            {item.consentStatus === 'active' && (
                                <>
                                    <p>согласие активно</p>
                                    <Button variant={'contained'} color={'error'}
                                            onClick={() => handleDeleteConsent('213', '321')} disabled={isPending}>отозвать
                                        согласие</Button>
                                </>
                            )}

                            {item.consentStatus === 'pending' && (
                                <>
                                    <p className={'text-orange-400'}>согласие в обработке</p>
                                </>
                            )}
                        </div>
                    )))
                )
            }

            <Snackbar
                open={snackbarOpen}
                autoHideDuration={6000}
                onClose={handleSnackbarClose}
                anchorOrigin={{vertical: 'bottom', horizontal: 'right'}}
            >
                <Alert onClose={handleSnackbarClose} severity="error" sx={{width: '100%'}}>
                    {snackbarMessage}
                </Alert>
            </Snackbar>
        </>
    )
}
