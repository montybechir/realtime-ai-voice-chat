import { useContext } from 'react';
import { AuthContext, IAuthContext } from './AuthContext';

const useAuth = (): IAuthContext => {
    const authContext = useContext(AuthContext);

    if (!authContext) {
        throw new Error('useAuth must be used with an auth provider');
    }
    return authContext;
};

export default useAuth;
