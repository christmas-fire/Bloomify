import { useState } from 'react';
import './App.css';
import { apiFetch } from './api';

function LoginForm({ onLoginSuccess, onBack }) {
    const [form, setForm] = useState({ username: '', password: '' });
    const [error, setError] = useState('');

    const handleChange = (e) => {
        setForm({ ...form, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');
        try {
            const response = await apiFetch('/auth/sign-in', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(form),
            });
            if (response.ok) {
                const data = await response.json();
                onLoginSuccess(data.access_token);
            } else {
                const error = await response.json();
                setError(error.message || 'Ошибка входа');
            }
        } catch (err) {
            setError('Ошибка соединения с сервером');
        }
    };

    return (
        <form className="start-screen" onSubmit={handleSubmit} style={{gap: 14}}>
            <div className="start-title" style={{fontSize: '2rem', marginBottom: 0}}>
                <span role="img" aria-label="flower">🌸</span> Вход
            </div>
            <div className="form-group">
                <label htmlFor="username">Имя пользователя</label>
                <input
                    type="text"
                    id="username"
                    name="username"
                    value={form.username}
                    onChange={handleChange}
                    required
                />
            </div>
            <div className="form-group">
                <label htmlFor="password">Пароль</label>
                <input
                    type="password"
                    id="password"
                    name="password"
                    value={form.password}
                    onChange={handleChange}
                    required
                />
            </div>
            <div className="form-buttons">
                <button type="submit" className="start-btn">Войти</button>
                <button type="button" className="start-btn" style={{background:'#f0f1f3', color:'#4f5d75', border:'1.5px solid #e0e3e8'}} onClick={onBack}>Назад</button>
            </div>
            {error && <div className="success-message" style={{color:'red'}}>{error}</div>}
        </form>
    );
}

export default LoginForm; 