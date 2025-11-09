import {create} from "zustand/react";
import {devtools, persist} from "zustand/middleware";

interface User {
    firstName: string;
    lastName: string;
    patronic: string;
    email: string;
}

interface AuthState {
    token: string | null;
    expiresAt: number | null;
    user: User | null;
    isAuthenticated: boolean;
    setUser: (user: Me) => void;
    login: (token: string, expiresIn: number, user?: User) => void;
    logout: () => void;
    checkTokenValidity: () => boolean;
}

interface Me {
    birthdate: string
    created_at: string
    email: string
    first_name: string
    id: number
    last_name: string
    patronymic: string
    updated_at: string
}

export const useAuthStore = create<AuthState>()(devtools(
    persist(
        (set, get) => ({
            token: null,
            expiresAt: null,
            user: null,
            isAuthenticated: false,
            setUser: (data) => set({
                user: {
                    firstName: data.first_name,
                    lastName: data.last_name,
                    patronic: data.patronymic,
                    email: data.email,
                }
            }),
            login: (token: string, expiresIn: number, user?: User) => {
                const expiresAt = Date.now() + expiresIn * 1000; // Преобразуем в timestamp
                set({
                    token,
                    expiresAt,
                    user: user || null,
                    isAuthenticated: true,
                });
            },
            logout: () => set({
                token: null,
                expiresAt: null,
                user: null,
                isAuthenticated: false,
            }),
            checkTokenValidity: () => {
                const {expiresAt} = get();
                return expiresAt ? Date.now() < expiresAt : false;
            },
        }),
        {name: 'auth-storage'} // Сохраняется в localStorage
    ))
);