export type LoginForm = {
    login: string,
    password: string,
}

export type Consents = {
    bankName: string,
    consentStatus: 'active' | 'pending' | null,
}