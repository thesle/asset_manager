<script>
  import { onMount } from 'svelte';
  import { api, notifications } from '../stores.js';
  import Card from '../../../shared/components/Card.svelte';
  import DataTable from '../../../shared/components/DataTable.svelte';
  import Button from '../../../shared/components/Button.svelte';
  import Modal from '../../../shared/components/Modal.svelte';
  import FormField from '../../../shared/components/FormField.svelte';
  import Loading from '../../../shared/components/Loading.svelte';

  let assets = [];
  let persons = [];
  let loading = true;
  let showAssignModal = false;

  let form = { AssetID: '', PersonID: '', Notes: '', EffectiveDate: '' };
  
  // Get today's date in YYYY-MM-DD format for the date input
  function getTodayDate() {
    const today = new Date();
    return today.toISOString().split('T')[0];
  }

  const columns = [
    { key: 'Name', label: 'Asset', sortable: true },
    { key: 'AssetTypeName', label: 'Type', sortable: true },
    { key: 'CurrentAssignee', label: 'Assigned To', sortable: true },
    { 
      key: 'actions', 
      label: 'Actions',
      render: (_, row) => `
        <div class="buttons are-small">
          <button class="button is-primary is-outlined assign-btn" data-id="${row.ID}">
            <span class="icon"><i class="fas fa-user-plus"></i></span>
            <span>${row.CurrentAssignee && row.CurrentAssignee !== 'Unassigned' ? 'Reassign' : 'Assign'}</span>
          </button>
          ${row.CurrentAssignee && row.CurrentAssignee !== 'Unassigned' ? `
            <button class="button is-warning is-outlined unassign-btn" data-id="${row.ID}">
              <span class="icon"><i class="fas fa-user-minus"></i></span>
              <span>Unassign</span>
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
    
    if (assignBtn) {
      const id = parseInt(assignBtn.dataset.id);
      openAssign(assets.find(a => a.ID === id));
    } else if (unassignBtn) {
      const id = parseInt(unassignBtn.dataset.id);
      handleUnassign(id);
    }
  }

  function openAssign(asset) {
    form = { AssetID: asset.ID, PersonID: '', Notes: '', EffectiveDate: getTodayDate() };
    showAssignModal = true;
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

  async function handleUnassign(assetId) {
    try {
      await api.unassignAsset(assetId);
      notifications.success('Asset unassigned');
      await loadData();
    } catch (err) {
      notifications.error(err.message);
    }
  }

  $: personOptions = persons.filter(p => p.Name !== 'Unassigned').map(p => ({ value: p.ID, label: p.Name }));
  $: selectedAsset = assets.find(a => a.ID === form.AssetID);
</script>

<h1 class="title">Asset Assignments</h1>

<Card>
  {#if loading}
    <Loading />
  {:else}
    <DataTable {columns} data={assets} emptyMessage="No assets found" />
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
