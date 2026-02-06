import { writable } from 'svelte/store';

const TOKEN_KEY = 'asset_manager_token';
const USER_KEY = 'asset_manager_user';

export function createAuthStore() {
  // Initialize from localStorage
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
  };
}
