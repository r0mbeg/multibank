import axios from "axios";

export const AxiosApiInstance = axios.create({
    baseURL: 'https://your-prod-api.com/api',
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true,
})

export class Api {
    static getConsents() {
        return AxiosApiInstance.get('/consents')
    }
}