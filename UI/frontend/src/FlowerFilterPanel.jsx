import { useState, useEffect, useRef } from 'react';
import './App.css';

function FlowerFilterPanel({ onFilter, onAdd }) {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [price, setPrice] = useState('');
    const [stock, setStock] = useState('');
    const debounceRef = useRef();

    useEffect(() => {
        if (debounceRef.current) clearTimeout(debounceRef.current);
        debounceRef.current = setTimeout(() => {
            onFilter({ name, description, price, stock });
        }, 300);
        return () => clearTimeout(debounceRef.current);
    }, [name, description, price, stock]);

    const handleReset = () => {
        setName(''); setDescription(''); setPrice(''); setStock('');
        onFilter({ name: '', description: '', price: '', stock: '' });
    };

    return (
        <div style={{overflowX:'auto', width:'100%'}}>
            <form className="filter-panel" onSubmit={e => e.preventDefault()} style={{display:'flex',flexWrap:'nowrap',alignItems:'center',gap:8,marginBottom:0,minWidth:600}}>
                <input
                    className="filter-input"
                    type="text"
                    placeholder="Название..."
                    value={name}
                    onChange={e => setName(e.target.value)}
                    style={{minWidth:90,maxWidth:140,padding:'6px 7px',fontSize:'0.98rem'}}
                />
                <input
                    className="filter-input"
                    type="text"
                    placeholder="Описание..."
                    value={description}
                    onChange={e => setDescription(e.target.value)}
                    style={{minWidth:90,maxWidth:140,padding:'6px 7px',fontSize:'0.98rem'}}
                />
                <input
                    className="filter-input"
                    type="number"
                    min="0"
                    placeholder="Цена..."
                    value={price}
                    onChange={e => setPrice(e.target.value)}
                    style={{minWidth:70,maxWidth:100,padding:'6px 7px',fontSize:'0.98rem'}}
                />
                <input
                    className="filter-input"
                    type="number"
                    min="0"
                    placeholder="Наличие..."
                    value={stock}
                    onChange={e => setStock(e.target.value)}
                    style={{minWidth:70,maxWidth:100,padding:'6px 7px',fontSize:'0.98rem'}}
                />
                <button type="button" className="filter-btn filter-btn-reset" onClick={handleReset} style={{padding:'6px 12px',fontSize:'0.98rem'}}>Сбросить</button>
                <div style={{flexShrink:0,marginLeft:42}}>
                  <button type="button" className="filter-btn" style={{background:'#e9eef6',color:'#1b2636',fontWeight:600,minWidth:90,padding:'7px 0',fontSize:'1.01rem',whiteSpace:'nowrap'}} onClick={onAdd}>Добавить</button>
                </div>
            </form>
        </div>
    );
}

export default FlowerFilterPanel; 