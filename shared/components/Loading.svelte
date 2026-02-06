<script>
  export let size = 'medium'; // small, medium, large
  export let text = 'Loading...';
  export let fullPage = false;
</script>

{#if fullPage}
  <div class="loading-overlay">
    <div class="loading-content">
      <span class="icon is-{size}">
        <i class="fas fa-spinner fa-spin fa-2x"></i>
      </span>
      {#if text}
        <p class="mt-3">{text}</p>
      {/if}
    </div>
  </div>
{:else}
  <div class="has-text-centered">
    <span class="icon is-{size}">
      <i class="fas fa-spinner fa-spin"></i>
    </span>
    {#if text}
      <span class="ml-2">{text}</span>
    {/if}
  </div>
{/if}

<style>
  .loading-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.9);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
  }
  .loading-content {
    text-align: center;
  }
</style>
