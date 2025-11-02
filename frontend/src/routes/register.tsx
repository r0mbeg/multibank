import {createFileRoute, createLink} from '@tanstack/react-router'
import {Controller, type SubmitErrorHandler, type SubmitHandler, useForm} from "react-hook-form";
import {
    Button,
    FormControl,
    IconButton,
    InputAdornment,
    InputLabel,
    Link as MUILink,
    OutlinedInput,
} from "@mui/material";
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import {Visibility, VisibilityOff, WarningAmber} from "@mui/icons-material";
import React, {useState} from "react";
import ErrorText from "../components/ErrorText.tsx";
import type { RegisterForm} from "../types/types.ts";
import {useMutation, useQuery} from "@tanstack/react-query";
import {Api} from "../api/api.ts";
import dayjs, { Dayjs } from 'dayjs';
import {LocalizationProvider} from "@mui/x-date-pickers";
import {AdapterDayjs} from "@mui/x-date-pickers/AdapterDayjs";
import {ruRU} from "@mui/x-date-pickers/locales";

export const Route = createFileRoute('/register')({
  component: RouteComponent,
})

const CustomLink = createLink(MUILink);

function RouteComponent() {
    const [showPassword, setShowPassword] = useState(false);
    const {
        control,
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<RegisterForm>()

    const registrationMutation = useMutation({
        mutationFn: async (registerData: RegisterForm) => {
            const response = await Api.registration(registerData);
            return response.data;
        },
        onSuccess: (data) => {
            console.log('Registration successful:', data);
        },
        onError: (error) => {
            console.error('Registration failed:', error);
        },
        retry: false,
    });
    const onSubmit: SubmitHandler<RegisterForm> = (registerData) => {
        console.log(registerData);
        registrationMutation.mutate(registerData);
    }

    const onError: SubmitErrorHandler<RegisterForm> = (errors) => {
        console.log(errors)
    }

    const handleClickShowPassword = () => setShowPassword((show) => !show);

    const handleMouseDownPassword = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
    };

    const handleMouseUpPassword = (event: React.MouseEvent<HTMLButtonElement>) => {
        event.preventDefault();
    };

    return (
        <div className={'flex flex-col justify-center items-center h-full'}>
            <h1 className={'text-6xl mb-8'}>Регистрация в Multibank APP</h1>
            <form className={'flex flex-col max-w-80 gap-y-4 p-4 rounded-md bg-white shadow-md'} onSubmit={handleSubmit(onSubmit, onError)}>
                <Controller
                    name='first_name'
                    control={control}
                    defaultValue={''}
                    render={({ field, fieldState: { error } }) => (
                        <FormControl
                            variant="outlined"
                            error={!!error}
                            size={'small'}
                        >
                            <InputLabel htmlFor="outlined-adornment-first-name">Имя</InputLabel>
                            <OutlinedInput
                                id={"outlined-adornment-first-name"}
                                {...field}
                                label={'имя'}
                                {...register("first_name", { required: true })}
                            />
                        </FormControl>
                    )}
                />

                <Controller
                    name='last_name'
                    control={control}
                    defaultValue={''}
                    render={({ field, fieldState: { error } }) => (
                        <FormControl
                            variant="outlined"
                            error={!!error}
                            size={'small'}
                        >
                            <InputLabel htmlFor="outlined-adornment-last-name">Фамилия</InputLabel>
                            <OutlinedInput
                                id={"outlined-adornment-last-name"}
                                {...field}
                                label={'фамилия'}
                                {...register("last_name", { required: true })}
                            />
                        </FormControl>
                    )}
                />

                <Controller
                    name='patronymic'
                    control={control}
                    defaultValue={''}
                    render={({ field, fieldState: { error } }) => (
                        <FormControl
                            variant="outlined"
                            error={!!error}
                            size={'small'}
                        >
                            <InputLabel htmlFor="outlined-adornment-patronymic">Отчество</InputLabel>
                            <OutlinedInput
                                id={"outlined-adornment-patronymic"}
                                {...field}
                                label={'отчество'}
                                {...register("patronymic", { required: true })}
                            />
                        </FormControl>
                    )}
                />

                <Controller
                    name='birthdate'
                    control={control}
                    defaultValue={''}
                    rules={{
                        required: 'Дата рождения обязательна',
                        validate: (value) => {
                            if (value && dayjs(value).isAfter(dayjs())) {
                                return 'Дата рождения не может быть в будущем';
                            }
                            if (value && dayjs().diff(dayjs(value), 'years') < 18) {
                                return 'Вы должны быть старше 18 лет';
                            }
                            return true;
                        }
                    }}
                    render={({ field, fieldState: { error } }) => (
                        <LocalizationProvider dateAdapter={AdapterDayjs} localeText={ruRU.components.MuiLocalizationProvider.defaultProps.localeText}>
                            <DatePicker
                                label="Дата рождения"
                                value={field.value ? dayjs(field.value) : null}
                                onChange={(newValue: Dayjs | null) => {
                                    field.onChange(newValue ? newValue.format('YYYY-MM-DD') : '');
                                }}
                                disableFuture
                                slotProps={{
                                    textField: {
                                        variant: 'outlined',
                                        size: 'small',
                                        error: !!error,
                                        helperText: error?.message,
                                        fullWidth: true,
                                    },
                                }}
                            />
                        </LocalizationProvider>
                    )}
                />

                <Controller
                    name='email'
                    control={control}
                    defaultValue={''}
                    render={({ field, fieldState: { error } }) => (
                        <FormControl
                            variant="outlined"
                            error={!!error}
                            size={'small'}
                        >
                            <InputLabel htmlFor="outlined-adornment-email">Эл. почта</InputLabel>
                            <OutlinedInput
                                id={"outlined-adornment-email"}
                                {...field}
                                label={'эл. почта'}
                                {...register("email", { required: true })}
                            />
                        </FormControl>
                    )}
                />

                <Controller
                    name={'password'}
                    control={control}
                    defaultValue={''}
                    render={({ field, fieldState: { error } }) => (
                        <FormControl
                            variant="outlined"
                            error={!!error}
                            size={'small'}
                        >
                            <InputLabel htmlFor="outlined-adornment-password">Пароль</InputLabel>
                            <OutlinedInput
                                id="outlined-adornment-password"
                                {...field}
                                type={showPassword ? 'text' : 'password'}
                                label={'пароль'}
                                {...register("password", { required: true })}
                                endAdornment={
                                    <InputAdornment position="end">
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
                                                    <VisibilityOff color={!!error ? 'error' : 'inherit'} />
                                                    :
                                                    <Visibility color={!!error ? 'error' : 'inherit'} />
                                            }
                                        </IconButton>
                                    </InputAdornment>
                                }
                            />
                        </FormControl>
                    )}
                />

                {/*{(errors.login || errors.password) && <ErrorText><WarningAmber /> Заполните все поля</ErrorText>}*/}

                <Button variant="contained" type={'submit'}>Зарегистрироваться</Button>

                <p style={{textAlign: 'center'}}>Есть аккаунт? <CustomLink to={'/login'} underline={'hover'} search={{redirect: undefined}}>Войти</CustomLink></p>
            </form>
        </div>
    );
}
