import axios from "axios";
import type {LoginForm, RegisterForm} from "../types/types.ts";
import {useAuthStore} from "../stores/authStore.ts";

export const AxiosApiInstance = axios.create({
    baseURL: 'http://localhost:8080/',
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true,
})

AxiosApiInstance.interceptors.request.use(
    (config) => {
        const token = useAuthStore.getState().token;
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

AxiosApiInstance.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            useAuthStore.getState().logout();
        }
        return Promise.reject(error);
    }
);

export class Api {
    static getConsents() {
        return AxiosApiInstance.get('/consents')
    }

    static login({email, password}: LoginForm) {
        return AxiosApiInstance.post('/auth/login', {
            email,
            password,
        })
    }

    static registration({first_name, last_name, birthdate, email, patronymic, password}: RegisterForm) {
        return AxiosApiInstance.post('/auth/register', {
            first_name,
            last_name,
            birthdate,
            email,
            patronymic,
            password,
        })
    }

    static getMe() {
        return AxiosApiInstance.get('/me')
    }

    static getProducts() {
        return AxiosApiInstance.get('/products')
    }

    static getBanks() {
        return AxiosApiInstance.get('/banks')
    }

    static newConsent(bank_code: string, client_id: string) {
        return AxiosApiInstance.post('/consent', {
            bank_code,
            client_id,
        })
    }

    static deleteConsent(bank_code: string, client_id: string) {
        return AxiosApiInstance.delete(`/consent`, {
            data: {
                bank_code,
                client_id,
            }
        })
    }
}