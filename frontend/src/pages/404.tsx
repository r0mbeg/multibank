import ErrorIcon from '@mui/icons-material/Error';
import { Link } from '@tanstack/react-router'

const NotFoundPage = () => {
    return (
        <div style={{height:'100%', display:'flex', justifyContent:'center', alignItems: 'center'}}>
            <h1>Страница не найдена <ErrorIcon/></h1>
            <Link to={'/'}>Домой</Link>
        </div>
    );
};

export default NotFoundPage;