import { createAuthStore, createNotificationStore, createApiClient } from '../../shared/index.js';

// Create stores
export const auth = createAuthStore();
export const notifications = createNotificationStore();

// Create API client
export const api = createApiClient(
  '', // Use relative URLs since we proxy through Vite
  () => auth.getToken(),
  () => {
    auth.logout();
    window.location.hash = '#/login';
  }
);
