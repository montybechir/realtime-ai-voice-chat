import useAuth from '@renderer/features/Auth/hooks/useAuth';
import { Navigate } from 'react-router-dom';

interface ProtectRouteProps {
    children: JSX.Element;
}

const ProtectedRoute: React.FC<ProtectRouteProps> = ({ children }) => {
    const { isAuthenticated } = useAuth();
    if (!isAuthenticated) {
        return <Navigate to="/login" replace />;
    }
    return children;
};

export default ProtectedRoute;
