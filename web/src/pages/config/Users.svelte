<script>
  import { onMount } from 'svelte';
  import { api, notifications } from '../../stores.js';
  import Card from '../../../../shared/components/Card.svelte';
  import DataTable from '../../../../shared/components/DataTable.svelte';
  import Button from '../../../../shared/components/Button.svelte';
  import Modal from '../../../../shared/components/Modal.svelte';
  import FormField from '../../../../shared/components/FormField.svelte';
  import ConfirmDialog from '../../../../shared/components/ConfirmDialog.svelte';

  let users = [];
  let loading = true;
  let showModal = false;
  let showDeleteConfirm = false;
  let showResetPassword = false;
  let editing = null;
  let deleteTarget = null;
  let resetTarget = null;
  let saving = false;

  let form = { Username: '', Email: '', Password: '', IsActive: true };
  let newPassword = '';

  const columns = [
    { key: 'Username', label: 'Username', sortable: true },
    { key: 'Email', label: 'Email', sortable: true },
    { key: 'IsActive', label: 'Active', render: (v) => v ? '<span class="tag is-success">Yes</span>' : '<span class="tag is-danger">No</span>' },
    { 
      key: 'actions', 
      label: 'Actions',
      render: (_, row) => `
        <div class="buttons are-small">
          <button class="button is-warning is-outlined edit-btn" data-id="${row.ID}">
            <span class="icon"><i class="fas fa-edit"></i></span>
          </button>
          <button class="button is-info is-outlined reset-btn" data-id="${row.ID}">
            <span class="icon"><i class="fas fa-key"></i></span>
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
      const result = await api.getUsers();
      users = result || [];
    } catch (err) {
      notifications.error('Failed to load users');
    } finally {
      loading = false;
    }
  }

  function handleTableClick(e) {
    const editBtn = e.target.closest('.edit-btn');
    const deleteBtn = e.target.closest('.delete-btn');
    const resetBtn = e.target.closest('.reset-btn');
    
    if (editBtn) {
      const id = parseInt(editBtn.dataset.id);
      openEdit(users.find(u => u.ID === id));
    } else if (deleteBtn) {
      const id = parseInt(deleteBtn.dataset.id);
      confirmDelete(users.find(u => u.ID === id));
    } else if (resetBtn) {
      const id = parseInt(resetBtn.dataset.id);
      openResetPassword(users.find(u => u.ID === id));
    }
  }

  function openNew() {
    editing = null;
    form = { Username: '', Email: '', Password: '', IsActive: true };
    showModal = true;
  }

  function openEdit(user) {
    editing = user;
    form = { Username: user.Username, Email: user.Email, Password: '', IsActive: user.IsActive };
    showModal = true;
  }

  function openResetPassword(user) {
    resetTarget = user;
    newPassword = '';
    showResetPassword = true;
  }

  function confirmDelete(user) {
    deleteTarget = user;
    showDeleteConfirm = true;
  }

  async function handleSave() {
    saving = true;
    try {
      if (editing) {
        await api.updateUser(editing.ID, { Username: form.Username, Email: form.Email, IsActive: form.IsActive });
        notifications.success('User updated');
      } else {
        await api.createUser(form);
        notifications.success('User created');
      }
      showModal = false;
      await loadData();
    } catch (err) {
      notifications.error(err.message);
    } finally {
      saving = false;
    }
  }

  async function handleResetPassword() {
    if (!newPassword || newPassword.length < 6) {
      notifications.error('Password must be at least 6 characters');
      return;
    }
    try {
      await api.resetUserPassword(resetTarget.ID, newPassword);
      notifications.success('Password reset successfully');
      showResetPassword = false;
    } catch (err) {
      notifications.error(err.message);
    }
  }

  async function handleDelete() {
    try {
      await api.deleteUser(deleteTarget.ID);
      notifications.success('User deleted');
      showDeleteConfirm = false;
      await loadData();
    } catch (err) {
      notifications.error(err.message);
    }
  }
</script>

<div class="level">
  <div class="level-left">
    <h1 class="title">Users</h1>
  </div>
  <div class="level-right">
    <Button color="primary" on:click={openNew}>
      <span class="icon"><i class="fas fa-plus"></i></span>
      <span>New User</span>
    </Button>
  </div>
</div>

<Card>
  <DataTable {columns} data={users} {loading} emptyMessage="No users found" />
</Card>

<Modal bind:active={showModal} title={editing ? 'Edit User' : 'New User'} size="small">
  <form on:submit|preventDefault={handleSave}>
    <FormField label="Username" name="username" bind:value={form.Username} required />
    <FormField label="Email" type="email" name="email" bind:value={form.Email} required />
    {#if !editing}
      <FormField label="Password" type="password" name="password" bind:value={form.Password} required />
    {/if}
    <FormField type="checkbox" name="isActive" bind:value={form.IsActive} placeholder="Active" />
  </form>
  
  <svelte:fragment slot="footer">
    <Button color="primary" loading={saving} on:click={handleSave}>Save</Button>
    <Button on:click={() => showModal = false}>Cancel</Button>
  </svelte:fragment>
</Modal>

<Modal bind:active={showResetPassword} title="Reset Password" size="small">
  {#if resetTarget}
    <p class="mb-4">Reset password for <strong>{resetTarget.Username}</strong></p>
    <FormField label="New Password" type="password" name="newPassword" bind:value={newPassword} required />
  {/if}
  
  <svelte:fragment slot="footer">
    <Button color="primary" on:click={handleResetPassword}>Reset Password</Button>
    <Button on:click={() => showResetPassword = false}>Cancel</Button>
  </svelte:fragment>
</Modal>

<ConfirmDialog
  bind:active={showDeleteConfirm}
  title="Delete User"
  message="Are you sure you want to delete this user?"
  onConfirm={handleDelete}
/>
