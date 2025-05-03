export function apiFetch(path, options) {
    const isDesktop = typeof window !== 'undefined' && window.runtime !== undefined;
    const baseURL = isDesktop ? 'http://127.0.0.1:8080' : '';
    return fetch(baseURL + path, options);
} 