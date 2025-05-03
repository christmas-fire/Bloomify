import { useEffect, useState, useCallback } from 'react';
import './App.css';
import './FlowerSidebar.css';
import { apiFetch } from './api';
import FlowerSidebar from './FlowerSidebar';
import FlowerFilterPanel from './FlowerFilterPanel';
import FlowerAddForm from './FlowerAddForm';
import flowerLogo from './assets/images/flower-logo.png';
import FlowerCard from './FlowerCard';

function MainScreen({ accessToken, onLogout, onTokenRefresh }) {
    const [flowers, setFlowers] = useState([]);
    const [error, setError] = useState('');
    const [addModalOpen, setAddModalOpen] = useState(false);
    const [cart, setCart] = useState({}); // {flowerId: count}
    const [profileOpen, setProfileOpen] = useState(false);
    const [user, setUser] = useState({ username: '', email: '' });
    const [profileError, setProfileError] = useState('');
    const [profileSuccess, setProfileSuccess] = useState('');
    const [profileForm, setProfileForm] = useState({ username: '', oldPassword: '', password: '', password2: '' });
    const [changeUsernameOpen, setChangeUsernameOpen] = useState(false);
    const [changePasswordOpen, setChangePasswordOpen] = useState(false);

    const fetchFlowers = useCallback(async () => {
        try {
            const response = await apiFetch('/api/v1/flowers/', {
                headers: {
                    'Authorization': `Bearer ${accessToken}`
                }
            }, onTokenRefresh);
            if (response.ok) {
                const data = await response.json();
                setFlowers(Array.isArray(data) ? data : []);
                setError('');
            } else {
                setError('Ошибка загрузки цветов');
                setFlowers([]);
            }
        } catch (err) {
            setError('Ошибка соединения с сервером');
            setFlowers([]);
        }
    }, [accessToken, onTokenRefresh]);

    // Получить корзину пользователя
    const fetchCart = useCallback(async () => {
        try {
            const response = await apiFetch('/api/v1/orders/user_id', {
                headers: { 'Authorization': `Bearer ${accessToken}` }
            }, onTokenRefresh);
            if (response.ok) {
                const data = await response.json();
                // Преобразуем в {flowerId: count}
                const cartObj = {};
                if (Array.isArray(data)) {
                    data.forEach(order => {
                        if (order.flowers) {
                            order.flowers.forEach(f => {
                                cartObj[f.id] = f.quantity || 1;
                            });
                        }
                    });
                }
                setCart(cartObj);
            }
        } catch {}
    }, [accessToken, onTokenRefresh]);

    // Получить инфо о пользователе
    const fetchUser = useCallback(async () => {
        try {
            const resp = await apiFetch('/api/v1/users/me', {
                headers: { 'Authorization': `Bearer ${accessToken}` }
            }, onTokenRefresh);
            if (resp.ok) {
                const data = await resp.json();
                setUser(data);
            }
        } catch {}
    }, [accessToken, onTokenRefresh]);

    const handleFilter = async (filters) => {
        setError('');
        let url = '';
        if (filters.name) {
            url = '/api/v1/flowers/name?name=' + encodeURIComponent(filters.name);
        } else if (filters.description) {
            url = '/api/v1/flowers/description?description=' + encodeURIComponent(filters.description);
        } else if (filters.price) {
            url = '/api/v1/flowers/price?price=' + encodeURIComponent(filters.price);
        } else if (filters.stock) {
            url = '/api/v1/flowers/stock?stock=' + encodeURIComponent(filters.stock);
        } else {
            url = '/api/v1/flowers/';
        }
        try {
            const response = await apiFetch(url, {
                headers: {
                    'Authorization': `Bearer ${accessToken}`
                }
            }, onTokenRefresh);
            if (response.ok) {
                const data = await response.json();
                setFlowers(Array.isArray(data) ? data : []);
                setError('');
            } else {
                setError('Ошибка загрузки цветов');
                setFlowers([]);
            }
        } catch (err) {
            setError('Ошибка соединения с сервером');
            setFlowers([]);
        }
    };

    useEffect(() => {
        fetchFlowers();
        fetchCart();
        if (profileOpen) fetchUser();
    }, [fetchFlowers, fetchCart, profileOpen, fetchUser]);

    const handleFlowerAdded = () => {
        fetchFlowers();
    };

    // Добавить в корзину
    const handleAddToCart = async (flower) => {
        try {
            await apiFetch('/api/v1/orders/', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${accessToken}` },
                body: JSON.stringify({ flower_id: flower.id, quantity: 1 })
            }, onTokenRefresh);
            fetchCart();
        } catch {}
    };
    // Удалить из корзины
    const handleRemoveFromCart = async (flower) => {
        // Для простоты: ищем orderId по user_id, потом DELETE /api/v1/orders/:id
        // (или PATCH /api/v1/orders/:id/flower_id если реализовано)
        // Здесь просто обновим корзину после действия
        try {
            await apiFetch(`/api/v1/orders/`, {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${accessToken}` },
                body: JSON.stringify({ flower_id: flower.id })
            }, onTokenRefresh);
            fetchCart();
        } catch {}
    };

    // Обработка изменения профиля
    const handleProfileChange = e => {
        setProfileForm({ ...profileForm, [e.target.name]: e.target.value });
    };
    const handleProfileSave = async e => {
        e.preventDefault();
        setProfileError(''); setProfileSuccess('');
        // Смена имени
        if (profileForm.username && profileForm.username !== user.username) {
            const resp = await apiFetch(`/api/v1/users/${user.id}/username`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${accessToken}` },
                body: JSON.stringify({ old_username: user.username, new_username: profileForm.username })
            }, onTokenRefresh);
            if (!resp.ok) {
                setProfileError('Ошибка смены имени');
                return;
            }
            setProfileSuccess('Имя пользователя обновлено!');
        }
        // Смена пароля
        if (profileForm.password && profileForm.password === profileForm.password2) {
            const resp = await apiFetch(`/api/v1/users/${user.id}/password`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${accessToken}` },
                body: JSON.stringify({username:user.username, old_password: profileForm.oldPassword, new_password: profileForm.password })
            }, onTokenRefresh);
            if (!resp.ok) {
                setProfileError('Ошибка смены пароля');
                return;
            }
            setProfileSuccess('Пароль обновлён!');
        } else if (profileForm.password || profileForm.password2) {
            setProfileError('Пароли не совпадают');
            return;
        }
    };

    const safeFlowers = Array.isArray(flowers) ? flowers : [];

    return (
        <div style={{display: 'flex', flexDirection: 'column', height: '100vh', minHeight: 500}}>
            {/* Заголовок */}
            <div style={{display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '18px 32px 0 32px'}}>
                <div className="start-title" style={{fontSize: '2.2rem', marginBottom: 0, display: 'flex', alignItems: 'center', gap: '0.5em'}}>
                    <img src={flowerLogo} alt="Bloomify" style={{width: 38, height: 38, marginRight: 4}} />
                    <span>Bloomify</span>
                </div>
                <div style={{display:'flex', flexDirection:'column', alignItems:'flex-end', gap:8}}>
                    <button className="start-btn" style={{width: 'auto', minWidth: 90, padding: '8px 18px', marginBottom: 0, background:'#e9eef6', color:'#1b2636', border:'1.5px solid #e0e3e8', fontWeight:600}} onClick={() => { setProfileForm({ username: '', oldPassword: '', password: '', password2: '' }); setProfileOpen(true); }}>Профиль</button>
                    <button className="start-btn" style={{width: 'auto', minWidth: 90, padding: '8px 18px', marginBottom: 0, background:'#f0f1f3', color:'#4f5d75', border:'1.5px solid #e0e3e8'}} onClick={onLogout}>Выйти</button>
                </div>
            </div>
            {/* Фильтры и кнопка Добавить */}
            <div style={{display: 'flex', alignItems: 'flex-end', gap: 18, padding: '0 32px', marginTop: 18}}>
                <div style={{flex: 1}}>
                    <FlowerFilterPanel onFilter={handleFilter} onAdd={() => setAddModalOpen(true)} />
                </div>
            </div>
            {/* Содержимое + скролл */}
            <div style={{display: 'flex', flex: 1, width: '100%', minHeight: 0, marginTop: 0, overflow:'hidden'}}>
                <div style={{flex: 1, marginLeft: 0, minWidth: 0, padding: '0 32px', overflowY:'auto', height:'100%'}}>
                    <div style={{marginTop: 18}}>
                        {error && (
                            <div className="success-message" style={{color: 'red', margin:'16px 0'}}>{error}</div>
                        )}
                        {safeFlowers.length === 0 && (
                            <div style={{color:'#6b7280', textAlign:'center', fontSize:'1.1rem', marginTop: 40}}>
                                Пока нет цветов
                            </div>
                        )}
                        {safeFlowers.length > 0 && (
                            <ul style={{padding: 0, listStyle: 'none', display:'grid', gridTemplateColumns:'repeat(auto-fit, minmax(260px, 1fr))', gap: '18px'}}>
                                {safeFlowers.map((flower) => (
                                    <FlowerCard
                                        key={flower.id}
                                        flower={flower}
                                        inCartCount={cart[flower.id] || 0}
                                        onAddToCart={() => handleAddToCart(flower)}
                                        onRemoveFromCart={() => handleRemoveFromCart(flower)}
                                    />
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
            </div>
            {/* Футер */}
            <footer style={{width:'100%', background:'#f6f7f9', color:'#4f5d75', textAlign:'center', padding:'14px 0 10px 0', fontSize:'1.08rem', fontWeight:500, letterSpacing:'0.5px', borderTop:'1.5px solid #e0e3e8'}}>
                christmas-fire 2025
            </footer>
            {/* Модальное окно для добавления цветка */}
            {addModalOpen && (
                <div className="modal-overlay" style={{position:'fixed',top:0,left:0,width:'100vw',height:'100vh',background:'rgba(27,38,54,0.32)',zIndex:1000,display:'flex',alignItems:'center',justifyContent:'center'}}>
                    <div className="modal-content" style={{background:'#fff',borderRadius:16,padding:'32px 28px',boxShadow:'0 8px 32px rgba(0,0,0,0.18)',minWidth:320,maxWidth:400,position:'relative'}}>
                        <button onClick={()=>setAddModalOpen(false)} style={{position:'absolute',top:12,right:16,background:'none',border:'none',fontSize:22,cursor:'pointer',color:'#4f5d75'}}>&times;</button>
                        {/* Сюда вставим форму добавления цветка */}
                        {/* FlowerAddForm будет добавлен отдельно */}
                        <FlowerAddForm onFlowerAdded={()=>{setAddModalOpen(false);handleFlowerAdded();}} onTokenRefresh={onTokenRefresh} />
                    </div>
                </div>
            )}
            {/* Модальное окно профиля */}
            {profileOpen && (
                <div className="modal-overlay" style={{position:'fixed',top:0,left:0,width:'100vw',height:'100vh',background:'rgba(27,38,54,0.32)',zIndex:2000,display:'flex',alignItems:'center',justifyContent:'center'}}>
                    <div className="modal-content" style={{background:'#fff',borderRadius:16,padding:'40px 36px',boxShadow:'0 8px 32px rgba(0,0,0,0.18)',minWidth:340,maxWidth:480,width:'100%',position:'relative'}}>
                        <button onClick={()=>setProfileOpen(false)} style={{position:'absolute',top:12,right:16,background:'none',border:'none',fontSize:22,cursor:'pointer',color:'#4f5d75'}}>&times;</button>
                        <div style={{fontSize:'2rem', fontWeight:700, color:'#1b2636', marginBottom:32, textAlign:'center', letterSpacing:'0.5px'}}>
                            <span role="img" aria-label="flower">🌸</span> Профиль
                        </div>
                        <div style={{marginBottom:12, color:'#4f5d75', fontSize:'1.08rem'}}><b>Имя пользователя:</b> {user.username}</div>
                        <div style={{marginBottom:32, color:'#4f5d75', fontSize:'1.08rem'}}><b>Email:</b> {user.email}</div>
                        <button className="flower-btn" style={{width:'60%', alignSelf:'center', marginBottom:16}} onClick={()=>{ setProfileForm({ username: '', oldPassword: '', password: '', password2: '' }); setChangeUsernameOpen(true); }}>Сменить имя пользователя</button>
                        <button className="flower-btn" style={{width:'60%', alignSelf:'center', marginBottom:24}} onClick={()=>{ setProfileForm({ username: '', oldPassword: '', password: '', password2: '' }); setChangePasswordOpen(true); }}>Сменить пароль</button>
                        <button className="start-btn" style={{width:'60%', alignSelf:'center', background:'#f0f1f3', color:'#4f5d75', border:'1.5px solid #e0e3e8', marginTop:10}} onClick={onLogout}>Выйти</button>
                    </div>
                </div>
            )}
            {/* Модальное окно смены имени */}
            {changeUsernameOpen && (
                <div className="modal-overlay" style={{position:'fixed',top:0,left:0,width:'100vw',height:'100vh',background:'rgba(27,38,54,0.32)',zIndex:2100,display:'flex',alignItems:'center',justifyContent:'center'}}>
                    <div className="modal-content" style={{background:'#fff',borderRadius:16,padding:'40px 36px',boxShadow:'0 8px 32px rgba(0,0,0,0.18)',minWidth:340,maxWidth:480,width:'100%',position:'relative'}}>
                        <button onClick={()=>setChangeUsernameOpen(false)} style={{position:'absolute',top:12,right:16,background:'none',border:'none',fontSize:22,cursor:'pointer',color:'#4f5d75'}}>&times;</button>
                        <div style={{fontSize:'1.3rem', fontWeight:600, color:'#1b2636', marginBottom:18, textAlign:'center'}}>Сменить имя пользователя</div>
                        <form onSubmit={handleProfileSave} style={{display:'flex', flexDirection:'column', gap:10, marginBottom:18}}>
                            <label style={{color:'#6b7280', fontSize:'1rem'}}>Новое имя пользователя</label>
                            <input type="text" name="username" value={profileForm.username} onChange={handleProfileChange} className="flower-input" style={{marginBottom:6}} />
                            <button type="submit" className="flower-btn" style={{marginTop:10, fontWeight:600, fontSize:'1.08rem', width:'50%', alignSelf:'center'}}>Сохранить</button>
                        </form>
                        {profileError && <div style={{color:'#d32f2f', marginBottom:8, textAlign:'center'}}>{profileError}</div>}
                        {profileSuccess && <div style={{color:'#388e3c', marginBottom:8, textAlign:'center'}}>{profileSuccess}</div>}
                    </div>
                </div>
            )}
            {/* Модальное окно смены пароля */}
            {changePasswordOpen && (
                <div className="modal-overlay" style={{position:'fixed',top:0,left:0,width:'100vw',height:'100vh',background:'rgba(27,38,54,0.32)',zIndex:2100,display:'flex',alignItems:'center',justifyContent:'center'}}>
                    <div className="modal-content" style={{background:'#fff',borderRadius:16,padding:'40px 36px',boxShadow:'0 8px 32px rgba(0,0,0,0.18)',minWidth:340,maxWidth:480,width:'100%',position:'relative'}}>
                        <button onClick={()=>setChangePasswordOpen(false)} style={{position:'absolute',top:12,right:16,background:'none',border:'none',fontSize:22,cursor:'pointer',color:'#4f5d75'}}>&times;</button>
                        <div style={{fontSize:'1.3rem', fontWeight:600, color:'#1b2636', marginBottom:18, textAlign:'center'}}>Сменить пароль</div>
                        <form onSubmit={handleProfileSave} style={{display:'flex', flexDirection:'column', gap:10, marginBottom:18}}>
                            <label style={{color:'#6b7280', fontSize:'1rem'}}>Старый пароль</label>
                            <input type="password" name="oldPassword" value={profileForm.oldPassword} onChange={handleProfileChange} className="flower-input" />
                            <label style={{color:'#6b7280', fontSize:'1rem'}}>Новый пароль</label>
                            <input type="password" name="password" value={profileForm.password} onChange={handleProfileChange} className="flower-input" />
                            <label style={{color:'#6b7280', fontSize:'1rem'}}>Повторите пароль</label>
                            <input type="password" name="password2" value={profileForm.password2} onChange={handleProfileChange} className="flower-input" />
                            <button type="submit" className="flower-btn" style={{marginTop:10, fontWeight:600, fontSize:'1.08rem', width:'50%', alignSelf:'center'}}>Сохранить</button>
                        </form>
                        {profileError && <div style={{color:'#d32f2f', marginBottom:8, textAlign:'center'}}>{profileError}</div>}
                        {profileSuccess && <div style={{color:'#388e3c', marginBottom:8, textAlign:'center'}}>{profileSuccess}</div>}
                    </div>
                </div>
            )}
        </div>
    );
}

export default MainScreen; 