import React from "react";

type PageTitleProps = {
    children: React.ReactNode
}

const PageTitle: React.FC<PageTitleProps> = ({children}) => {
    return (
        <h1 className={'text-4xl'}>{children}</h1>
    );
};

export default PageTitle;