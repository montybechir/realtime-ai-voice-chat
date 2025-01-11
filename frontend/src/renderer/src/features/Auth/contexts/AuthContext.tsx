import { createContext } from 'react';
import { User } from '../../Profile/types';

export const AuthContext = createContext<AuthContextType | undefined>(undefined);
export interface AuthContextType {
    isAuthenticated: boolean;
    user: User | undefined;
    loading: boolean;
    error: Error | undefined;
    login: () => Promise<void>;
    register: () => Promise<void>;
    logout: () => Promise<void>;
}
