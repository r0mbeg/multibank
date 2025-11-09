import {createFileRoute, useNavigate} from '@tanstack/react-router'
import {Alert, Button, Skeleton, Snackbar, Stack} from "@mui/material";
import {Api} from '../api/api';
import {useEffect, useState} from "react";
import {useQuery} from "@tanstack/react-query";
import type {Consents} from "../types/types.ts";
import PageTitle from "../components/PageTitle.tsx";
import {useDeleteConsent} from "../hooks/useDeleteConsent.ts";
import ConsentItem from "../components/ConsentItem.tsx";

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

    const {mutate: deleteConsent, isError} = useDeleteConsent(
        setSnackbarMessage,
        setSnackbarOpen,
    );

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

    const handleDeleteConsent = (consent_id: string) => {
        deleteConsent({consent_id})
    }

    return (
        <>
            <div className={'flex justify-between items-center'}>
                <PageTitle>Согласия</PageTitle>
                <Button variant={'contained'} onClick={handleAddAccount}>Добавить новый аккаунт</Button>
            </div>

            {isLoading ?
                (
                    <Stack className={'mt-4'}>
                        <Skeleton/>
                        <Skeleton/>
                        <Skeleton/>
                    </Stack>
                ) : (
                    consents && consents.length > 0 ? (
                        <div className={'mt-4'}>
                            {consents.map((consent) => (
                            <ConsentItem key={consent.consent_id} consent={consent}
                                         handleDeleteConsent={handleDeleteConsent}/>
                            ))}
                        </div>
                    ) : (
                        <p>Список согласий пуст</p>
                    )
                )
            }

            <Snackbar
                open={snackbarOpen}
                autoHideDuration={6000}
                onClose={handleSnackbarClose}
                anchorOrigin={{vertical: 'bottom', horizontal: 'right'}}
            >
                <Alert onClose={handleSnackbarClose} severity={`${isError || error ? 'error' : 'success'}`}
                       sx={{width: '100%'}}>
                    {snackbarMessage}
                </Alert>
            </Snackbar>
        </>
    )
}
