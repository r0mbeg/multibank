import { createFileRoute } from '@tanstack/react-router'
import {Alert, Button, Skeleton, Snackbar} from "@mui/material";
import { Api } from '../api/api';
import {useEffect, useState} from "react";
import {useQuery} from "@tanstack/react-query";
import type {Consents} from "../types/types.ts";
import PageTitle from "../components/PageTitle.tsx";

export const Route = createFileRoute('/_auth/consents')({
  component: RouteComponent,
})

const consents = [
    {
        bankName: 'abank',
        consentStatus: 'active',
    },
    {
        bankName: 'bbank',
        consentStatus: 'pending',
    },
    {
        bankName: 'cbank',
        consentStatus: null,
    }
]

function RouteComponent() {
    const [snackbarOpen, setSnackbarOpen] = useState(false); // Состояние для открытия Snackbar
    const [snackbarMessage, setSnackbarMessage] = useState('');

    // const { data: consents, isLoading, error } = useQuery<Consents[]>({
    //     queryKey: ['consents'],
    //     queryFn: async (): Promise<Consents[]> => {
    //         const response = await Api.getConsents();
    //         return response.data;
    //     },
    //     retry: false,
    // })

    //
    // useEffect(() => {
    //     if (error) {
    //         console.log(error)
    //         setSnackbarMessage(error.message || 'Произошла ошибка'); // Устанавливаем сообщение
    //         setSnackbarOpen(true);
    //     }
    // }, [error])


    const handleSnackbarClose = () => {
        setSnackbarOpen(false); // Закрываем Snackbar
    };

  return (
      <>
          <PageTitle>Согласия</PageTitle>
          {/*{isLoading ?*/}
          {/*    (*/}
          {/*        <Skeleton/>*/}
          {/*    ) : consents && consents.length > 0 ? (*/}
          {consents.map((item, idx) => (
                      <div key={idx} className={'flex justify-between border mb-4 shadow-md rounded-md p-4'}>
                          <p className={'text-xl'}>{item.bankName}</p>

                          {item.consentStatus === 'active' && (
                              <>
                                <p>согласие активно</p>
                                <Button variant={'contained'} color={'error'}>отозвать согласие</Button>
                              </>
                          )}

                          {item.consentStatus === 'pending' && (
                              <>
                                  <p className={'text-orange-400'}>согласие в обработке</p>
                              </>
                          )}

                          {item.consentStatus === null && (
                              <>
                                  <Button variant={'contained'} color={'success'}>выпустить согласие</Button>
                              </>
                          )}
                      </div>
                  ))}
              {/*) : (*/}
              {/*    <p>Нет согласий</p>*/}
              {/*)*/}

          <Snackbar
              open={snackbarOpen}
              autoHideDuration={6000}
              onClose={handleSnackbarClose}
              anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
          >
              <Alert onClose={handleSnackbarClose} severity="error" sx={{ width: '100%' }}>
                  {snackbarMessage}
              </Alert>
          </Snackbar>
      </>
  )
}
