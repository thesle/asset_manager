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

  let person = null;
  let attributes = [];
  let allAttributes = [];
  let currentAssets = [];
  let loading = true;
  let showAttributeModal = false;

  let attributeForm = { AttributeID: '', Value: '' };

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    loading = true;
    try {
      const id = params.id;
      const [personResult, attrsResult, allAttrsResult, assetsResult] = await Promise.all([
        api.getPerson(id),
        api.getPersonAttributes(id),
        api.getAttributes(),
        api.getCurrentPersonAssignments(id)
      ]);
      person = personResult;
      attributes = attrsResult || [];
      allAttributes = allAttrsResult || [];
      currentAssets = assetsResult || [];
    } catch (err) {
      notifications.error('Failed to load person');
    } finally {
      loading = false;
    }
  }

  async function handleAddAttribute() {
    try {
      await api.setPersonAttribute(params.id, {
        AttributeID: parseInt(attributeForm.AttributeID),
        Value: attributeForm.Value
      });
      notifications.success('Attribute added');
      showAttributeModal = false;
      attributeForm = { AttributeID: '', Value: '' };
      attributes = await api.getPersonAttributes(params.id);
    } catch (err) {
      notifications.error(err.message);
    }
  }

  $: attributeOptions = allAttributes.map(a => ({ value: a.ID, label: a.Name }));
  $: selectedAttribute = allAttributes.find(a => a.ID === parseInt(attributeForm.AttributeID));

  const assetColumns = [
    { key: 'AssetName', label: 'Asset' },
    { key: 'EffectiveFrom', label: 'Assigned Since', render: (v) => v ? new Date(v).toLocaleDateString() : '' }
  ];
</script>

{#if loading}
  <Loading />
{:else if person}
  <div class="level">
    <div class="level-left">
      <div>
        <h1 class="title">{person.Name}</h1>
        <p class="subtitle">{person.Email || 'No email'}</p>
      </div>
    </div>
    <div class="level-right">
      <a href="#/persons" class="button">
        <span class="icon"><i class="fas fa-arrow-left"></i></span>
        <span>Back to Persons</span>
      </a>
    </div>
  </div>

  <div class="columns">
    <div class="column is-6">
      <Card title="Details">
        <table class="table is-fullwidth">
          <tbody>
            <tr><th>Name</th><td>{person.Name}</td></tr>
            <tr><th>Email</th><td>{person.Email || '-'}</td></tr>
            <tr><th>Phone</th><td>{person.Phone || '-'}</td></tr>
          </tbody>
        </table>
      </Card>
    </div>

    <div class="column is-6">
      <Card title="Current Assets">
        <DataTable columns={assetColumns} data={currentAssets} emptyMessage="No assets assigned" />
      </Card>
    </div>
  </div>

  <Card title="Attributes">
    <svelte:fragment slot="headerIcon">
      <Button size="small" color="primary" on:click={() => showAttributeModal = true}>
        <span class="icon"><i class="fas fa-plus"></i></span>
      </Button>
    </svelte:fragment>
    
    {#if attributes.length === 0}
      <p class="has-text-grey">No attributes</p>
    {:else}
      <table class="table is-fullwidth">
        <thead>
          <tr><th>Attribute</th><th>Value</th></tr>
        </thead>
        <tbody>
          {#each attributes as attr}
            <tr>
              <td>{attr.AttributeName}</td>
              <td>
                {#if attr.DataType === 'boolean'}
                  {attr.Value === 'true' ? 'Yes' : 'No'}
                {:else}
                  {attr.Value}
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    {/if}
  </Card>
{/if}

<Modal bind:active={showAttributeModal} title="Add Attribute" size="small">
  <FormField
    label="Attribute"
    type="select"
    name="attribute"
    bind:value={attributeForm.AttributeID}
    options={attributeOptions}
    required
  />
  {#if selectedAttribute}
    <DynamicField
      label="Value"
      name="value"
      dataType={selectedAttribute.DataType}
      enumOptions={selectedAttribute.EnumOptions}
      bind:value={attributeForm.Value}
    />
  {/if}
  
  <svelte:fragment slot="footer">
    <Button color="primary" on:click={handleAddAttribute}>Add</Button>
    <Button on:click={() => showAttributeModal = false}>Cancel</Button>
  </svelte:fragment>
</Modal>
