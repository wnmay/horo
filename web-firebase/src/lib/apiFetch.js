// lib/apiFetch.js
import { getCurrentIdToken } from "../firebase/auth";

export async function apiFetch(url, options = {}) {
    const token = await getCurrentIdToken(); // null if not logged in
    const headers = new Headers(options.headers || {});
    if (token) headers.set("Authorization", `Bearer ${token}`);

    const res = await fetch(url, { ...options, headers });

    if (res.status === 401) {
        // one retry with forced refresh
        const retryToken = await getCurrentIdToken(true);
        if (retryToken) {
            const retryHeaders = new Headers(options.headers || {});
            retryHeaders.set("Authorization", `Bearer ${retryToken}`);
            return fetch(url, { ...options, headers: retryHeaders });
        }
    }
    return res;
}
