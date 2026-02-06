import { writable } from 'svelte/store';

export function createNotificationStore() {
  const { subscribe, update } = writable([]);

  let nextId = 1;

  return {
    subscribe,

    add(message, type = 'info', timeout = 5000) {
      const id = nextId++;
      const notification = { id, message, type };
      
      update(notifications => [...notifications, notification]);

      if (timeout > 0) {
        setTimeout(() => {
          this.remove(id);
        }, timeout);
      }

      return id;
    },

    success(message, timeout = 5000) {
      return this.add(message, 'success', timeout);
    },

    error(message, timeout = 8000) {
      return this.add(message, 'danger', timeout);
    },

    warning(message, timeout = 6000) {
      return this.add(message, 'warning', timeout);
    },

    info(message, timeout = 5000) {
      return this.add(message, 'info', timeout);
    },

    remove(id) {
      update(notifications => notifications.filter(n => n.id !== id));
    },

    clear() {
      update(() => []);
    },
  };
}
