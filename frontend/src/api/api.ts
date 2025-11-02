import axios from "axios";
import type {LoginForm, RegisterForm} from "../types/types.ts";

export const AxiosApiInstance = axios.create({
    baseURL: 'http://localhost:8080/',
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true,
})

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
}