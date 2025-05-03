import { useEffect, useState } from 'react';
import './App.css';
import { apiFetch } from './api';

function MainScreen({ accessToken, onLogout }) {
    const [flowers, setFlowers] = useState([]);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchFlowers = async () => {
            try {
                const response = await apiFetch('/api/v1/flowers/', {
                    headers: {
                        'Authorization': `Bearer ${accessToken}`
                    }
                });
                if (response.ok) {
                    const data = await response.json();
                    setFlowers(data);
                } else {
                    setError('Ошибка загрузки цветов');
                }
            } catch (err) {
                setError('Ошибка соединения с сервером');
            }
        };
        fetchFlowers();
    }, [accessToken]);

    return (
        <div className="start-screen" style={{maxWidth: 700, minHeight: 400, alignItems: 'stretch'}}>
            <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8}}>
                <div className="start-title" style={{fontSize: '2rem', marginBottom: 0}}>
                    <span role="img" aria-label="flower">🌸</span> Все цветы
                </div>
                <button className="start-btn" style={{width: 'auto', minWidth: 90, padding: '8px 18px', marginBottom: 0, background:'#f0f1f3', color:'#4f5d75', border:'1.5px solid #e0e3e8'}} onClick={onLogout}>Выйти</button>
            </div>
            {error && <div className="success-message" style={{color: 'red', margin:'16px 0'}}>{error}</div>}
            <div style={{marginTop: 18}}>
                {flowers.length === 0 && !error && <div style={{color:'#6b7280', textAlign:'center'}}>Нет цветов для отображения.</div>}
                <ul style={{padding: 0, listStyle: 'none', display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(260px, 1fr))', gap: '18px'}}>
                    {flowers.map((flower) => (
                        <li key={flower.id} style={{background:'#f6f7f9', boxShadow:'0 2px 12px rgba(0,0,0,0.06)', borderRadius:12, padding:'18px 18px 14px 18px', display:'flex', flexDirection:'column', gap:6}}>
                            <div style={{fontWeight:700, fontSize:'1.18rem', color:'#1b2636', marginBottom:2}}>{flower.name}</div>
                            <div style={{color:'#6b7280', fontSize:'0.98rem', marginBottom:4}}>{flower.description}</div>
                            <div style={{fontSize:'1.05rem', color:'#4f5d75'}}>Цена: <b>{flower.price}</b></div>
                            <div style={{fontSize:'0.98rem', color:'#4f5d75'}}>В наличии: <b>{flower.stock}</b></div>
                        </li>
                    ))}
                </ul>
            </div>
        </div>
    );
}

export default MainScreen; 