<script>
  import { onMount } from 'svelte';
  import { api, notifications } from '../stores.js';
  import Card from '../../../shared/components/Card.svelte';
  import DataTable from '../../../shared/components/DataTable.svelte';
  import Button from '../../../shared/components/Button.svelte';
  import Modal from '../../../shared/components/Modal.svelte';
  import FormField from '../../../shared/components/FormField.svelte';
  import SearchInput from '../../../shared/components/SearchInput.svelte';
  import ConfirmDialog from '../../../shared/components/ConfirmDialog.svelte';
  import CustomFields from '../../../shared/components/CustomFields.svelte';

  let persons = [];
  let attributes = [];
  let loading = true;
  let showModal = false;
  let showDeleteConfirm = false;
  let editingPerson = null;
  let deleteTarget = null;
  let saving = false;

  let form = { Name: '', Email: '', Phone: '' };
  let customFieldValues = {};

  const columns = [
    { key: 'Name', label: 'Name', sortable: true },
    { key: 'Email', label: 'Email', sortable: true },
    { key: 'Phone', label: 'Phone' },
    { 
      key: 'actions', 
      label: 'Actions',
      render: (_, row) => row.Name === 'Unassigned' ? '' : `
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
      const [personsResult, attrsResult] = await Promise.all([
        api.getPersons(),
        api.getAttributes()
      ]);
      persons = personsResult || [];
      attributes = attrsResult || [];
    } catch (err) {
      notifications.error('Failed to load persons');
    } finally {
      loading = false;
    }
  }

  function handleTableClick(e) {
    const editBtn = e.target.closest('.edit-btn');
    const deleteBtn = e.target.closest('.delete-btn');
    
    if (editBtn) {
      const id = parseInt(editBtn.dataset.id);
      openEdit(persons.find(p => p.ID === id));
    } else if (deleteBtn) {
      const id = parseInt(deleteBtn.dataset.id);
      confirmDelete(persons.find(p => p.ID === id));
    }
  }

  function openNew() {
    editingPerson = null;
    form = { Name: '', Email: '', Phone: '' };
    customFieldValues = {};
    showModal = true;
  }

  async function openEdit(person) {
    editingPerson = person;
    form = { Name: person.Name, Email: person.Email || '', Phone: person.Phone || '' };
    
    // Load existing attribute values
    customFieldValues = {};
    try {
      const personAttrs = await api.getPersonAttributes(person.ID);
      if (personAttrs) {
        for (const attr of personAttrs) {
          customFieldValues[attr.AttributeID] = attr.Value;
        }
      }
    } catch (err) {
      // Ignore - just show empty custom fields
    }
    showModal = true;
  }

  function confirmDelete(person) {
    deleteTarget = person;
    showDeleteConfirm = true;
  }

  async function handleSave() {
    saving = true;
    try {
      let personId;
      
      if (editingPerson) {
        await api.updatePerson(editingPerson.ID, form);
        personId = editingPerson.ID;
      } else {
        const created = await api.createPerson(form);
        personId = created.ID;
      }
      
      // Save custom field values
      for (const [attrId, value] of Object.entries(customFieldValues)) {
        if (value !== undefined && value !== '') {
          await api.setPersonAttribute(personId, { AttributeID: parseInt(attrId), Value: String(value) });
        }
      }
      
      notifications.success(editingPerson ? 'Person updated' : 'Person created');
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
      await api.deletePerson(deleteTarget.ID);
      notifications.success('Person deleted');
      showDeleteConfirm = false;
      await loadData();
    } catch (err) {
      notifications.error(err.message);
    }
  }

  async function handleSearch(term) {
    if (!term) {
      await loadData();
      return;
    }
    loading = true;
    try {
      const result = await api.searchPersons(term);
      persons = result || [];
    } catch (err) {
      notifications.error('Search failed');
    } finally {
      loading = false;
    }
  }
</script>

<div class="level">
  <div class="level-left">
    <h1 class="title">Persons</h1>
  </div>
  <div class="level-right">
    <Button color="primary" on:click={openNew}>
      <span class="icon"><i class="fas fa-plus"></i></span>
      <span>New Person</span>
    </Button>
  </div>
</div>

<Card>
  <div class="mb-4">
    <SearchInput placeholder="Search persons..." onSearch={handleSearch} />
  </div>
  
  <DataTable {columns} data={persons} {loading} emptyMessage="No persons found" />
</Card>

<Modal bind:active={showModal} title={editingPerson ? 'Edit Person' : 'New Person'} size="wide">
  <form on:submit|preventDefault={handleSave}>
    <div class="columns">
      <div class="column">
        <h6 class="title is-6 mb-3">Person Details</h6>
        <FormField label="Name" name="name" bind:value={form.Name} required />
        <FormField label="Email" type="email" name="email" bind:value={form.Email} />
        <FormField label="Phone" name="phone" bind:value={form.Phone} />
      </div>
      <div class="column">
        <h6 class="title is-6 mb-3">Custom Attributes</h6>
        <CustomFields 
          definitions={attributes} 
          bind:values={customFieldValues}
        />
      </div>
    </div>
  </form>
  
  <svelte:fragment slot="footer">
    <Button color="primary" loading={saving} on:click={handleSave}>Save</Button>
    <Button on:click={() => showModal = false}>Cancel</Button>
  </svelte:fragment>
</Modal>

<ConfirmDialog
  bind:active={showDeleteConfirm}
  title="Delete Person"
  message="Are you sure you want to delete this person?"
  onConfirm={handleDelete}
/>
