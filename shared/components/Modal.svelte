<script>
  export let active = false;
  export let title = '';
  export let size = ''; // '', 'small', 'large', 'wide'
  export let onClose = () => {};

  function handleClose() {
    active = false;
    onClose();
  }

  function handleKeydown(event) {
    if (event.key === 'Escape') {
      handleClose();
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="modal" class:is-active={active}>
  <div class="modal-background" on:click={handleClose} on:keydown={handleKeydown} role="button" tabindex="-1" aria-label="Close modal"></div>
  <div class="modal-card" class:is-small={size === 'small'} class:is-large={size === 'large'} class:is-wide={size === 'wide'}>
    <header class="modal-card-head">
      <p class="modal-card-title">{title}</p>
      <button class="delete" aria-label="close" on:click={handleClose}></button>
    </header>
    <section class="modal-card-body">
      <slot></slot>
    </section>
    {#if $$slots.footer}
      <footer class="modal-card-foot">
        <slot name="footer"></slot>
      </footer>
    {/if}
  </div>
</div>

<style>
  .modal-card.is-small {
    width: 400px;
  }
  .modal-card.is-large {
    width: 900px;
  }
  .modal-card.is-wide {
    width: 800px;
  }
  .modal-card {
    max-height: 90vh;
  }
  .modal-card-body {
    overflow-y: auto;
  }
</style>
