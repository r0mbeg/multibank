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
            useAuthStore.getState().login(data.accsess_token, data.expires_in);
            navigate({to: '/'})
            console.log('Registration successful:', data);
        },
        onError: (error) => {
            console.error('Registration failed:', error);
            setSnackbarMessage(error.response.data.error || 'Произошла ошибка'); // Устанавливаем сообщение
            setSnackbarOpen(true);
        },
        retry: false,
    });

    return {mutate, isPending};
}