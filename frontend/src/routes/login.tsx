import {createFileRoute, createLink, useNavigate, useSearch} from "@tanstack/react-router";
import {Controller, type SubmitErrorHandler, type SubmitHandler, useForm} from "react-hook-form";
import {useAuthStore} from "../stores/authStore.ts";
import React, {useState} from "react";
import {
    Button,
    FormControl,
    IconButton,
    InputAdornment,
    InputLabel,
    OutlinedInput,
    Link as MUILink,
} from "@mui/material";
import {VisibilityOff, Visibility, WarningAmber} from "@mui/icons-material";
import ErrorText from "../components/ErrorText.tsx";
import type {LoginForm} from "../types/types.ts";

export const Route = createFileRoute('/login')({
    component: LoginPage,
    validateSearch: (search: Record<string, unknown>) => ({
        redirect: (typeof search.redirect === 'string' ? search.redirect : undefined) as string | undefined,
    })
})

const CustomLink = createLink(MUILink);

function LoginPage() {
    const [showPassword, setShowPassword] = useState(false);
    const navigate = useNavigate();
    const { login } = useAuthStore();
    const { redirect } = useSearch({ from: '/login' });
    const {
        control,
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<LoginForm>()

    const onSubmit: SubmitHandler<LoginForm> = (data) => {
        login();
        navigate({to: redirect || '/'})
        console.log(data)
    }

    const onError: SubmitErrorHandler<LoginForm> = (errors) => {
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
            <h1 className={'text-6xl mb-8'}>Multibank APP</h1>
            <form className={'flex flex-col max-w-80 gap-y-4 p-4 rounded-md bg-white shadow-md'} onSubmit={handleSubmit(onSubmit, onError)}>
                <Controller
                    name='login'
                    control={control}
                    defaultValue={''}
                    render={({ field, fieldState: { error } }) => (
                        <FormControl
                            variant="outlined"
                            error={!!error}
                            size={'small'}
                        >
                            <InputLabel htmlFor="outlined-adornment-login">Логин</InputLabel>
                            <OutlinedInput
                                id={"outlined-adornment-login"}
                                {...field}
                                label={'логин'}
                                {...register("login", { required: true })}
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

                {(errors.login || errors.password) && <ErrorText><WarningAmber /> Заполните все поля</ErrorText>}

                <Button variant="contained" type={'submit'}>Войти</Button>

                <CustomLink to={'/register'} underline={'hover'} sx={{textAlign: 'center'}}>Зарегистрироваться</CustomLink>
            </form>
        </div>
    );
}

export default LoginPage;