import { useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';

export const Navbar = () => {
    const { user, logout } = useContext(AuthContext);
    const navigate = useNavigate();

    const handleLogout = async () => {
        await logout();
        navigate('/login');
    };

    return (
        <nav>
            {user ? (
                <>
                    <Link to="/profile">Профиль</Link>
                    <button onClick={handleLogout}>Выйти</button>
                </>
            ) : (
                <>
                    <Link to="/login">Вход</Link>
                    <Link to="/register">Регистрация</Link>
                </>
            )}
        </nav>
    );
};