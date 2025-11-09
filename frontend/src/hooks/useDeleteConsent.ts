import {Api} from "../api/api";
import {useMutation, useQueryClient} from "@tanstack/react-query";

export const useDeleteConsent = (setSnackbarMessage: (message: string) => void, setSnackbarOpen: (open: boolean) => void) => {
    const queryClient = useQueryClient();

    const {mutate, isPending, isError} = useMutation({
        mutationFn: async ({consent_id}: { consent_id: string }) => {
            const response = await Api.deleteConsent(consent_id);
            return response.data;
        },
        onSuccess: () => {
            setSnackbarMessage('Отзыв согласия успешно выполнен');
            setSnackbarOpen(true);
            queryClient.invalidateQueries({queryKey: ['consents']});
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