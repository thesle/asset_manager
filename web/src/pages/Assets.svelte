<script>
  import { onMount } from 'svelte';
  import { querystring } from 'svelte-spa-router';
  import { api, notifications } from '../stores.js';
  import Card from '../../../shared/components/Card.svelte';
  import DataTable from '../../../shared/components/DataTable.svelte';
  import Button from '../../../shared/components/Button.svelte';
  import Modal from '../../../shared/components/Modal.svelte';
  import FormField from '../../../shared/components/FormField.svelte';
  import SearchInput from '../../../shared/components/SearchInput.svelte';
  import ConfirmDialog from '../../../shared/components/ConfirmDialog.svelte';
  import CustomFields from '../../../shared/components/CustomFields.svelte';

  let assets = [];
  let assetTypes = [];
  let properties = [];
  let loading = true;
  let showModal = false;
  let showDeleteConfirm = false;
  let editingAsset = null;
  let deleteTarget = null;
  let saving = false;
  let initialEditHandled = false;

  let form = {
    AssetTypeID: '',
    Name: '',
    Model: '',
    SerialNumber: '',
    OrderNo: '',
    LicenseNumber: '',
    Notes: ''
  };
  let customFieldValues = {};

  const columns = [
    { key: 'Name', label: 'Name', sortable: true },
    { key: 'AssetTypeName', label: 'Type', sortable: true },
    { key: 'Model', label: 'Model', sortable: true },
    { key: 'SerialNumber', label: 'Serial Number', sortable: true },
    { key: 'CurrentAssignee', label: 'Assigned To', sortable: true },
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
    
    // Check for edit query parameter
    checkEditParam();
    
    // Event delegation for table buttons
    document.addEventListener('click', handleTableClick);
    return () => document.removeEventListener('click', handleTableClick);
  });
  
  // React to querystring changes
  $: if ($querystring && !initialEditHandled) {
    checkEditParam();
  }
  
  function checkEditParam() {
    const params = new URLSearchParams($querystring);
    const editId = params.get('edit');
    if (editId && assets.length > 0) {
      const asset = assets.find(a => a.ID === parseInt(editId));
      if (asset) {
        initialEditHandled = true;
        openEdit(asset);
        // Clear the query param from URL
        window.history.replaceState(null, '', '#/assets');
      }
    }
  }

  async function loadData() {
    loading = true;
    try {
      const [assetsResult, typesResult, propsResult] = await Promise.all([
        api.getAssetsWithAssignments(),
        api.getAssetTypes(),
        api.getProperties()
      ]);
      assets = assetsResult || [];
      assetTypes = typesResult || [];
      properties = propsResult || [];
    } catch (err) {
      notifications.error('Failed to load assets');
    } finally {
      loading = false;
    }
  }

  function handleTableClick(e) {
    const editBtn = e.target.closest('.edit-btn');
    const deleteBtn = e.target.closest('.delete-btn');
    
    if (editBtn) {
      const id = parseInt(editBtn.dataset.id);
      openEdit(assets.find(a => a.ID === id));
    } else if (deleteBtn) {
      const id = parseInt(deleteBtn.dataset.id);
      confirmDelete(assets.find(a => a.ID === id));
    }
  }

  function openNew() {
    editingAsset = null;
    form = {
      AssetTypeID: '',
      Name: '',
      Model: '',
      SerialNumber: '',
      OrderNo: '',
      LicenseNumber: '',
      Notes: ''
    };
    customFieldValues = {};
    showModal = true;
  }

  async function openEdit(asset) {
    editingAsset = asset;
    form = {
      AssetTypeID: asset.AssetTypeID,
      Name: asset.Name,
      Model: asset.Model || '',
      SerialNumber: asset.SerialNumber || '',
      OrderNo: asset.OrderNo || '',
      LicenseNumber: asset.LicenseNumber || '',
      Notes: asset.Notes || ''
    };
    
    // Load existing property values
    customFieldValues = {};
    try {
      const assetProps = await api.getAssetProperties(asset.ID);
      if (assetProps) {
        for (const prop of assetProps) {
          customFieldValues[prop.PropertyID] = prop.Value;
        }
      }
    } catch (err) {
      // Ignore - just show empty custom fields
    }
    showModal = true;
  }

  function confirmDelete(asset) {
    deleteTarget = asset;
    showDeleteConfirm = true;
  }

  async function handleSave() {
    saving = true;
    try {
      const data = { ...form, AssetTypeID: parseInt(form.AssetTypeID) };
      let assetId;
      
      if (editingAsset) {
        await api.updateAsset(editingAsset.ID, data);
        assetId = editingAsset.ID;
      } else {
        const created = await api.createAsset(data);
        assetId = created.ID;
      }
      
      // Save custom field values
      for (const [propId, value] of Object.entries(customFieldValues)) {
        if (value !== undefined && value !== '') {
          await api.setAssetProperty(assetId, { PropertyID: parseInt(propId), Value: String(value) });
        }
      }
      
      notifications.success(editingAsset ? 'Asset updated' : 'Asset created');
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
      await api.deleteAsset(deleteTarget.ID);
      notifications.success('Asset deleted');
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
      const result = await api.searchAssets(term);
      assets = result || [];
    } catch (err) {
      notifications.error('Search failed');
    } finally {
      loading = false;
    }
  }

  $: assetTypeOptions = assetTypes.map(t => ({ value: t.ID, label: t.Name }));
</script>

<div class="level">
  <div class="level-left">
    <h1 class="title">Assets</h1>
  </div>
  <div class="level-right">
    <Button color="primary" on:click={openNew}>
      <span class="icon"><i class="fas fa-plus"></i></span>
      <span>New Asset</span>
    </Button>
  </div>
</div>

<Card>
  <div class="mb-4">
    <SearchInput placeholder="Search assets..." onSearch={handleSearch} />
  </div>
  
  <DataTable {columns} data={assets} {loading} emptyMessage="No assets found" />
</Card>

<Modal bind:active={showModal} title={editingAsset ? 'Edit Asset' : 'New Asset'} size="wide">
  <form on:submit|preventDefault={handleSave}>
    <div class="columns">
      <div class="column">
        <h6 class="title is-6 mb-3">Asset Details</h6>
        <FormField
          label="Asset Type"
          type="select"
          name="assetType"
          bind:value={form.AssetTypeID}
          options={assetTypeOptions}
          required
        />
        <FormField label="Name" name="name" bind:value={form.Name} required />
        <FormField label="Model" name="model" bind:value={form.Model} />
        <FormField label="Serial Number" name="serialNumber" bind:value={form.SerialNumber} />
        <FormField label="Order No" name="orderNo" bind:value={form.OrderNo} />
        <FormField label="License Number" name="licenseNumber" bind:value={form.LicenseNumber} />
        <FormField label="Notes" type="textarea" name="notes" bind:value={form.Notes} />
      </div>
      <div class="column">
        <h6 class="title is-6 mb-3">Custom Properties</h6>
        <CustomFields 
          definitions={properties} 
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
  title="Delete Asset"
  message="Are you sure you want to delete this asset? This action cannot be undone."
  onConfirm={handleDelete}
/>
