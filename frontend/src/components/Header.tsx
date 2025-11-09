import {Link, useNavigate} from "@tanstack/react-router";
import LogoutIcon from '@mui/icons-material/Logout';
import {useAuthStore} from "../stores/authStore.ts";
import {Button} from "@mui/material";

const Header = () => {
    const {logout} = useAuthStore.getState();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout()
        navigate({
            to: '/login',
        })
    }

    return (
        <header className="shadow-md bg-white rounded-md flex items-center justify-between p-4">
            <Link to={'/'}>MultiBankAPP</Link>
            <Button
                startIcon={<LogoutIcon/>}
                onClick={handleLogout}
            >
                Выйти
            </Button>
        </header>
    );
};

export default Header;