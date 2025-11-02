import { Link } from "@tanstack/react-router";
import {Button} from "@mui/material";
import { AccountBalance, ReceiptLong, Storefront, Description } from "@mui/icons-material";

const links = [
    { to: '/accounts', icon: <AccountBalance />, text: 'Счета' },
    { to: '/transactions', icon: <ReceiptLong />, text: 'Переводы' },
    { to: '/showcase', icon: <Storefront />, text: 'Витрина' },
    { to: '/consents', icon: <Description />, text: 'Согласия' },
];

const Sidebar = () => {
    return (
        <nav className="w-64 bg-white rounded-xl shadow-md p-4">
            {links.map((link, idx) => (
                <Link to={link.to} key={idx}>
                    <Button
                        variant={'text'}
                        startIcon={link.icon}
                        sx={{
                            width: '100%',
                            justifyContent: 'flex-start',
                            marginBottom: idx < links.length - 1 ? 2 : 0,
                        }}
                    >
                        {link.text}
                    </Button>
                </Link>
            ))}
        </nav>
    );
};

export default Sidebar;