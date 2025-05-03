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
    const [cartOpen, setCartOpen] = useState(false);
    const [cartError, setCartError] = useState(''); // Новое состояние для ошибки корзины
    const [showPurchaseSuccess, setShowPurchaseSuccess] = useState(false); // Состояние для модалки "Спасибо за покупку"

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
            const response = await apiFetch('/api/v1/orders/', {
                headers: { 'Authorization': `Bearer ${accessToken}` }
            }, onTokenRefresh);
            if (response.ok) {
                const data = await response.json();
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
            } else {
                console.error('[MainScreen] fetchCart failed with status:', response.status);
            }
        } catch (err) {
            console.error('[MainScreen] Error in fetchCart:', err);
        }
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
        setCartError(''); // Очищаем предыдущую ошибку
        try {
            const response = await apiFetch('/api/v1/orders/', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${accessToken}` },
                body: JSON.stringify({ flower_id: flower.id, quantity: 1 })
            }, onTokenRefresh);

            if (response.ok) {
                fetchCart(); // Обновляем корзину при успехе
            } else {
                // Обрабатываем ошибку от сервера
                try {
                    const errorBody = await response.json();
                    if (errorBody.message && errorBody.message.includes('insufficient stock')) {
                        // Формируем более читаемое сообщение
                        const flowerName = flower?.name || `цветок ID ${flower.id}`; 
                        // Пытаемся извлечь доступное количество из сообщения бэкенда
                        const availableMatch = errorBody.message.match(/Available: (\d+)/);
                        const availableStock = availableMatch ? availableMatch[1] : '?';
                        const userMessage = `Недостаточно '${flowerName}'. Доступно: ${availableStock}.`;
                        setCartError(userMessage); // Устанавливаем ошибку нехватки товара
                        setTimeout(() => setCartError(''), 5000);
                    } else {
                        console.error('[MainScreen] handleAddToCart failed:', errorBody.message || response.status);
                        setCartError('Не удалось добавить товар.'); // Общая ошибка
                        setTimeout(() => setCartError(''), 5000);
                    }
                } catch (jsonError) {
                    console.error('[MainScreen] handleAddToCart failed to parse error JSON:', jsonError);
                    console.error('[MainScreen] handleAddToCart response status:', response.status);
                    setCartError('Не удалось обработать ответ сервера.');
                     setTimeout(() => setCartError(''), 5000);
                }
            }
        } catch (err) {
            console.error('[MainScreen] Network error in handleAddToCart:', err);
            setCartError('Ошибка сети при добавлении товара.');
             setTimeout(() => setCartError(''), 5000);
        }
    };
    // Удалить из корзины
    const handleRemoveFromCart = async (flower) => {
        setCartError(''); // Очищаем ошибку
        try {
            const response = await apiFetch(`/api/v1/orders/flower/${flower.id}/`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${accessToken}` }
            }, onTokenRefresh);
            if (response.ok) {
                fetchCart();
            } else {
                console.error(`[MainScreen] Failed to remove flower ${flower.id}. Status:`, response.status);
            }
        } catch (err) {
            console.error('[MainScreen] Error in handleRemoveFromCart:', err);
        }
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

    // TODO: Реализовать эндпоинты на бэкенде PATCH /api/v1/orders/flower/:id/increment и decrement
    const handleIncreaseQuantity = async (flowerId) => {
        setCartError(''); // Очищаем предыдущую ошибку
        try {
            const response = await apiFetch(`/api/v1/orders/flower/${flowerId}/increment/`, {
                 method: 'PATCH',
                 headers: { 'Authorization': `Bearer ${accessToken}` }
             }, onTokenRefresh);
             if (response.ok) {
                fetchCart(); // Обновляем корзину при успехе
             } else {
                // Обрабатываем ошибку от сервера
                try {
                    const errorBody = await response.json();
                    if (errorBody.message && errorBody.message.includes('insufficient stock')) {
                         // Находим имя цветка
                        const flower = safeFlowers.find(f => f.id === parseInt(flowerId, 10));
                        const flowerName = flower?.name || `цветок ID ${flowerId}`;
                        // Извлекаем детали из сообщения бэкенда
                        const availableMatch = errorBody.message.match(/Available: (\d+)/);
                        const inCartMatch = errorBody.message.match(/In cart: (\d+)/);
                        const availableStock = availableMatch ? availableMatch[1] : '?';
                        const inCartQuantity = inCartMatch ? inCartMatch[1] : '?';

                        const userMessage = `Недостаточно '${flowerName}'. Доступно: ${availableStock}, в корзине: ${inCartQuantity}.`;
                        setCartError(userMessage); // Устанавливаем ошибку нехватки товара
                        setTimeout(() => setCartError(''), 5000);
                    } else {
                        console.error(`[MainScreen] Failed to increment flower ${flowerId}:`, errorBody.message || response.status);
                        setCartError('Не удалось увеличить количество.'); // Общая ошибка
                         setTimeout(() => setCartError(''), 5000);
                    }
                } catch (jsonError) {
                    console.error(`[MainScreen] Failed to increment flower ${flowerId}, could not parse error JSON:`, jsonError);
                    console.error(`[MainScreen] Failed to increment flower ${flowerId}, response status:`, response.status);
                    setCartError('Не удалось обработать ответ сервера.');
                     setTimeout(() => setCartError(''), 5000);
                }
             }
        } catch (err) {
            console.error('[MainScreen] Network error in handleIncreaseQuantity:', err);
            setCartError('Ошибка сети при увеличении количества.');
             setTimeout(() => setCartError(''), 5000);
        }
    };

    const handleDecreaseQuantity = async (flowerId) => {
        setCartError(''); // Очищаем ошибку
        try {
            const response = await apiFetch(`/api/v1/orders/flower/${flowerId}/decrement/`, {
                 method: 'PATCH',
                 headers: { 'Authorization': `Bearer ${accessToken}` }
            }, onTokenRefresh);
            if (response.ok || response.status === 204) {
                fetchCart();
            } else {
                console.error(`[MainScreen] Failed to decrement flower ${flowerId}. Status:`, response.status);
            }
        } catch (err) {
            console.error('[MainScreen] Error in handleDecreaseQuantity:', err);
        }
    };

    // Обработчик "покупки"
    const handlePurchase = async () => {
        setCartError(''); // Очищаем предыдущие ошибки
        try {
            const response = await apiFetch('/api/v1/orders/active', {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${accessToken}` }
            }, onTokenRefresh);

            if (response.ok || response.status === 204) {
                // Успешно удалили заказ на бэкенде
                // alert('Спасибо за покупку!'); // Убираем стандартный alert
                setCart({}); // Очищаем корзину на фронте
                setCartOpen(false); // Закрываем модалку корзины
                setShowPurchaseSuccess(true); // Показываем модалку "Спасибо"
                 // Опционально: закрыть модалку через 3 секунды
                setTimeout(() => setShowPurchaseSuccess(false), 3000);
            } else {
                // Обрабатываем ошибку от сервера
                let errorMsg = 'Не удалось завершить покупку.';
                try {
                    const errorBody = await response.json();
                    errorMsg = errorBody.message || errorMsg;
                } catch (e) { /* Игнорируем ошибку парсинга JSON */ }
                console.error('[MainScreen] handlePurchase failed:', errorMsg, `Status: ${response.status}`);
                setCartError(errorMsg);
                setTimeout(() => setCartError(''), 5000); // Показываем ошибку на 5 сек
            }
        } catch (err) {
            // Ошибка сети
            console.error('[MainScreen] Network error in handlePurchase:', err);
            setCartError('Ошибка сети при завершении покупки.');
             setTimeout(() => setCartError(''), 5000);
        }
    };

    const safeFlowers = Array.isArray(flowers) ? flowers : [];

    return (
        <div style={{display: 'flex', flexDirection: 'column', height: '100vh', minHeight: 500}}>
            {/* Заголовок */}
            <div style={{display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', padding: '18px 32px 0 32px'}}>
                {/* Левая часть: Название и Корзина */} 
                <div style={{display: 'flex', flexDirection: 'column', alignItems: 'flex-start'}}>
                    <div className="start-title" style={{fontSize: '2.2rem', marginBottom: 10, display: 'flex', alignItems: 'center', gap: '0.5em'}}>
                        <img src={flowerLogo} alt="Bloomify" style={{width: 38, height: 38, marginRight: 4}} />
                        <span>Bloomify</span>
                    </div>
                    {/* Кнопка Корзина */} 
                    <button
                        className="start-btn"
                        style={{
                            width: 'auto',
                            minWidth: 110,
                            padding: '8px 18px',
                            background: Object.values(cart).reduce((a,b)=>a+b,0) > 0 ? '#e9fbe5' : '#e9eef6',
                            color: Object.values(cart).reduce((a,b)=>a+b,0) > 0 ? '#388e3c' : '#1b2636',
                            border: Object.values(cart).reduce((a,b)=>a+b,0) > 0 ? '1.5px solid #b6e7b0' : '1.5px solid #e0e3e8',
                            fontWeight:600,
                            position:'relative'
                        }}
                        onClick={()=>setCartOpen(true)}
                    >
                        Корзина
                        {Object.values(cart).reduce((a,b)=>a+b,0) > 0 && (
                            <span style={{position:'absolute',top:-7,right:-12,background:'#388e3c',color:'#fff',borderRadius:'50%',fontSize:'0.98rem',fontWeight:700,minWidth:22,height:22,display:'flex',alignItems:'center',justifyContent:'center',padding:'0 6px',boxShadow:'0 1px 4px rgba(56,142,60,0.12)'}}>
                                {Object.values(cart).reduce((a,b)=>a+b,0)}
                            </span>
                        )}
                    </button>
                </div>
                {/* Правая часть: Профиль и Выйти */} 
                <div style={{display:'flex', flexDirection:'column', alignItems:'flex-end', gap:8}}>
                    <button className="start-btn" style={{width: 'auto', minWidth: 90, padding: '8px 18px', background:'#e9eef6', color:'#1b2636', border:'1.5px solid #e0e3e8', fontWeight:600}} onClick={() => { setProfileForm({ username: '', oldPassword: '', password: '', password2: '' }); setProfileOpen(true); }}>Профиль</button>
                    <button className="start-btn" style={{width: 'auto', minWidth: 90, padding: '8px 18px', background:'#f0f1f3', color:'#4f5d75', border:'1.5px solid #e0e3e8'}} onClick={onLogout}>Выйти</button>
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
            {/* Модальное окно корзины */}
            {cartOpen && (
                <div className="modal-overlay" style={{position:'fixed',top:0,left:0,width:'100vw',height:'100vh',background:'rgba(27,38,54,0.32)',zIndex:1500,display:'flex',alignItems:'center',justifyContent:'center'}}>
                    <div className="modal-content" style={{background:'#fff',borderRadius:16,padding:'36px 32px',boxShadow:'0 8px 32px rgba(0,0,0,0.18)',minWidth:320,maxWidth:420,position:'relative'}}>
                        <button onClick={()=>{setCartOpen(false); setCartError('');}} style={{position:'absolute',top:12,right:16,background:'none',border:'none',fontSize:22,cursor:'pointer',color:'#4f5d75'}}>&times;</button>
                        <div style={{fontSize:'1.4rem', fontWeight:700, color:'#1b2636', marginBottom:18, textAlign:'center'}}>Корзина</div>
                        {/* Отображение ошибки корзины (с резервированием места) */} 
                        <div 
                            style={{
                                color: 'red', 
                                textAlign: 'center', 
                                marginBottom: cartError ? 10 : 0, // Отступ только если есть ошибка
                                padding: cartError ? '8px' : '0', // Паддинг только если есть ошибка
                                background: cartError ? '#fff0f0' : 'transparent', 
                                border: cartError ? '1px solid #f3bcbc' : 'none', 
                                borderRadius: 6,
                                minHeight: '1.5em', // Резервируем высоту примерно под одну строку
                                fontSize: '0.95rem', // Чуть уменьшим шрифт ошибки
                                transition: 'all 0.2s ease-out' // Плавное появление/исчезновение
                            }}
                        >
                            {cartError} {/* Показываем текст ошибки, если он есть */}
                        </div>
                        {Object.keys(cart).length === 0 || Object.values(cart).reduce((a,b)=>a+b,0) === 0 ? (
                            <div style={{color:'#6b7280', textAlign:'center', fontSize:'1.08rem', marginTop: 24}}>
                                Корзина пуста
                            </div>
                        ) : (
                            <ul style={{listStyle:'none',padding:0,margin:0,marginBottom:18}}>
                                {Object.entries(cart).map(([flowerId, quantity]) => {
                                    const flower = safeFlowers.find(f => f.id === parseInt(flowerId, 10));
                                    if (!flower) return null;
                                    return (
                                        <li key={flowerId} style={{display:'flex',alignItems:'center',justifyContent:'space-between',gap:10,padding:'8px 0',borderBottom:'1px solid #f0f1f3'}}>
                                            <span style={{color:'#1b2636', flex: 1, textAlign: 'left'}}>{flower.name}</span>
                                            <div style={{display:'flex', alignItems:'center', gap: 8}}>
                                                <button onClick={() => handleDecreaseQuantity(flower.id)} className="flower-btn" style={{padding:'2px 8px', fontSize:'1rem', minWidth:26, background:'#f0f1f3', border:'1px solid #e0e3e8', color:'#4f5d75'}}>-</button>
                                                <span style={{color:'#4f5d75', minWidth: 15, textAlign: 'center'}}>{quantity}</span>
                                                <button onClick={() => handleIncreaseQuantity(flower.id)} className="flower-btn" style={{padding:'2px 8px', fontSize:'1rem', minWidth:26, background:'#f0f1f3', border:'1px solid #e0e3e8', color:'#4f5d75'}}>+</button>
                                            </div>
                                            <span style={{color:'#1b2636',fontWeight:600, minWidth: 50, textAlign: 'right'}}>{(flower.price * quantity).toFixed(2)}</span>
                                        </li>
                                    );
                                })}
                            </ul>
                        )}
                        {Object.values(cart).reduce((a,b)=>a+b,0) > 0 && (
                            <>
                                <div style={{fontWeight:600, fontSize:'1.12rem', textAlign:'right', marginBottom:10}}>
                                    Итого: {(safeFlowers.filter(f => cart[f.id]).reduce((sum, f) => sum + (f.price * cart[f.id]), 0)).toFixed(2)}
                                </div>
                                {/* Кнопка Купить (теперь внутри условия) */} 
                                <button 
                                    className="flower-btn" 
                                    style={{
                                        width:'100%', 
                                        marginTop:15, 
                                        padding: '10px',
                                        fontSize: '1.1rem',
                                        fontWeight: 600,
                                        background: '#388e3c',
                                        color: '#fff',
                                        border: 'none'
                                    }}
                                    onClick={handlePurchase}
                                >
                                    Купить
                                </button>
                            </>
                        )}
                    </div>
                </div>
            )}
            {/* Модальное окно "Спасибо за покупку" */} 
            {showPurchaseSuccess && (
                 <div className="modal-overlay" style={{position:'fixed',top:0,left:0,width:'100vw',height:'100vh',background:'rgba(27,38,54,0.32)',zIndex:2500,display:'flex',alignItems:'center',justifyContent:'center'}}>
                    <div className="modal-content" style={{background:'#fff',borderRadius:16,padding:'40px 36px',boxShadow:'0 8px 32px rgba(0,0,0,0.18)',minWidth:300,maxWidth:400,textAlign:'center',position:'relative'}}>
                         <button onClick={()=>setShowPurchaseSuccess(false)} style={{position:'absolute',top:12,right:16,background:'none',border:'none',fontSize:22,cursor:'pointer',color:'#4f5d75'}}>&times;</button>
                        <div style={{fontSize:'1.8rem', fontWeight:700, color:'#388e3c', marginBottom:12}}>🎉</div>
                        <div style={{fontSize:'1.3rem', fontWeight:600, color:'#1b2636', marginBottom:24}}>Спасибо за покупку!</div>
                         <button 
                             className="flower-btn" 
                             style={{padding:'8px 24px', fontSize:'1rem', fontWeight:600, background:'#e9eef6', color:'#1b2636', border:'1.5px solid #e0e3e8'}}
                            onClick={()=>setShowPurchaseSuccess(false)}
                         >
                             Закрыть
                         </button>
                    </div>
                </div>
            )}
        </div>
    );
}

export default MainScreen; 