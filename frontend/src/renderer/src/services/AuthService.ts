import { User } from '@renderer/features/profile/types';

interface IAuthService {
    register(): Promise<User>;
    login(): Promise<User>;
    logout(): Promise<void>;
    getCurrentUser(): Promise<User>;
}

export const AuthService: IAuthService = {
    async register(): Promise<User> {},
    async login(): Promise<User> {},
    async logout(): Promise<void> {},
    async getCurrentUser(): Promise<User> {
        return null;
    }
};
