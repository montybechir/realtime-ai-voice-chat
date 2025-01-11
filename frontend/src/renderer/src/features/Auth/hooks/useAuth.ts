import { useContext } from 'react';
import { AuthContext, AuthContextType } from '../contexts/AuthContext';

const useAuth = (): AuthContextType => {
    const authContext = useContext(AuthContext);

    if (!authContext) {
        throw new Error('useAuth must be used with an auth provider');
    }
    return authContext;
};

export default useAuth;
