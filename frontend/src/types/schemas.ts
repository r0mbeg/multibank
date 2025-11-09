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
    first_name: z.string().min(1, 'поле не может быть пустым'),
    last_name: z.string().min(1, 'поле не может быть пустым'),
    email: z.email().min(1, 'поле не может быть пустым'),
    patronymic: z.string().min(1, 'поле не может быть пустым'),
    password: z.string().min(8, 'пароль должен содержать хотя бы 8 символов'),
    birthdate: z.string().min(1, 'поле не может быть пустым'),
})