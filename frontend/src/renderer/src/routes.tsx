import { Routes as RouterRoutes, Route } from 'react-router-dom';
import ProtectedRoute from './components/ProtectedRoute';
import NotFound from './pages/NotFound';
import Home from './pages/Home';
import Login from './pages/Login';

const Routes: React.FC = () => {
    return (
        <RouterRoutes>
            <Route
                path="/"
                element={
                    <ProtectedRoute>
                        <Home />
                    </ProtectedRoute>
                }
            />
            <Route path="/login" element={<Login />} />
            <Route path="*" element={<NotFound />} />
        </RouterRoutes>
    );
};

export default Routes;
