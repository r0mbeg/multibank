import React from 'react';
import type {Consents} from "../types/types.ts";
import {Button} from "@mui/material";

interface IConsentItemProps {
    consent: Consents
    handleDeleteConsent: (consent_id: string) => void
}

const ConsentItem: React.FC<IConsentItemProps> = ({consent, handleDeleteConsent}) => {
    return (
        <div className={'flex justify-between border mb-4 shadow-md rounded-md p-4 mt-4'}>
            <div>
                <p className={'text-xl'}>{consent.bank_code} ({consent.client_id})</p>
                {consent.status === 'Authorized' &&
                    <p className={'text-green-500'}>согласие активно</p>}
                {consent.status === 'AwaitingAuthorization' &&
                    <p className={'text-orange-500'}>согласие в обработке</p>}
            </div>

            {consent.status === 'Authorized' && (
                <>

                    <Button variant={'contained'} color={'error'}
                            onClick={() => handleDeleteConsent(consent.consent_id)}>отозвать
                        согласие</Button>
                </>
            )}
        </div>
    );
};

export default ConsentItem;