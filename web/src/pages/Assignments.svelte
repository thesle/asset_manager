<script>
  import { onMount } from 'svelte';
  import { api, notifications } from '../stores.js';
  import Card from '../../../shared/components/Card.svelte';
  import DataTable from '../../../shared/components/DataTable.svelte';
  import Button from '../../../shared/components/Button.svelte';
  import Modal from '../../../shared/components/Modal.svelte';
  import FormField from '../../../shared/components/FormField.svelte';
  import Loading from '../../../shared/components/Loading.svelte';
  import ConfirmDialog from '../../../shared/components/ConfirmDialog.svelte';
  import SearchInput from '../../../shared/components/SearchInput.svelte';

  let assets = [];
  let persons = [];
  let loading = true;
  let showAssignModal = false;
  let showEditModal = false;
  let showUnassignConfirm = false;
  let editingAssignment = null;
  let unassignTarget = null;
  let searchTerm = '';
  let showUnassignedOnly = false;

  let form = { AssetID: '', PersonID: '', Notes: '', EffectiveDate: '' };
  let editForm = { ID: '', AssetID: '', PersonID: '', EffectiveDate: '', Notes: '' };
  
  // Get today's date in YYYY-MM-DD format for the date input
  function getTodayDate() {
    const today = new Date();
    return today.toISOString().split('T')[0];
  }
  
  function formatDateForInput(dateStr) {
    if (!dateStr) return getTodayDate();
    const date = new Date(dateStr);
    return date.toISOString().split('T')[0];
  }

  const columns = [
    { 
      key: 'Name', 
      label: 'Asset', 
      sortable: true,
      render: (val, row) => `<a href="#/assets?edit=${row.ID}" class="has-text-link">${val}</a>`
    },
    { key: 'AssetTypeName', label: 'Type', sortable: true },
    { 
      key: 'PurchasedAt', 
      label: 'Purchased',
      render: (v) => v?.Time ? new Date(v.Time).toLocaleDateString() : '-'
    },
    { 
      key: 'CurrentAssignee', 
      label: 'Assigned To', 
      sortable: true,
      render: (val, row) => {
        if (!val || val === 'Unassigned') return `<span class="has-text-grey">${val || 'Unassigned'}</span>`;
        return `<a href="#/persons?edit=${row.CurrentAssigneeID}" class="has-text-link">${val}</a>`;
      }
    },
    { 
      key: 'AssignedFrom', 
      label: 'Assigned Since',
      render: (v) => v ? new Date(v).toLocaleDateString() : '-'
    },
    { 
      key: 'actions', 
      label: 'Actions',
      render: (_, row) => `
        <div class="buttons are-small">
          ${row.CurrentAssignee && row.CurrentAssignee !== 'Unassigned' ? `
            <button class="button is-info is-outlined edit-btn" data-id="${row.ID}">
              <span class="icon"><i class="fas fa-edit"></i></span>
            </button>
          ` : ''}
          <button class="button is-primary is-outlined assign-btn" data-id="${row.ID}">
            <span class="icon"><i class="fas fa-user-plus"></i></span>
            <span>${row.CurrentAssignee && row.CurrentAssignee !== 'Unassigned' ? 'Reassign' : 'Assign'}</span>
          </button>
          ${row.CurrentAssignee && row.CurrentAssignee !== 'Unassigned' ? `
            <button class="button is-warning is-outlined unassign-btn" data-id="${row.ID}">
              <span class="icon"><i class="fas fa-user-minus"></i></span>
            </button>
          ` : ''}
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
      const [assetsResult, personsResult] = await Promise.all([
        api.getAssetsWithAssignments(),
        api.getPersons()
      ]);
      assets = assetsResult || [];
      persons = personsResult || [];
    } catch (err) {
      notifications.error('Failed to load data');
    } finally {
      loading = false;
    }
  }

  function handleTableClick(e) {
    const assignBtn = e.target.closest('.assign-btn');
    const unassignBtn = e.target.closest('.unassign-btn');
    const editBtn = e.target.closest('.edit-btn');
    
    if (editBtn) {
      const id = parseInt(editBtn.dataset.id);
      openEdit(assets.find(a => a.ID === id));
    } else if (assignBtn) {
      const id = parseInt(assignBtn.dataset.id);
      openAssign(assets.find(a => a.ID === id));
    } else if (unassignBtn) {
      const id = parseInt(unassignBtn.dataset.id);
      confirmUnassign(assets.find(a => a.ID === id));
    }
  }

  function openAssign(asset) {
    form = { AssetID: asset.ID, PersonID: '', Notes: '', EffectiveDate: getTodayDate() };
    showAssignModal = true;
  }

  async function openEdit(asset) {
    // Fetch the current assignment for this asset
    try {
      const assignment = await api.getCurrentAssetAssignment(asset.ID);
      if (assignment) {
        editingAssignment = { ...assignment, AssetName: asset.Name };
        editForm = {
          ID: assignment.ID,
          AssetID: assignment.AssetID,
          PersonID: assignment.PersonID,
          EffectiveDate: formatDateForInput(assignment.EffectiveFrom),
          Notes: assignment.Notes || ''
        };
        showEditModal = true;
      }
    } catch (err) {
      notifications.error('Failed to load assignment details');
    }
  }

  function confirmUnassign(asset) {
    unassignTarget = asset;
    showUnassignConfirm = true;
  }

  async function handleAssign() {
    try {
      // Convert date string to ISO datetime for the API
      const effectiveDate = form.EffectiveDate ? new Date(form.EffectiveDate).toISOString() : null;
      await api.assignAsset(form.AssetID, parseInt(form.PersonID), form.Notes, effectiveDate);
      notifications.success('Asset assigned');
      showAssignModal = false;
      await loadData();
    } catch (err) {
      notifications.error(err.message);
    }
  }

  async function handleEditSave() {
    try {
      const effectiveFrom = editForm.EffectiveDate ? new Date(editForm.EffectiveDate).toISOString() : null;
      await api.updateAssignment(editForm.ID, {
        AssetID: parseInt(editForm.AssetID),
        PersonID: parseInt(editForm.PersonID),
        EffectiveFrom: effectiveFrom,
        Notes: editForm.Notes
      });
      notifications.success('Assignment updated');
      showEditModal = false;
      editingAssignment = null;
      await loadData();
    } catch (err) {
      notifications.error(err.message);
    }
  }

  async function handleUnassign() {
    if (!unassignTarget) return;
    try {
      await api.unassignAsset(unassignTarget.ID);
      notifications.success('Asset unassigned');
      showUnassignConfirm = false;
      unassignTarget = null;
      await loadData();
    } catch (err) {
      notifications.error(err.message);
    }
  }

  $: personOptions = persons.filter(p => p.Name !== 'Unassigned').map(p => ({ value: p.ID, label: p.Name }));
  $: selectedAsset = assets.find(a => a.ID === form.AssetID);
  $: editPersonOptions = persons.filter(p => p.Name !== 'Unassigned').map(p => ({ value: p.ID, label: p.Name }));

  $: filteredAssets = assets.filter(asset => {
    // Apply unassigned filter
    if (showUnassignedOnly && asset.CurrentAssignee && asset.CurrentAssignee !== 'Unassigned') {
      return false;
    }
    // Apply search filter
    if (searchTerm.trim()) {
      const term = searchTerm.toLowerCase();
      return (
        (asset.Name || '').toLowerCase().includes(term) ||
        (asset.AssetTypeName || '').toLowerCase().includes(term) ||
        (asset.CurrentAssignee || '').toLowerCase().includes(term) ||
        (asset.Model || '').toLowerCase().includes(term) ||
        (asset.SerialNumber || '').toLowerCase().includes(term)
      );
    }
    return true;
  });

  function handleSearch(term) {
    searchTerm = term;
  }

  function toggleUnassignedOnly() {
    showUnassignedOnly = !showUnassignedOnly;
  }
</script>

<h1 class="title">Asset Assignments</h1>

<Card>
  <div class="mb-4">
    <div class="field is-grouped">
      <div class="control is-expanded">
        <SearchInput placeholder="Search assets..." onSearch={handleSearch} bind:value={searchTerm} />
      </div>
      <div class="control">
        <button
          class="button"
          class:is-info={showUnassignedOnly}
          class:is-outlined={!showUnassignedOnly}
          on:click={toggleUnassignedOnly}
        >
          <span class="icon is-small">
            <i class="fas fa-user-slash"></i>
          </span>
          <span>{showUnassignedOnly ? 'Show All' : 'Unassigned Only'}</span>
        </button>
      </div>
    </div>
  </div>

  {#if loading}
    <Loading />
  {:else}
    <DataTable {columns} data={filteredAssets} emptyMessage="No assets found" />
  {/if}
</Card>

<Modal bind:active={showAssignModal} title="Assign Asset" size="small">
  {#if selectedAsset}
    <div class="notification is-info is-light">
      Assigning: <strong>{selectedAsset.Name}</strong>
      {#if selectedAsset.CurrentAssignee && selectedAsset.CurrentAssignee !== 'Unassigned'}
        <br>Currently assigned to: {selectedAsset.CurrentAssignee}
      {/if}
    </div>
  {/if}
  
  <FormField
    label="Assign to"
    type="select"
    name="person"
    bind:value={form.PersonID}
    options={personOptions}
    required
  />
  <FormField 
    label="Effective Date" 
    type="date" 
    name="effectiveDate" 
    bind:value={form.EffectiveDate}
    required
  />
  <FormField label="Notes" type="textarea" name="notes" bind:value={form.Notes} />
  
  <svelte:fragment slot="footer">
    <Button color="primary" on:click={handleAssign}>Assign</Button>
    <Button on:click={() => showAssignModal = false}>Cancel</Button>
  </svelte:fragment>
</Modal>

<Modal bind:active={showEditModal} title="Edit Assignment" size="small">
  {#if editingAssignment}
    <div class="notification is-info is-light">
      Editing assignment for: <strong>{editingAssignment.AssetName}</strong>
    </div>
  {/if}
  
  <FormField
    label="Assigned to"
    type="select"
    name="editPerson"
    bind:value={editForm.PersonID}
    options={editPersonOptions}
    required
  />
  <FormField 
    label="Assigned Date" 
    type="date" 
    name="editEffectiveDate" 
    bind:value={editForm.EffectiveDate}
    required
  />
  <FormField label="Notes" type="textarea" name="editNotes" bind:value={editForm.Notes} />
  
  <svelte:fragment slot="footer">
    <Button color="primary" on:click={handleEditSave}>Save</Button>
    <Button on:click={() => { showEditModal = false; editingAssignment = null; }}>Cancel</Button>
  </svelte:fragment>
</Modal>

<ConfirmDialog
  bind:active={showUnassignConfirm}
  title="Unassign Asset"
  message={unassignTarget ? `Are you sure you want to unassign "${unassignTarget.Name}" from ${unassignTarget.CurrentAssignee}?` : ''}
  confirmText="Unassign"
  confirmColor="warning"
  onConfirm={handleUnassign}
  onCancel={() => { showUnassignConfirm = false; unassignTarget = null; }}
/>
