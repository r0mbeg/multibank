import {loginBankSchema, loginSchema, registerSchema} from "./schemas.ts";
import {z} from "zod";

export type RegisterForm = z.infer<typeof registerSchema>

export type LoginForm = z.infer<typeof loginSchema>

export type LoginBankForm = z.infer<typeof loginBankSchema>

export type Consents = {
    bank_code: string,
    status: 'Authorized' | 'AwaitingAuthorization' | null,
    consent_id: string,
    client_id: string,
}

export type Accounts = {
    account_id: string,
    account_sub_type: string,
    amount: string,
    bank_code: string,
    client_id: string,
    currency: string,
    nickname: string,
    opening_date: string,
    status: string,
}

export type Product = {
    product_id: string,
    product_name: string,
    description: string,
    bank_name: string,
    is_recommended?: boolean,
}

export type Bank = {
    id: number,
    name: string,
    code: string,
    is_enabled: boolean,
    authorized: boolean,
    token_expires: string
}