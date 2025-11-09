import React from 'react';
import type {Product} from "../types/types.ts";

interface IProductItemProps {
    product: Product;
}

const ProductItem: React.FC<IProductItemProps> = ({product}) => {
    return (
        <div
            className={`flex flex-col rounded-md shadow-md p-4 cursor-pointer ${product.is_recommended ? 'bg-amber-100 animate-pulse hover:opacity-100' : 'bg-slate-50'} hover:bg-slate-100`}
        >
            <p>{product.bank_name}</p>
            <p>{product.product_name}</p>
            <p>{product.description}</p>
            {product.is_recommended && (
                <p className={'mt-8 text-bold'}>
                    Мы рекомендуем!
                </p>
            )}
        </div>
    );
};

export default ProductItem;