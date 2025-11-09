import {createFileRoute} from '@tanstack/react-router'
import PageTitle from "../components/PageTitle.tsx";
// import { Api } from "../api/api.ts";
// import type {Product} from "../types/types.ts";
// import {useQuery} from "@tanstack/react-query";

export const Route = createFileRoute('/_auth/products')({
    component: RouteComponent,
})

const products = [
    {
        product_id: "prod-abank-deposit-001",
        productType: "deposit",
        product_name: "Накопительный депозит",
        description: "Выгодная ставка",
        interestRate: 8.5,
        minAmount: 50000,
        termMonths: 12,
        bank_id: 1,
        bank_code: "abank",
        bank_name: "Awesome Bank",
        isRecommended: true
    },
    {
        product_id: "prod-abank-card-001",
        productType: "credit_card",
        product_name: "Кредитная карта Premium",
        description: "Кэшбэк и бонусы",
        interestRate: 15.9,
        bank_id: 1,
        bank_code: "abank",
        bank_name: "Awesome Bank",
    },
    {
        product_id: "prod-abank-loan-001",
        productType: "loan",
        product_name: "Потребительский кредит",
        description: "Быстрое одобрение",
        interestRate: 12.9, minAmount: 50000,
        termMonths: 36, bank_id: 1,
        bank_code: "abank",
        bank_name: "Awesome Bank",
    },
    {
        product_id: "deposit-reliable",
        productType: "deposit",
        product_name: "Вклад \"Надежный\"",
        description: "Классический вклад с гарантированной доходностью",
        interestRate: 8.5,
        minAmount: 10000,
        maxAmount: 10000000,
        termMonths: 12,
        bank_id: 1,
        bank_code: "abank",
        bank_name: "Awesome Bank",
    },
    {
        product_id: "loan-consumer",
        productType: "loan",
        product_name: "Потребительский кредит",
        description: "Кредит на любые цели",
        interestRate: 12.9,
        minAmount: 50000,
        maxAmount: 3000000,
        termMonths: 60,
        bank_id: 1,
        bank_code: "abank",
        bank_name: "Awesome Bank",
    },
]

function RouteComponent() {
    // const { data: products, isLoading, error } = useQuery<Product[]>({
    //     queryKey: ['products'],
    //     queryFn: async (): Promise<Product[]> => {
    //         const response = await Api.getProducts();
    //         return response.data;
    //     },
    //     retry: false,
    // })
    const isLoading = false;

    return (
        <>
            <PageTitle>Витрина предложений</PageTitle>
            {isLoading ? (
                <span>Загрузка...</span>
            ) : (
                <div className="grid grid-cols-4 gap-8 mt-4">
                    {products && products.length > 0 ? (
                        products.map((product) => (
                            <div
                                key={product.product_id}
                                className={`flex flex-col rounded-md shadow-md p-4 cursor-pointer ${product.isRecommended ? 'bg-amber-50' : 'bg-slate-50'} hover:bg-slate-100`}
                            >
                                <p>{product.bank_name}</p>
                                <p>{product.product_name}</p>
                                <p>{product.description}</p>
                            </div>
                        ))
                    ) : (
                        <p>Нет доступных предложений :(</p>
                    )
                    }
                </div>
            )}
        </>
    )
}
