import {type ReactNode} from 'react';

interface ErrorTextProps {
    children: ReactNode;
}

const ErrorText = ({children}: ErrorTextProps) => {
    return (
        <div style={{display: 'inline-flex', alignItems: 'end', margin: '0 auto', color: 'red'}}>
            {children}
        </div>
    );
};

export default ErrorText;