import {useMutation} from "@tanstack/react-query";
import type {RegisterForm} from "../types/types.ts";
import {Api} from "../api/api.ts";
import {useAuthStore} from "../stores/authStore.ts";
import {useNavigate} from "@tanstack/react-router";

export const useRegistration = (setSnackbarMessage: (message: string) => void, setSnackbarOpen: (open: boolean) => void) => {
    const navigate = useNavigate();

    const {mutate, isPending} = useMutation({
        mutationFn: async (registerData: RegisterForm) => {
            const response = await Api.registration(registerData);
            return response.data;
        },
        onSuccess: (data) => {
            useAuthStore.getState().login(data.access_token, data.expires_in);
            navigate({to: '/'})
        },
        onError: (error) => {
            console.error('Registration failed:', error);
            setSnackbarMessage('Произошла ошибка');
            setSnackbarOpen(true);
        },
        retry: false,
    });

    return {mutate, isPending};
}