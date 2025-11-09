import {useNavigate} from "@tanstack/react-router";
import {Api} from "../api/api";
import {useMutation} from "@tanstack/react-query";

export const useNewConsent = (setSnackbarMessage: (message: string) => void, setSnackbarOpen: (open: boolean) => void) => {
    const navigate = useNavigate();

    const {mutate, isPending, isError} = useMutation({
        mutationFn: async ({bank_code, client_id}: { bank_code: string, client_id: string }) => {
            const response = await Api.newConsent(bank_code, client_id);
            return response.data;
        },
        onSuccess: () => {
            setSnackbarMessage('Успешная авторизация, согласие обрабатывается... вы будете перенаправленны на страницу согласий через 6 секунд...');
            setSnackbarOpen(true);
            setTimeout(() => {
                navigate({to: '/consents'})
            }, 6000);
        },
        onError: (error) => {
            console.error('Ошибка при подключении банка:', error);
            setSnackbarMessage('Произошла ошибка при подключении банка');
            setSnackbarOpen(true);
        },
        retry: false,
    });

    return {mutate, isPending, isError};
}