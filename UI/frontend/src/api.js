let accessToken = null;

export function setAccessToken(token) {
    accessToken = token;
}

export function getAccessToken() {
    return accessToken;
}

export function setRefreshToken(token) {
    localStorage.setItem('refreshToken', token);
}

export function getRefreshToken() {
    return localStorage.getItem('refreshToken');
}

export async function apiFetch(path, options = {}, onTokenRefresh) {
    const isDesktop = typeof window !== 'undefined' && window.runtime !== undefined;
    const baseURL = isDesktop ? 'http://127.0.0.1:8080' : '';
    let headers = options.headers || {};
    if (accessToken) {
        headers['Authorization'] = `Bearer ${accessToken}`;
    }
    try {
        let response = await fetch(baseURL + path, { ...options, headers });
        if (response.status === 401 && getRefreshToken()) {
            // Попробуем обновить access token
            const refreshResp = await fetch(baseURL + '/auth/refresh', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh_token: getRefreshToken() })
            });
            if (refreshResp.ok) {
                const tokens = await refreshResp.json();
                setAccessToken(tokens.access_token);
                setRefreshToken(tokens.refresh_token);
                if (onTokenRefresh) onTokenRefresh(tokens.access_token);
                // Повторяем исходный запрос с новым access token
                headers['Authorization'] = `Bearer ${tokens.access_token}`;
                response = await fetch(baseURL + path, { ...options, headers });
            } else {
                // refresh не сработал — разлогиниваем
                setAccessToken(null);
                setRefreshToken('');
                throw new Error('Session expired');
            }
        }
        return response;
    } catch (err) {
        throw err;
    }
} 