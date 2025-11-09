import {loginBankSchema, loginSchema, registerSchema} from "./schemas.ts";
import {z} from "zod";

export type RegisterForm = z.infer<typeof registerSchema>

export type LoginForm = z.infer<typeof loginSchema>

export type LoginBankForm = z.infer<typeof loginBankSchema>

export type Consents = {
    bankName: string,
    consentStatus: 'active' | 'pending' | null,
}

export type Product = {
    product_id: string,
    product_name: string,
    description: string,
    bank_name: string,
}

export type Bank = {
    id: number,
    name: string,
    code: string,
    is_enabled: boolean,
    authorized: boolean,
    token_expires: string
}