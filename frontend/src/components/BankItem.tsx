import React from 'react';
import {Button} from "@mui/material";
import type {Bank} from "../types/types.ts";

interface IBankItemProps {
    bank: Bank
    handleClickOpenConsent: (bankCode: string) => void
}

const BankItem: React.FC<IBankItemProps> = ({bank, handleClickOpenConsent}) => {
    return (
        <div className={'rounded-md shadow-md p-4 flex justify-between items-center'}>
            {bank.name}
            <Button variant={'contained'} color={'success'}
                    onClick={() => handleClickOpenConsent(bank.code)}>Выпустить согласие</Button>
        </div>
    );
}

export default BankItem;