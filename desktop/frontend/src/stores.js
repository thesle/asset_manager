import { writable, get } from 'svelte/store';
import { createApiClient, createNotificationStore } from '../../../shared/index.js';

// Config store for API URL
export const config = writable({
  apiUrl: '',
  configured: false
});

// Auth store (similar to web but uses config for API URL)
const TOKEN_KEY = 'asset_manager_token';
const USER_KEY = 'asset_manager_user';

function createDesktopAuthStore() {
  const storedToken = typeof localStorage !== 'undefined' ? localStorage.getItem(TOKEN_KEY) : null;
  const storedUser = typeof localStorage !== 'undefined' ? localStorage.getItem(USER_KEY) : null;

  const { subscribe, set, update } = writable({
    token: storedToken,
    user: storedUser ? JSON.parse(storedUser) : null,
    isAuthenticated: !!storedToken,
  });

  return {
    subscribe,
    
    login(token, user) {
      localStorage.setItem(TOKEN_KEY, token);
      localStorage.setItem(USER_KEY, JSON.stringify(user));
      set({
        token,
        user,
        isAuthenticated: true,
      });
    },

    logout() {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
      set({
        token: null,
        user: null,
        isAuthenticated: false,
      });
    },

    getToken() {
      let token = null;
      subscribe(state => {
        token = state.token;
      })();
      return token;
    },

    updateUser(user) {
      localStorage.setItem(USER_KEY, JSON.stringify(user));
      update(state => ({
        ...state,
        user,
      }));
    },

    setFromConfig(token) {
      if (token) {
        localStorage.setItem(TOKEN_KEY, token);
        update(state => ({
          ...state,
          token,
          isAuthenticated: true,
        }));
      }
    }
  };
}

export const auth = createDesktopAuthStore();
export const notifications = createNotificationStore();

// API client - created dynamically based on config
let apiClient = null;

export function getApi() {
  const cfg = get(config);
  if (!apiClient || apiClient._baseUrl !== cfg.apiUrl) {
    apiClient = createApiClient(
      cfg.apiUrl,
      () => auth.getToken(),
      () => {
        auth.logout();
        window.location.hash = '#/login';
      }
    );
    apiClient._baseUrl = cfg.apiUrl;
  }
  return apiClient;
}

// Export api as a proxy that delegates to getApi() for compatibility with web pages
export const api = new Proxy({}, {
  get(target, prop) {
    const client = getApi();
    if (typeof client[prop] === 'function') {
      return client[prop].bind(client);
    }
    return client[prop];
  }
});

// Initialize config from Wails backend
export async function initConfig() {
  if (window.go && window.go.main && window.go.main.App) {
    try {
      const hasConfig = await window.go.main.App.HasConfig();
      if (hasConfig) {
        const cfg = await window.go.main.App.GetConfig();
        config.set({
          apiUrl: cfg.APIUrl,
          configured: true
        });
        if (cfg.Token) {
          auth.setFromConfig(cfg.Token);
        }
      }
    } catch (err) {
      console.error('Failed to load config:', err);
    }
  }
}

// Save config to Wails backend
export async function saveConfig(apiUrl, token) {
  if (window.go && window.go.main && window.go.main.App) {
    try {
      await window.go.main.App.SaveConfig(apiUrl, token);
      config.set({
        apiUrl,
        configured: true
      });
    } catch (err) {
      console.error('Failed to save config:', err);
      throw err;
    }
  }
}
