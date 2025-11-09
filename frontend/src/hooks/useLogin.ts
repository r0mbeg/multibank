import {useNavigate} from "@tanstack/react-router";
import type {LoginForm} from "../types/types";
import {Api} from "../api/api";
import {useAuthStore} from "../stores/authStore.ts";
import {useMutation} from "@tanstack/react-query";

export const useLogin = (setSnackbarMessage: (message: string) => void, setSnackbarOpen: (open: boolean) => void) => {
    const navigate = useNavigate();

    const {mutate, isPending} = useMutation({
        mutationFn: async (registerData: LoginForm) => {
            const response = await Api.login(registerData);
            return response.data;
        },
        onSuccess: (data) => {
            useAuthStore.getState().login(data.access_token, data.expires_in);
            navigate({to: '/consents'})
        },
        onError: (error) => {
            console.error('Login failed:', error);
            setSnackbarMessage('Произошла ошибка');
            setSnackbarOpen(true);
        },
        retry: false,
    });

    return {mutate, isPending};
}