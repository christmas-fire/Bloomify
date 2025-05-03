// Глобальные переменные для управления обновлением токена
let isRefreshing = false;
let refreshTokenPromise = null;

let _getAccessToken = () => null;

export function setAccessTokenGetter(fn) {
    _getAccessToken = fn;
}

export function setAccessToken(token) {
    // для обратной совместимости, ничего не делаем
}

export function getAccessToken() {
    return _getAccessToken();
}

export function setRefreshToken(token) {
    localStorage.setItem('refreshToken', token);
}

export function getRefreshToken() {
    return localStorage.getItem('refreshToken');
}

// Функция для выполнения обновления токена
async function performTokenRefresh(baseURL, onTokenRefresh) {
    try {
        const refreshResp = await fetch(baseURL + '/auth/refresh', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: getRefreshToken() })
        });

        if (refreshResp.ok) {
            const tokens = await refreshResp.json();
            if (onTokenRefresh) onTokenRefresh(tokens.access_token); // Обновляем access token в App
            setRefreshToken(tokens.refresh_token); // Обновляем refresh token в localStorage
            return tokens.access_token; // Возвращаем новый access token
        } else {
            setRefreshToken(''); // Очищаем старый refresh token
            throw new Error('Session expired or invalid refresh token');
        }
    } catch (refreshErr) {
        console.error("Error during token refresh:", refreshErr);
        setRefreshToken(''); // Убедимся, что токен очищен при любой ошибке
        throw refreshErr; // Перебрасываем ошибку
    }
}

export async function apiFetch(path, options = {}, onTokenRefresh) {
    const isDesktop = typeof window !== 'undefined' && window.runtime !== undefined;
    const baseURL = isDesktop ? 'http://127.0.0.1:8080' : '';
    
    async function executeFetch(isRetry = false) {
        let headers = { ...(options.headers || {}) }; // Клонируем заголовки
        const token = getAccessToken();
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        try {
            const response = await fetch(baseURL + path, { ...options, headers });

            if (response.status === 401 && getRefreshToken() && !isRetry) {
                 // Получили 401, есть refresh token и это не повторный запрос после обновления
                
                if (!isRefreshing) {
                    // Если никто другой не обновляет токен, начинаем обновление
                    isRefreshing = true;
                    refreshTokenPromise = performTokenRefresh(baseURL, onTokenRefresh)
                        .finally(() => {
                            isRefreshing = false; // Сбрасываем флаг после завершения
                        });
                }

                // Ждем завершения текущего (или только что запущенного) процесса обновления
                try {
                    await refreshTokenPromise; // Ждем успешного обновления
                    // Если обновление прошло успешно, повторяем исходный запрос
                    return executeFetch(true); // Передаем флаг, что это повторный запрос
                } catch (refreshError) {
                     // Ошибка при обновлении токена (например, сессия истекла)
                    // Не повторяем запрос, просто выбрасываем ошибку
                    console.error("Token refresh failed, original request aborted:", path);
                    throw refreshError; // Передаем ошибку дальше
                }
            } else if (!response.ok && response.status !== 401) {
                 // Обработка других ошибок сервера (не 401)
                 // Попытаемся прочитать тело ошибки для более детальной информации
                 let errorData = { message: `Request failed with status ${response.status}` };
                 try {
                     errorData = await response.json();
                 } catch (e) { /* игнорируем, если тело не JSON */ }
                 console.error("API Fetch error:", errorData);
                 // Создаем объект ошибки, похожий на то, что ожидается в catch блоках компонентов
                 const error = new Error(errorData.message || `HTTP error ${response.status}`);
                 error.response = response; // Добавляем сам ответ для возможного анализа
                 error.errorBody = errorData;
                 throw error;
             }

            // Если ответ OK или это был 401 при повторном запросе (что не должно случиться, но на всякий случай)
            return response; // Возвращаем успешный ответ или ошибку 401, если обновление не помогло

        } catch (err) {
            // Ловим ошибки сети или ошибки, выброшенные из блока try
             console.error("Network or other error in apiFetch for path:", path, err);
            // Перебрасываем ошибку, чтобы ее могли поймать вызывающие функции
             // Убедимся, что у ошибки есть поле response, если это возможно
             if (!err.response) {
                 err.response = { status: 0, statusText: 'Network Error' }; // Пример для сетевой ошибки
             }
            throw err; 
        }
    }

    return executeFetch();
} 