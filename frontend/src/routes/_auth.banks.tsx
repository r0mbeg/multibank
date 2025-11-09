import {createFileRoute} from '@tanstack/react-router'
import {useQuery} from "@tanstack/react-query";
import type {Bank, LoginBankForm} from "../types/types.ts";
import {Api} from "../api/api.ts";
import PageTitle from "../components/PageTitle.tsx";
import React, {useState} from "react";
import {
    Alert,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    FormControl, IconButton, InputAdornment, Snackbar,
    TextField
} from "@mui/material";
import {Controller, type SubmitErrorHandler, type SubmitHandler, useForm} from "react-hook-form";
import {Visibility, VisibilityOff} from "@mui/icons-material";
import {zodResolver} from "@hookform/resolvers/zod";
import {loginBankSchema} from "../types/schemas.ts";
import {useNewConsent} from "../hooks/useNewConsent.ts";

export const Route = createFileRoute('/_auth/banks')({
    component: RouteComponent,
})

function RouteComponent() {
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const [showPassword, setShowPassword] = useState<boolean>(false);
    const [snackbarOpen, setSnackbarOpen] = useState(false); // Состояние для открытия Snackbar
    const [snackbarMessage, setSnackbarMessage] = useState<string>('');
    const [selectedBank, setSelectedBank] = useState<string>('');

    const {mutate: newConsent, isPending, isError} = useNewConsent(setSnackbarMessage, setSnackbarOpen)

    const {
        control,
        register,
        handleSubmit,
        reset,
        formState: {errors},
    } = useForm<LoginBankForm>({
        resolver: zodResolver(loginBankSchema)
    })

    const {data: banks, isLoading, error} = useQuery<Bank[]>({
        queryKey: ['banks'],
        queryFn: async (): Promise<Bank[]> => {
            const response = await Api.getBanks();
            return response.data;
        },
        retry: false,
    })

    const handleClickOpenConsent = (bankCode: string) => {
        setSelectedBank(bankCode);
        setIsOpen(true);
    };

    const handleClose = () => {
        reset();
        setIsOpen(false);
    };

    const onSubmit: SubmitHandler<LoginBankForm> = (data) => {
        newConsent({bank_code: selectedBank, client_id: data.login})
    }

    const onError: SubmitErrorHandler<LoginBankForm> = (errors) => {
        console.log(errors)
    }

    const handleClickShowPassword = () => setShowPassword((show) => !show);

    const handleMouseDownPassword = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
    };

    const handleMouseUpPassword = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
    };

    const handleSnackbarClose = () => {
        setSnackbarOpen(false);
    };

    return (
        <>
            <PageTitle>Подключенные банки</PageTitle>
          {isLoading ? (
              <p>Загрузка...</p>
          ) : (
              banks && banks.length > 0 ? (
                  banks.map((bank) => (
                      <div key={bank.id}>{bank.name}</div>
                  ))
              ) : (
                  <p>Список банков пуст</p>
              )
          )}
            {error &&
                <p>Произошла непредвиденная ошибка при получении данных о банках, попробуйте обновить страницу...</p>}

            <Button variant="outlined" onClick={() => handleClickOpenConsent('abank')}>
                Open form dialog
            </Button>

            <Dialog open={isOpen} onClose={handleClose}>
                <DialogTitle>Выпуск согласия</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Для выпуска нового согласия - авторизуйтесь в банке
                    </DialogContentText>
                    <form onSubmit={handleSubmit(onSubmit, onError)} id="subscription-form" className={'flex flex-col'}>
                        <Controller
                            name='login'
                            control={control}
                            defaultValue={''}
                            render={({field, fieldState: {error}}) => (
                                <FormControl
                                    variant="standard"
                                    error={!!error}
                                    size={'small'}
                                >
                                    <TextField
                                        {...field}
                                        margin="dense"
                                        id="name"
                                        label="логин"
                                        {...register("login", {required: true})}
                                    />
                                </FormControl>
                            )}
                        />
                        {errors.login && <p className={'text-red-600'}>{errors.login.message}</p>}

                        <Controller
                            name='password'
                            control={control}
                            defaultValue={''}
                            render={({field, fieldState: {error}}) => (
                                <FormControl
                                    variant="standard"
                                    error={!!error}
                                    size={'small'}
                                >
                                    <TextField
                                        {...field}
                                        margin="dense"
                                        id="name"
                                        label="пароль"
                                        type={showPassword ? 'text' : 'password'}
                                        {...register("password", {required: true})}
                                        slotProps={{
                                            input: {
                                                endAdornment: <InputAdornment position="end">
                                                    <IconButton
                                                        aria-label={
                                                            showPassword ? 'hide the password' : 'display the password'
                                                        }
                                                        onClick={handleClickShowPassword}
                                                        onMouseDown={handleMouseDownPassword}
                                                        onMouseUp={handleMouseUpPassword}
                                                        edge="end"
                                                    >
                                                        {
                                                            showPassword
                                                                ?
                                                                <VisibilityOff/>
                                                                :
                                                                <Visibility/>
                                                        }
                                                    </IconButton>
                                                </InputAdornment>,
                                            },
                                        }}
                                    />
                                </FormControl>
                            )}
                        />
                        {errors.password && <p className={'text-red-600'}>{errors.password.message}</p>}
                    </form>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose} color={'error'} variant={'contained'}
                            disabled={isPending}>Отмена</Button>
                    <Button type="submit" form="subscription-form" color={'success'} variant={'contained'}
                            disabled={isPending}>
                        Выпустить согласие
                    </Button>
                </DialogActions>
            </Dialog>

            <Snackbar
                open={snackbarOpen}
                autoHideDuration={6000}
                onClose={handleSnackbarClose}
                anchorOrigin={{vertical: 'top', horizontal: 'center'}}
            >
                <Alert onClose={handleSnackbarClose} severity={`${isError ? 'error' : 'success'}`} sx={{width: '100%'}}>
                    {snackbarMessage}
                </Alert>
            </Snackbar>
        </>
    )
}
