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

  let form = { Name: '', DataType: 'string', EnumOptions: '' };

  const dataTypes = [
    { value: 'string', label: 'Text' },
    { value: 'int', label: 'Integer' },
    { value: 'decimal', label: 'Decimal' },
    { value: 'boolean', label: 'Yes/No' },
    { value: 'date', label: 'Date' },
    { value: 'datetime', label: 'Date & Time' },
    { value: 'enum', label: 'Dropdown' }
  ];

  const columns = [
    { key: 'Name', label: 'Name', sortable: true },
    { key: 'DataType', label: 'Data Type', render: (v) => dataTypes.find(d => d.value === v)?.label || v },
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
      const result = await api.getProperties();
      items = result || [];
    } catch (err) {
      notifications.error('Failed to load properties');
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
    form = { Name: '', DataType: 'string', EnumOptions: '' };
    showModal = true;
  }

  function openEdit(item) {
    editing = item;
    let enumOpts = '';
    if (item.EnumOptions) {
      try {
        enumOpts = JSON.parse(item.EnumOptions).join(', ');
      } catch {}
    }
    form = { Name: item.Name, DataType: item.DataType, EnumOptions: enumOpts };
    showModal = true;
  }

  function confirmDelete(item) {
    deleteTarget = item;
    showDeleteConfirm = true;
  }

  async function handleSave() {
    saving = true;
    try {
      const data = { ...form };
      if (form.DataType === 'enum' && form.EnumOptions) {
        data.EnumOptions = JSON.stringify(form.EnumOptions.split(',').map(s => s.trim()).filter(Boolean));
      } else {
        data.EnumOptions = '';
      }
      
      if (editing) {
        await api.updateProperty(editing.ID, data);
        notifications.success('Property updated');
      } else {
        await api.createProperty(data);
        notifications.success('Property created');
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
      await api.deleteProperty(deleteTarget.ID);
      notifications.success('Property deleted');
      showDeleteConfirm = false;
      await loadData();
    } catch (err) {
      notifications.error(err.message);
    }
  }
</script>

<div class="level">
  <div class="level-left">
    <h1 class="title">Properties</h1>
  </div>
  <div class="level-right">
    <Button color="primary" on:click={openNew}>
      <span class="icon"><i class="fas fa-plus"></i></span>
      <span>New Property</span>
    </Button>
  </div>
</div>

<p class="subtitle">Custom properties that can be attached to assets</p>

<Card>
  <DataTable {columns} data={items} {loading} emptyMessage="No properties found" />
</Card>

<Modal bind:active={showModal} title={editing ? 'Edit Property' : 'New Property'} size="small">
  <form on:submit|preventDefault={handleSave}>
    <FormField label="Name" name="name" bind:value={form.Name} required />
    <FormField 
      label="Data Type" 
      type="select" 
      name="dataType" 
      bind:value={form.DataType} 
      options={dataTypes}
      required 
    />
    {#if form.DataType === 'enum'}
      <FormField 
        label="Options" 
        name="enumOptions" 
        bind:value={form.EnumOptions}
        help="Comma-separated list of options"
        placeholder="Option 1, Option 2, Option 3"
      />
    {/if}
  </form>
  
  <svelte:fragment slot="footer">
    <Button color="primary" loading={saving} on:click={handleSave}>Save</Button>
    <Button on:click={() => showModal = false}>Cancel</Button>
  </svelte:fragment>
</Modal>

<ConfirmDialog
  bind:active={showDeleteConfirm}
  title="Delete Property"
  message="Are you sure you want to delete this property?"
  onConfirm={handleDelete}
/>
