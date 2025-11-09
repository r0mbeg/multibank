import {z} from 'zod';

export const loginBankSchema = z.object({
    login: z.string().min(1, 'заполните поле "логин"'),
    password: z.string().min(1, 'заполните поле "пароль"'),
})

export const loginSchema = z.object({
    email: z.email().min(1, 'заполните поле "Эл. почта"'),
    password: z.string().min(8, 'пароль должен содержать хотя бы 8 символов'),
})

export const registerSchema = z.object({
    first_name: z.string(),
    last_name: z.string(),
    email: z.email(),
    patronymic: z.string(),
    password: z.string().min(8, 'пароль должен содержать хотя бы 8 символов'),
    birthdate: z.string(),
})