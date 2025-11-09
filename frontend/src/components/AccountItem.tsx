import type { Accounts } from "../types/types.ts";
import React from "react";

interface IAccountItemProps {
    account: Accounts;
};

const AccountItem: React.FC<IAccountItemProps> = ({ account }) => {
    return (
        <div className={'rounded-md shadow-md p-4'}>
            <div>
                <p className={'font-bold'}>{account.bank_code} <span
                    className={'font-normal'}>(номер счета: {account.account_id})</span>
                </p>
                <p>{account.nickname} пользователя {account.client_id}</p>
            </div>
            <p>Баланс: <span className={'font-bold'}>{account.amount}</span> {account.currency}
            </p>
        </div>
    );
};

export default AccountItem;
