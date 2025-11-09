import {Link} from "@tanstack/react-router";
import {Button} from "@mui/material";
import {AccountBalance, Storefront, Description, Payments} from "@mui/icons-material";
import {useAuthStore} from "../stores/authStore.ts";

const links = [
    {to: '/accounts', icon: <Payments/>, text: 'Счета'},
    {to: '/products', icon: <Storefront/>, text: 'Витрина'},
    {to: '/consents', icon: <Description/>, text: 'Согласия'},
    {to: '/banks', icon: <AccountBalance/>, text: 'Банки'},
];

const Sidebar = () => {
    const {user} = useAuthStore();

    return (
        <nav className="w-64 bg-white rounded-xl shadow-md p-4">
            <p className={'mb-4'}>{user?.lastName} {user?.firstName.slice(0, 1).toUpperCase()}. {user?.patronic.slice(0, 1).toUpperCase()}.</p>
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