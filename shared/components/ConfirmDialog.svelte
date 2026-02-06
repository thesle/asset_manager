<script>
  import Modal from './Modal.svelte';
  import Button from './Button.svelte';

  export let active = false;
  export let title = 'Confirm';
  export let message = 'Are you sure?';
  export let confirmText = 'Confirm';
  export let cancelText = 'Cancel';
  export let confirmColor = 'danger';
  export let onConfirm = () => {};
  export let onCancel = () => {};

  let loading = false;

  async function handleConfirm() {
    loading = true;
    try {
      await onConfirm();
      active = false;
    } finally {
      loading = false;
    }
  }

  function handleCancel() {
    active = false;
    onCancel();
  }
</script>

<Modal bind:active {title} size="small" onClose={handleCancel}>
  <p>{message}</p>
  
  <svelte:fragment slot="footer">
    <Button color={confirmColor} {loading} on:click={handleConfirm}>
      {confirmText}
    </Button>
    <Button on:click={handleCancel}>
      {cancelText}
    </Button>
  </svelte:fragment>
</Modal>
