import {create} from "zustand/react";
import {persist} from "zustand/middleware";

interface User {
    userName: string;
    id: number;
}

interface AuthState {
    token: string | null;
    user: User | null;
    login: () => void;
    logout: () => void;
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            token: null,
            user: null,
            login: () => set({token: '12345678', user: {id: 1, userName: 'UserOne'}}),
            logout: () => set({token: null, user: null}),
        }),
        { name: 'auth-storage'}
    )
)