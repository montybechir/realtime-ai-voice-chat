import { createContext } from 'react';
import { User } from '../profile/types';

export const AuthContext = createContext<IAuthContext | undefined>(undefined);
export interface IAuthContext {
	isAuthenticated: boolean;
	user: User | undefined;
	loading: boolean;
	error: Error | undefined;
	login: () => Promise<void>;
	register: () => Promise<void>;
	logout: () => Promise<void>;
}
