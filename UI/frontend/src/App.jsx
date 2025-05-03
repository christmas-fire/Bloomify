import './App.css';
import { useState } from 'react';
import RegistrationForm from './RegistrationForm';
import LoginForm from './LoginForm';
import MainScreen from './MainScreen';

function App() {
    const [screen, setScreen] = useState('start'); // start | register | login | main
    const [accessToken, setAccessToken] = useState(null);

    const handleRegistrationSuccess = () => {
        setScreen('login');
    };

    const handleLoginSuccess = (token) => {
        setAccessToken(token);
        setScreen('main');
    };

    const handleLogout = () => {
        setAccessToken(null);
        setScreen('start');
    };

    if (screen === 'main' && accessToken) {
        return <MainScreen accessToken={accessToken} onLogout={handleLogout} />;
    }

    return (
        <div id="App">
            {screen === 'start' && (
                <div className="start-screen">
                    <div className="start-title">
                        <span role="img" aria-label="flower">🌸</span> Bloomify
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
        </div>
    );
}

export default App;
