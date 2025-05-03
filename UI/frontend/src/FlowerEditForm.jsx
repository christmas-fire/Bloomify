import { useState, useEffect } from 'react';
import './FlowerSidebar.css'; // Используем те же стили, что и для добавления
import { apiFetch } from './api';

function FlowerEditForm({ flower, onClose, onFlowerUpdated, accessToken, onTokenRefresh }) {
    const [form, setForm] = useState({ name: '', description: '', price: '', stock: '' });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    // Инициализация формы данными цветка при монтировании или изменении пропа flower
    useEffect(() => {
        if (flower) {
            setForm({
                name: flower.name || '',
                description: flower.description || '',
                price: flower.price !== undefined ? String(flower.price) : '',
                stock: flower.stock !== undefined ? String(flower.stock) : ''
            });
        }
    }, [flower]);

    const handleChange = (e) => {
        setForm({ ...form, [e.target.name]: e.target.value });
        // Сбрасываем сообщение об успехе/ошибке при изменении полей
        setError('');
        setSuccess('');
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!flower) return; // На всякий случай

        setLoading(true);
        setError('');
        setSuccess('');

        const updatePromises = [];
        const updatedFields = [];

        // Проверяем каждое поле и добавляем промис обновления, если оно изменилось
        if (form.name !== flower.name) {
            updatedFields.push('name');
            updatePromises.push(apiFetch(`/api/v1/flowers/${flower.id}/name`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${accessToken}` },
                body: JSON.stringify({ new_name: form.name })
            }, onTokenRefresh));
        }
        if (form.description !== flower.description) {
             updatedFields.push('description');
            updatePromises.push(apiFetch(`/api/v1/flowers/${flower.id}/description`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${accessToken}` },
                body: JSON.stringify({ new_description: form.description })
            }, onTokenRefresh));
        }
        const formPrice = parseFloat(form.price);
        if (formPrice !== flower.price) {
             updatedFields.push('price');
            updatePromises.push(apiFetch(`/api/v1/flowers/${flower.id}/price`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${accessToken}` },
                body: JSON.stringify({ new_price: formPrice })
            }, onTokenRefresh));
        }
        const formStock = parseInt(form.stock, 10);
        if (formStock !== flower.stock) {
             updatedFields.push('stock');
            updatePromises.push(apiFetch(`/api/v1/flowers/${flower.id}/stock`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${accessToken}` },
                body: JSON.stringify({ new_stock: formStock })
            }, onTokenRefresh));
        }

        if (updatePromises.length === 0) {
            setSuccess('Нет изменений для сохранения.');
            setLoading(false);
            setTimeout(() => setSuccess(''), 3000);
            return;
        }

        try {
            const results = await Promise.allSettled(updatePromises);

            const failedUpdates = results.filter(result => result.status === 'rejected' || (result.status === 'fulfilled' && !result.value.ok));

            if (failedUpdates.length > 0) {
                // Пытаемся получить сообщение об ошибке из первого неудачного запроса
                let errorMsg = 'Ошибка при обновлении некоторых полей.';
                const firstFailure = failedUpdates[0];
                 if (firstFailure.status === 'rejected') {
                    errorMsg = firstFailure.reason?.message || errorMsg;
                 } else if (firstFailure.value && !firstFailure.value.ok) {
                    try {
                        const body = await firstFailure.value.json();
                        errorMsg = body.message || errorMsg;
                    } catch { /* ignore json parsing error */ }
                }
                setError(errorMsg);
                console.error("Failed updates:", failedUpdates);
            } else {
                 setSuccess(`Данные успешно обновлены!`);
                 // Вызываем колбэк обновления после короткой паузы, чтобы пользователь увидел сообщение
                 setTimeout(() => {
                    if (onFlowerUpdated) onFlowerUpdated();
                 }, 1500);
            }

        } catch (err) {
            // Общая ошибка (маловероятно, т.к. используем Promise.allSettled)
             console.error("General error during flower update:", err);
            setError('Произошла ошибка при обновлении.');
        } finally {
            setLoading(false);
             // Скрываем сообщение об ошибке через 5 секунд
            if (error) setTimeout(() => setError(''), 5000);
        }
    };

    if (!flower) return null; // Не рендерим форму, если нет данных цветка

    return (
        <form className="flower-form" onSubmit={handleSubmit} style={{minWidth: 260}}>
            {/* Не добавляем заголовок формы, он будет в модалке */}
            <label className="flower-label" style={{ display: 'block', marginBottom: '3px', fontWeight: 500, color: '#4f5d75' }}>Название</label>
            <input className="flower-input" name="name" type="text" placeholder="Название" value={form.name} onChange={handleChange} required />
            
            <label className="flower-label" style={{ display: 'block', marginBottom: '3px', fontWeight: 500, color: '#4f5d75' }}>Описание</label>
            <input className="flower-input" name="description" type="text" placeholder="Описание" value={form.description} onChange={handleChange} required />
            
            <label className="flower-label" style={{ display: 'block', marginBottom: '3px', fontWeight: 500, color: '#4f5d75' }}>Цена</label>
            <input className="flower-input" name="price" type="number" min="0" step="0.01" placeholder="Цена" value={form.price} onChange={handleChange} required />
            
            <label className="flower-label" style={{ display: 'block', marginBottom: '3px', fontWeight: 500, color: '#4f5d75' }}>Остаток</label>
            <input className="flower-input" name="stock" type="number" min="0" step="1" placeholder="Остаток" value={form.stock} onChange={handleChange} required />
            
            <button className="flower-btn" type="submit" disabled={loading} style={{marginTop: 15}}>
                {loading ? 'Сохранение...' : 'Сохранить изменения'}
            </button>
            
            {success && <div className="flower-success" style={{marginTop: 10}}>{success}</div>}
            {error && <div className="flower-error" style={{marginTop: 10}}>{error}</div>}
        </form>
    );
}

export default FlowerEditForm; 