import { useState } from 'react';
import flowerLogo from './assets/images/flower-logo.png';

// Глоб для всех картинок
const images = import.meta.glob('./assets/images/*.{png,jpg,jpeg}', { eager: true, import: 'default' });

function getFlowerImage(name) {
    const png = images[`./assets/images/${name}.png`];
    const jpg = images[`./assets/images/${name}.jpg`];
    const jpeg = images[`./assets/images/${name}.jpeg`];
    return png || jpg || jpeg || flowerLogo;
}

function FlowerCard({ flower, inCartCount, onAddToCart, onRemoveFromCart }) {
    const imageSrc = getFlowerImage(flower.name);

    return (
        <li style={{background:'#f6f7f9', boxShadow:'0 2px 12px rgba(0,0,0,0.06)', borderRadius:12, padding:'18px 18px 14px 18px', display:'flex', flexDirection:'column', gap:6, position:'relative', zIndex:1}}>
            <div style={{display:'flex', justifyContent:'center', alignItems:'center', marginBottom:10}}>
                <div style={{width:100, height:100, background:'#f0f1f3', borderRadius:10, display:'flex', alignItems:'center', justifyContent:'center', overflow:'hidden'}}>
                    <img src={imageSrc} alt={flower.name} style={{maxWidth:'100%', maxHeight:'100%', objectFit:'cover'}} />
                </div>
            </div>
            
            <div style={{fontWeight:700, fontSize:'1.18rem', color:'#1b2636', marginBottom:2}}>{flower.name}</div>
            <div style={{color:'#6b7280', fontSize:'0.98rem', marginBottom:4}}>{flower.description}</div>
            <div style={{fontSize:'1.05rem', color:'#4f5d75'}}>Цена: <b>{flower.price}</b></div>
            <div style={{fontSize:'0.98rem', color:'#4f5d75'}}>В наличии: <b>{flower.stock}</b></div>
            <div style={{marginTop:'auto', display:'flex', justifyContent:'center', alignItems:'center', gap:8}}>
                {inCartCount > 0 ? (
                    <>
                        <button className="flower-btn" style={{background:'#e9fbe5', color:'#388e3c', border:'1.5px solid #b6e7b0', borderRadius:7, padding:'7px 18px', fontWeight:600, fontSize:'1.01rem', cursor:'default'}} disabled>
                            В корзине ({inCartCount})
                        </button>
                        <button className="flower-btn" style={{background:'#fff0f0', color:'#d32f2f', border:'1.5px solid #f3bcbc', borderRadius:7, padding:'7px 12px', fontWeight:600, fontSize:'1.01rem', cursor:'pointer'}} onClick={onRemoveFromCart}>
                            Убрать из корзины
                        </button>
                    </>
                ) : (
                    <button className="flower-btn" style={{background:'#e9fbe5', color:'#388e3c', border:'1.5px solid #b6e7b0', borderRadius:7, padding:'9px 22px', fontWeight:600, fontSize:'1.08rem', cursor:'pointer'}} onClick={onAddToCart}>
                        В корзину
                    </button>
                )}
            </div>
        </li>
    );
}

export default FlowerCard; 