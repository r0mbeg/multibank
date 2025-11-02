export type RegisterForm = {
    first_name: string,
    last_name: string,
    email: string,
    patronymic: string,
    password: string,
    birthdate: string,
}

export type LoginForm = {
    email: string,
    password: string,
}

export type Consents = {
    bankName: string,
    consentStatus: 'active' | 'pending' | null,
}