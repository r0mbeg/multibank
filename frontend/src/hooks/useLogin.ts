import { useNavigate } from "@tanstack/react-router";
import type { LoginForm } from "../types/types";
import { Api } from "../api/api";
import {useAuthStore} from "../stores/authStore.ts";
import {useMutation} from "@tanstack/react-query";
import type {AxiosError} from "axios";

export const useLogin = (setSnackbarMessage: (message: string) => void, setSnackbarOpen: (open: boolean) => void) => {
    const navigate = useNavigate();

    const {mutate, isPending} = useMutation({
        mutationFn: async (registerData: LoginForm) => {
            const response = await Api.login(registerData);
            return response.data;
        },
        onSuccess: (data) => {
            useAuthStore.getState().login(data.access_token, data.expires_in);
            navigate({to: '/'})
            console.log('Login successful:', data);
        },
        onError: (error: AxiosError) => {
            console.error('Login failed:', error);
            setSnackbarMessage(error.response.data.error || 'Произошла ошибка'); // Устанавливаем сообщение
            setSnackbarOpen(true);
        },
        retry: false,
    });

    return {mutate, isPending};
}