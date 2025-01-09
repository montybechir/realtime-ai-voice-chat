import { AuthContext, IAuthContext } from './AuthContext';
import { ReactNode, useEffect, useLayoutEffect, useState } from 'react';
import { User } from '../profile/types';
import { AuthService } from '@renderer/services/AuthService';

interface AuthProviderProps {
    children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps): ReactNode {
    const [user, setUser] = useState<User | undefined>(undefined);
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<Error | undefined>();

    const login = async (): Promise<void> => {
        try {
            const userData = await AuthService.login();
            setUser(userData);
            setLoading(false);
        } catch (e) {
            setError(e instanceof Error ? error : new Error('Login failed'));
        }
    };
    const register = async (): Promise<void> => {
        try {
            const userData = await AuthService.login();
            setUser(userData);
            setLoading(false);
        } catch (e) {
            setError(e instanceof Error ? error : new Error('Login failed'));
        }
    };

    const logout = async (): Promise<void> => {
        try {
            await AuthService.logout();
            setUser(undefined);
            setLoading(false);
        } catch (e) {
            setError(e instanceof Error ? error : new Error('Login failed'));
        }
    };

    const contextValue: IAuthContext = {
        isAuthenticated,
        user,
        loading,
        error,
        login,
        register,
        logout
    };

    useEffect(() => {
        if (user) {
            setIsAuthenticated(true);
        } else {
            setIsAuthenticated(false);
        }
    }, [user]);

    useLayoutEffect(() => {
        const checkAuth = async () => {
            try {
                const userData = await AuthService.getCurrentUser();
                setUser(userData);
            } catch (e) {
                setError(e instanceof Error ? e : new Error('Auth check failed'));
            } finally {
                setLoading(false);
            }
        };

        checkAuth();
    }, []);

    return <AuthContext.Provider value={contextValue}>{!loading && children}</AuthContext.Provider>;
}
