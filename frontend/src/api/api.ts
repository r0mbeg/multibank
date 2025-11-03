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

// Interceptor для запросов: динамически добавляем токен из store
AxiosApiInstance.interceptors.request.use(
    (config) => {
        const token = useAuthStore.getState().token; // Берём токен на момент запроса
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);
// Опционально: Interceptor для ответов — обрабатываем 401 (токен истёк)
AxiosApiInstance.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // Токен истёк или недействителен — разлогиниваем
            useAuthStore.getState().logout();
            // Опционально: перенаправьте на /login
            // window.location.href = '/login';
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
}