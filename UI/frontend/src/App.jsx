import './App.css';
import { useState, useEffect } from 'react';
import RegistrationForm from './RegistrationForm';
import LoginForm from './LoginForm';
import MainScreen from './MainScreen';
import { getRefreshToken, setAccessToken, apiFetch, setAccessTokenGetter } from './api';

function App() {
    const [screen, setScreen] = useState('start'); // start | register | login | main
    const [accessToken, setAccessTokenState] = useState(null);
    const [refreshChecked, setRefreshChecked] = useState(false);

    // При монтировании: если есть refresh token, пробуем refresh
    useEffect(() => {
        setAccessTokenGetter(() => accessToken);
    }, [accessToken]);

    useEffect(() => {
        const tryRefresh = async () => {
            const refreshToken = getRefreshToken();
            if (refreshToken) {
                try {
                    const resp = await apiFetch('/auth/refresh', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ refresh_token: refreshToken })
                    });
                    if (resp.ok) {
                        const tokens = await resp.json();
                        setAccessToken(tokens.access_token);
                        setAccessTokenState(tokens.access_token);
                        setScreen('main');
                    } else {
                        setScreen('login');
                    }
                } catch {
                    setScreen('login');
                }
            } else {
                setScreen('start');
            }
            setRefreshChecked(true);
        };
        tryRefresh();
    }, []);

    const handleRegistrationSuccess = () => {
        setScreen('login');
    };

    const handleLoginSuccess = (token) => {
        setAccessTokenState(token);
        setScreen('main');
    };

    const handleLogout = () => {
        setAccessTokenState(null);
        setAccessToken(null);
        setScreen('start');
    };

    // Обновление accessToken при refresh
    const handleTokenRefresh = (newToken) => {
        setAccessTokenState(newToken);
    };

    if (!refreshChecked) return null;

    return (
        <div id="App">
            {screen === 'start' && (
                <div className="start-screen">
                    <div className="start-title">
                        <span role="img" aria-label="flower">🌸 Bloomify</span>
                    </div>
                    <div className="start-subtitle">
                        Твой онлайн-магазин цветов<br/>Красота и забота — в каждом букете
                    </div>
                    <button className="start-btn" onClick={() => setScreen('login')}>Войти</button>
                    <button className="start-btn" onClick={() => setScreen('register')}>Зарегистрироваться</button>
                </div>
            )}
            {screen === 'register' && (
                <RegistrationForm onSuccess={handleRegistrationSuccess} onBack={() => setScreen('start')} />
            )}
            {screen === 'login' && (
                <LoginForm onLoginSuccess={handleLoginSuccess} onBack={() => setScreen('start')} />
            )}
            {screen === 'main' && accessToken && (
                <MainScreen accessToken={accessToken} onLogout={handleLogout} onTokenRefresh={handleTokenRefresh} />
            )}
        </div>
    );
}

export default App;
