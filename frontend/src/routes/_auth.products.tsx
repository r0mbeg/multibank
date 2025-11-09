import {createFileRoute} from '@tanstack/react-router'
import PageTitle from "../components/PageTitle.tsx";
import {Api} from "../api/api.ts";
import type {Product} from "../types/types.ts";
import {useQuery} from "@tanstack/react-query";
import ProductItem from "../components/ProductItem.tsx";
import {useMemo} from "react";

export const Route = createFileRoute('/_auth/products')({
    component: RouteComponent,
})

function RouteComponent() {
    const {data: products, isLoading, error} = useQuery<Product[]>({
        queryKey: ['products'],
        queryFn: async (): Promise<Product[]> => {
            const response = await Api.getProducts();
            return response.data;
        },
        retry: false,
    })

    const sortedProducts = useMemo(() => {
        return products?.sort((a, b) => {
            if (a.is_recommended && !b.is_recommended) return -1;
            if (!a.is_recommended && b.is_recommended) return 1;
            return 0;
        }) || [];
    }, [products]);

    return (
        <>
            <PageTitle>Витрина предложений</PageTitle>

            {isLoading ? (
                <span>Загрузка...</span>
            ) : (
                <div className="grid grid-cols-3 gap-4 mt-4">
                    {sortedProducts && sortedProducts.length > 0 ? (
                        sortedProducts.map((product) => (
                            <ProductItem key={product.product_id} product={product}/>
                        ))
                    ) : (
                        <p>Нет доступных предложений :(</p>
                    )
                    }
                </div>
            )}

            {error && <p>Произошла непредвиденная ошибка при получении данных о продуктах, попробуйте обновить
                страницу...</p>}
        </>
    )
}
