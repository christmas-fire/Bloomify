import { useState } from 'react';
import './FlowerSidebar.css';
import { apiFetch } from './api';

function FlowerAddForm({ onFlowerAdded, onTokenRefresh }) {
    const [form, setForm] = useState({ name: '', description: '', price: '', stock: '' });
    const [loading, setLoading] = useState(false);
    const [success, setSuccess] = useState(false);
    const [error, setError] = useState('');

    const handleChange = (e) => {
        setForm({ ...form, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        setSuccess(false);
        try {
            const response = await apiFetch('/api/v1/flowers/', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    name: form.name,
                    description: form.description,
                    price: parseFloat(form.price),
                    stock: parseInt(form.stock, 10)
                })
            }, onTokenRefresh);
            if (response.ok) {
                setSuccess(true);
                setForm({ name: '', description: '', price: '', stock: '' });
                if (onFlowerAdded) onFlowerAdded();
            } else {
                const err = await response.json();
                setError(err.message || 'Ошибка добавления');
            }
        } catch (err) {
            setError('Ошибка соединения с сервером');
        }
        setLoading(false);
    };

    return (
        <form className="flower-form" onSubmit={handleSubmit} style={{minWidth:260}}>
            <div className="flower-form-title">Добавить цветок</div>
            <input className="flower-input" name="name" type="text" placeholder="Название" value={form.name} onChange={handleChange} required />
            <input className="flower-input" name="description" type="text" placeholder="Описание" value={form.description} onChange={handleChange} required />
            <input className="flower-input" name="price" type="number" min="0" step="0.01" placeholder="Цена" value={form.price} onChange={handleChange} required />
            <input className="flower-input" name="stock" type="number" min="0" step="1" placeholder="Остаток" value={form.stock} onChange={handleChange} required />
            <button className="flower-btn" type="submit" disabled={loading}>{loading ? 'Добавление...' : 'Добавить'}</button>
            {success && <div className="flower-success">Цветок добавлен!</div>}
            {error && <div className="flower-error">{error}</div>}
        </form>
    );
}

export default FlowerAddForm; 