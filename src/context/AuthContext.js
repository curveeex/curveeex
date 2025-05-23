import { createContext, useState, useEffect } from 'react';
import * as authApi from '../api/auth';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const checkAuth = async () => {
            try {
                const response = await authApi.getProfile();
                setUser(response.data);
            } catch (error) {
                console.error('Auth check failed', error);
            } finally {
                setLoading(false);
            }
        };
        checkAuth();
    }, []);

    const register = async (email, password) => {
        await authApi.register(email, password);
    };

    const login = async (email, password) => {
        const user = await authApi.login(email, password);
        setUser(user);
    };

    const logout = async () => {
        await authApi.logout();
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ user, loading, register, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};