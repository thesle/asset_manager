<script>
  import { onMount } from 'svelte';
  import { api, notifications } from '../../stores.js';
  import Card from '../../../../shared/components/Card.svelte';
  import DataTable from '../../../../shared/components/DataTable.svelte';
  import Button from '../../../../shared/components/Button.svelte';
  import Modal from '../../../../shared/components/Modal.svelte';
  import FormField from '../../../../shared/components/FormField.svelte';
  import ConfirmDialog from '../../../../shared/components/ConfirmDialog.svelte';

  let items = [];
  let loading = true;
  let showModal = false;
  let showDeleteConfirm = false;
  let editing = null;
  let deleteTarget = null;
  let saving = false;

  let form = { Name: '', Description: '' };

  const columns = [
    { key: 'Name', label: 'Name', sortable: true },
    { key: 'Description', label: 'Description' },
    { 
      key: 'actions', 
      label: 'Actions',
      render: (_, row) => `
        <div class="buttons are-small">
          <button class="button is-warning is-outlined edit-btn" data-id="${row.ID}">
            <span class="icon"><i class="fas fa-edit"></i></span>
          </button>
          <button class="button is-danger is-outlined delete-btn" data-id="${row.ID}">
            <span class="icon"><i class="fas fa-trash"></i></span>
          </button>
        </div>
      `
    }
  ];

  onMount(async () => {
    await loadData();
    document.addEventListener('click', handleTableClick);
    return () => document.removeEventListener('click', handleTableClick);
  });

  async function loadData() {
    loading = true;
    try {
      const result = await api.getAssetTypes();
      items = result || [];
    } catch (err) {
      notifications.error('Failed to load asset types');
    } finally {
      loading = false;
    }
  }

  function handleTableClick(e) {
    const editBtn = e.target.closest('.edit-btn');
    const deleteBtn = e.target.closest('.delete-btn');
    
    if (editBtn) {
      const id = parseInt(editBtn.dataset.id);
      openEdit(items.find(i => i.ID === id));
    } else if (deleteBtn) {
      const id = parseInt(deleteBtn.dataset.id);
      confirmDelete(items.find(i => i.ID === id));
    }
  }

  function openNew() {
    editing = null;
    form = { Name: '', Description: '' };
    showModal = true;
  }

  function openEdit(item) {
    editing = item;
    form = { Name: item.Name, Description: item.Description || '' };
    showModal = true;
  }

  function confirmDelete(item) {
    deleteTarget = item;
    showDeleteConfirm = true;
  }

  async function handleSave() {
    saving = true;
    try {
      if (editing) {
        await api.updateAssetType(editing.ID, form);
        notifications.success('Asset type updated');
      } else {
        await api.createAssetType(form);
        notifications.success('Asset type created');
      }
      showModal = false;
      await loadData();
    } catch (err) {
      notifications.error(err.message);
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    try {
      await api.deleteAssetType(deleteTarget.ID);
      notifications.success('Asset type deleted');
      showDeleteConfirm = false;
      await loadData();
    } catch (err) {
      notifications.error(err.message);
    }
  }
</script>

<div class="level">
  <div class="level-left">
    <h1 class="title">Asset Types</h1>
  </div>
  <div class="level-right">
    <Button color="primary" on:click={openNew}>
      <span class="icon"><i class="fas fa-plus"></i></span>
      <span>New Asset Type</span>
    </Button>
  </div>
</div>

<Card>
  <DataTable {columns} data={items} {loading} emptyMessage="No asset types found" />
</Card>

<Modal bind:active={showModal} title={editing ? 'Edit Asset Type' : 'New Asset Type'} size="small">
  <form on:submit|preventDefault={handleSave}>
    <FormField label="Name" name="name" bind:value={form.Name} required />
    <FormField label="Description" type="textarea" name="description" bind:value={form.Description} />
  </form>
  
  <svelte:fragment slot="footer">
    <Button color="primary" loading={saving} on:click={handleSave}>Save</Button>
    <Button on:click={() => showModal = false}>Cancel</Button>
  </svelte:fragment>
</Modal>

<ConfirmDialog
  bind:active={showDeleteConfirm}
  title="Delete Asset Type"
  message="Are you sure you want to delete this asset type?"
  onConfirm={handleDelete}
/>
