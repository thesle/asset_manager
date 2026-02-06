<script>
  import { onMount } from 'svelte';
  import { api, notifications } from '../stores.js';
  import Card from '../../../shared/components/Card.svelte';
  import Button from '../../../shared/components/Button.svelte';
  import DataTable from '../../../shared/components/DataTable.svelte';
  import Modal from '../../../shared/components/Modal.svelte';
  import FormField from '../../../shared/components/FormField.svelte';
  import DynamicField from '../../../shared/components/DynamicField.svelte';
  import Loading from '../../../shared/components/Loading.svelte';

  export let params = {};

  let asset = null;
  let properties = [];
  let allProperties = [];
  let assignments = [];
  let persons = [];
  let loading = true;
  let showPropertyModal = false;
  let showAssignModal = false;

  let propertyForm = { PropertyID: '', Value: '' };
  let assignForm = { PersonID: '', Notes: '' };

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    loading = true;
    try {
      const id = params.id;
      const [assetResult, propsResult, assignResult, allPropsResult, personsResult] = await Promise.all([
        api.getAsset(id),
        api.getAssetProperties(id),
        api.getAssetAssignments(id),
        api.getProperties(),
        api.getPersons()
      ]);
      asset = assetResult;
      properties = propsResult || [];
      assignments = assignResult || [];
      allProperties = allPropsResult || [];
      persons = personsResult || [];
    } catch (err) {
      notifications.error('Failed to load asset');
    } finally {
      loading = false;
    }
  }

  async function handleAddProperty() {
    try {
      await api.setAssetProperty(params.id, {
        PropertyID: parseInt(propertyForm.PropertyID),
        Value: propertyForm.Value
      });
      notifications.success('Property added');
      showPropertyModal = false;
      propertyForm = { PropertyID: '', Value: '' };
      properties = await api.getAssetProperties(params.id);
    } catch (err) {
      notifications.error(err.message);
    }
  }

  async function handleAssign() {
    try {
      await api.assignAsset(parseInt(params.id), parseInt(assignForm.PersonID), assignForm.Notes);
      notifications.success('Asset assigned');
      showAssignModal = false;
      assignForm = { PersonID: '', Notes: '' };
      assignments = await api.getAssetAssignments(params.id);
    } catch (err) {
      notifications.error(err.message);
    }
  }

  async function handleUnassign() {
    try {
      await api.unassignAsset(params.id);
      notifications.success('Asset unassigned');
      assignments = await api.getAssetAssignments(params.id);
    } catch (err) {
      notifications.error(err.message);
    }
  }

  $: propertyOptions = allProperties.map(p => ({ value: p.ID, label: p.Name }));
  $: personOptions = persons.filter(p => p.Name !== 'Unassigned').map(p => ({ value: p.ID, label: p.Name }));
  $: selectedProperty = allProperties.find(p => p.ID === parseInt(propertyForm.PropertyID));
  $: currentAssignment = assignments.find(a => !a.EffectiveTo || new Date(a.EffectiveTo) > new Date());

  const assignmentColumns = [
    { key: 'PersonName', label: 'Person' },
    { key: 'EffectiveFrom', label: 'From', render: (v) => v ? new Date(v).toLocaleDateString() : '' },
    { key: 'EffectiveTo', label: 'To', render: (v) => v ? new Date(v).toLocaleDateString() : 'Current' },
    { key: 'Notes', label: 'Notes' }
  ];
</script>

{#if loading}
  <Loading />
{:else if asset}
  <div class="level">
    <div class="level-left">
      <div>
        <h1 class="title">{asset.Name}</h1>
        <p class="subtitle">{asset.AssetTypeName}</p>
      </div>
    </div>
    <div class="level-right">
      <a href="#/assets" class="button">
        <span class="icon"><i class="fas fa-arrow-left"></i></span>
        <span>Back to Assets</span>
      </a>
    </div>
  </div>

  <div class="columns">
    <div class="column is-6">
      <Card title="Details">
        <table class="table is-fullwidth">
          <tbody>
            <tr><th>Model</th><td>{asset.Model || '-'}</td></tr>
            <tr><th>Serial Number</th><td>{asset.SerialNumber || '-'}</td></tr>
            <tr><th>Order No</th><td>{asset.OrderNo || '-'}</td></tr>
            <tr><th>License Number</th><td>{asset.LicenseNumber || '-'}</td></tr>
            <tr><th>Notes</th><td>{asset.Notes || '-'}</td></tr>
          </tbody>
        </table>
      </Card>
    </div>

    <div class="column is-6">
      <Card title="Current Assignment">
        {#if currentAssignment && currentAssignment.PersonName !== 'Unassigned'}
          <div class="content">
            <p><strong>Assigned to:</strong> {currentAssignment.PersonName}</p>
            <p><strong>Since:</strong> {new Date(currentAssignment.EffectiveFrom).toLocaleDateString()}</p>
            {#if currentAssignment.Notes}
              <p><strong>Notes:</strong> {currentAssignment.Notes}</p>
            {/if}
          </div>
          <div class="buttons">
            <Button color="warning" on:click={() => showAssignModal = true}>Reassign</Button>
            <Button color="danger" outlined on:click={handleUnassign}>Unassign</Button>
          </div>
        {:else}
          <p class="has-text-grey">Not currently assigned</p>
          <Button color="primary" on:click={() => showAssignModal = true}>Assign</Button>
        {/if}
      </Card>
    </div>
  </div>

  <div class="columns">
    <div class="column is-6">
      <Card title="Custom Properties">
        <svelte:fragment slot="headerIcon">
          <Button size="small" color="primary" on:click={() => showPropertyModal = true}>
            <span class="icon"><i class="fas fa-plus"></i></span>
          </Button>
        </svelte:fragment>
        
        {#if properties.length === 0}
          <p class="has-text-grey">No custom properties</p>
        {:else}
          <table class="table is-fullwidth">
            <tbody>
              {#each properties as prop}
                <tr>
                  <th>{prop.PropertyName}</th>
                  <td>{prop.Value}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        {/if}
      </Card>
    </div>

    <div class="column is-6">
      <Card title="Assignment History">
        <DataTable columns={assignmentColumns} data={assignments} emptyMessage="No assignment history" />
      </Card>
    </div>
  </div>
{/if}

<Modal bind:active={showPropertyModal} title="Add Property" size="small">
  <FormField
    label="Property"
    type="select"
    name="property"
    bind:value={propertyForm.PropertyID}
    options={propertyOptions}
    required
  />
  {#if selectedProperty}
    <DynamicField
      label="Value"
      name="value"
      dataType={selectedProperty.DataType}
      enumOptions={selectedProperty.EnumOptions}
      bind:value={propertyForm.Value}
    />
  {/if}
  
  <svelte:fragment slot="footer">
    <Button color="primary" on:click={handleAddProperty}>Add</Button>
    <Button on:click={() => showPropertyModal = false}>Cancel</Button>
  </svelte:fragment>
</Modal>

<Modal bind:active={showAssignModal} title="Assign Asset" size="small">
  <FormField
    label="Person"
    type="select"
    name="person"
    bind:value={assignForm.PersonID}
    options={personOptions}
    required
  />
  <FormField label="Notes" type="textarea" name="notes" bind:value={assignForm.Notes} />
  
  <svelte:fragment slot="footer">
    <Button color="primary" on:click={handleAssign}>Assign</Button>
    <Button on:click={() => showAssignModal = false}>Cancel</Button>
  </svelte:fragment>
</Modal>
