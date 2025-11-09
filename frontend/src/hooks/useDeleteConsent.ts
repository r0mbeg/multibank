import {Api} from "../api/api";
import {useMutation} from "@tanstack/react-query";

export const useDeleteConsent = (setSnackbarMessage: (message: string) => void, setSnackbarOpen: (open: boolean) => void) => {
    const {mutate, isPending, isError} = useMutation({
        mutationFn: async ({bank_code, consent_id}: { bank_code: string, consent_id: string }) => {
            const response = await Api.deleteConsent(bank_code, consent_id);
            return response.data;
        },
        onSuccess: () => {
            setSnackbarMessage('Отзыв согласия успешно выполнен');
            setSnackbarOpen(true);
        },
        onError: (error) => {
            console.error('Ошибка при отзыве согласия:', error);
            setSnackbarMessage('Произошла ошибка при отзыве согласия');
            setSnackbarOpen(true);
        },
        retry: false,
    });

    return {mutate, isPending, isError};
}