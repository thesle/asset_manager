<script>
  export let notifications = [];
  export let onRemove = () => {};
</script>

<div class="notifications-container">
  {#each notifications as notification (notification.id)}
    <div class="notification is-{notification.type}">
      <button class="delete" on:click={() => onRemove(notification.id)}></button>
      {notification.message}
    </div>
  {/each}
</div>

<style>
  .notifications-container {
    position: fixed;
    top: 60px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 1000;
    max-width: 400px;
  }
  .notification {
    margin-bottom: 0.5rem;
    animation: slideIn 0.3s ease-out;
  }
  @keyframes slideIn {
    from {
      transform: translateY(-100%);
      opacity: 0;
    }
    to {
      transform: translateY(0);
      opacity: 1;
    }
  }
</style>
