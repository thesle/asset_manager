<script>
  import { onMount } from 'svelte';
  import { api, notifications } from '../../stores.js';
  import Button from '../../../../shared/components/Button.svelte';
  import Loading from '../../../../shared/components/Loading.svelte';
  import FormField from '../../../../shared/components/FormField.svelte';
  import { exportToCSV } from '../../../../shared/utils/csvExport.js';

  let assetTypes = [];
  let attributes = [];
  let properties = [];
  let persons = [];
  let loading = true;
  let searching = false;
  let selectedAssetTypeId = '';
  let expandedPersonId = null;
  let personAssets = {};
  let loadingAssets = {};

  onMount(async () => {
    try {
      assetTypes = await api.getAssetTypes() || [];
      attributes = await api.getAttributes() || [];
      properties = await api.getProperties() || [];
    } catch (err) {
      notifications.error('Failed to load initial data');
    } finally {
      loading = false;
    }
  });

  async function runReport() {
    if (!selectedAssetTypeId) {
      notifications.error('Please select an asset type');
      return;
    }

    searching = true;
    expandedPersonId = null;
    personAssets = {};
    try {
      const rawPersons = await api.getMultipleAssetsReport(selectedAssetTypeId) || [];
      
      // Load attributes for each person
      persons = await Promise.all(rawPersons.map(async (person) => {
        try {
          const personAttrs = await api.getPersonAttributes(person.id);
          return { ...person, attributes: personAttrs || [] };
        } catch {
          return { ...person, attributes: [] };
        }
      }));
      
      if (persons.length === 0) {
        notifications.info('No persons found with multiple assets of this type');
      }
    } catch (err) {
      notifications.error('Failed to run report: ' + err.message);
      persons = [];
    } finally {
      searching = false;
    }
  }

  async function toggleExpand(personId) {
    if (expandedPersonId === personId) {
      expandedPersonId = null;
      return;
    }

    expandedPersonId = personId;

    // Load current assets if not already loaded
    if (!personAssets[personId]) {
      loadingAssets[personId] = true;
      loadingAssets = loadingAssets;
      try {
        const assignments = await api.getCurrentPersonAssignments(personId);
        const assets = assignments || [];

        // Load properties for each assigned asset
        const assetsWithProps = await Promise.all(
          assets.map(async (assignment) => {
            try {
              const assetProps = await api.getAssetProperties(assignment.AssetID);
              return { ...assignment, properties: assetProps || [] };
            } catch {
              return { ...assignment, properties: [] };
            }
          })
        );

        personAssets[personId] = assetsWithProps;
        personAssets = personAssets;
      } catch {
        personAssets[personId] = [];
        personAssets = personAssets;
      } finally {
        loadingAssets[personId] = false;
        loadingAssets = loadingAssets;
      }
    }
  }

  function getPropertyValue(assignment, propId) {
    const prop = assignment.properties?.find(p => p.PropertyID === propId);
    return prop?.Value || '-';
  }

  function formatDate(dateStr) {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleDateString();
  }

  function handleExport() {
    if (persons.length === 0) {
      notifications.error('No data to export');
      return;
    }

    const selectedType = assetTypes.find(t => t.ID === parseInt(selectedAssetTypeId));
    const filename = `multiple-${selectedType?.Name || 'assets'}-report-${new Date().toISOString().split('T')[0]}.csv`;
    exportToCSV(persons, filename);
    notifications.success('Report exported');
  }

  function getAttributeValue(person, attrId) {
    const attr = person.attributes?.find((a) => a.AttributeID === attrId);
    if (!attr) return "-";
    return attr.Value || "-";
  }

  function isAttributeBoolean(attrId) {
    const attrDef = attributes.find((a) => a.ID === attrId);
    return attrDef?.DataType === "boolean";
  }

  function getAttributeBoolValue(person, attrId) {
    const attr = person.attributes?.find((a) => a.AttributeID === attrId);
    if (!attr) return null;
    return attr.Value === "true";
  }

  $: selectedAssetType = assetTypes.find(t => t.ID === parseInt(selectedAssetTypeId));
  $: assetTypeOptions = assetTypes.map(t => ({ value: t.ID, label: t.Name }));
</script>

<style>
  .accordion-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .accordion-item {
    border: 1px solid #dbdbdb;
    border-radius: 4px;
    overflow: hidden;
  }

  .accordion-item.is-expanded {
    border-color: #3273dc;
  }

  .accordion-header {
    display: flex;
    align-items: center;
    padding: 0.75rem 1rem;
    background-color: #f5f5f5;
    gap: 0.75rem;
  }

  .accordion-header:hover {
    background-color: #eeeeee;
  }

  .is-clickable {
    cursor: pointer;
  }

  .accordion-icon {
    flex-shrink: 0;
  }

  .accordion-title {
    flex-shrink: 0;
    min-width: 200px;
  }

  .accordion-details {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
    align-items: center;
  }

  .accordion-content {
    padding: 0;
    background-color: #fafafa;
  }

  .accordion-content .box {
    margin: 1rem;
    background-color: white;
  }
</style>

<div class="container">
  <section class="section">
    <h1 class="title">Multiple Assets Report</h1>
    <p class="subtitle">Find persons holding multiple assets of the same type</p>

    {#if loading}
      <Loading />
    {:else}
      <div class="box">
        <div class="columns">
          <div class="column is-half">
            <FormField 
              label="Asset Type" 
              type="select" 
              name="assetType" 
              bind:value={selectedAssetTypeId} 
              options={assetTypeOptions}
              placeholder="Select Asset Type"
              required 
            />
          </div>
          <div class="column is-half">
            <div class="field">
              <label class="label">&nbsp;</label>
              <div class="control">
                <Button color="primary" on:click={runReport} disabled={searching || !selectedAssetTypeId}>
                  {#if searching}
                    <span class="icon"><i class="fas fa-spinner fa-pulse"></i></span>
                  {:else}
                    <span class="icon"><i class="fas fa-search"></i></span>
                  {/if}
                  <span>Run Report</span>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {#if persons.length > 0}
        <div class="box">
          <div class="level mb-4">
            <div class="level-left">
              <div class="level-item">
                <p class="subtitle is-5">
                  Found {persons.length} person{persons.length !== 1 ? 's' : ''} with multiple {selectedAssetType?.Name || 'assets'}
                </p>
              </div>
            </div>
            <div class="level-right">
              <div class="level-item">
                <Button color="info" outlined on:click={handleExport}>
                  <span class="icon"><i class="fas fa-download"></i></span>
                  <span>Export CSV</span>
                </Button>
              </div>
            </div>
          </div>

          <div class="accordion-list">
            {#each persons as person}
              <div class="accordion-item" class:is-expanded={expandedPersonId === person.id}>
                <div
                  class="accordion-header is-clickable"
                  on:click={() => toggleExpand(person.id)}
                  on:keydown={(e) => e.key === 'Enter' && toggleExpand(person.id)}
                  role="button"
                  tabindex="0"
                >
                  <div class="accordion-icon">
                    <span class="icon">
                      <i
                        class="fas"
                        class:fa-chevron-down={expandedPersonId === person.id}
                        class:fa-chevron-right={expandedPersonId !== person.id}
                      ></i>
                    </span>
                  </div>
                  <div class="accordion-title">
                    <a
                      href="#/persons?edit={person.id}"
                      on:click|stopPropagation
                      class="has-text-link has-text-weight-semibold"
                    >
                      {person.name}
                    </a>
                  </div>
                  <div class="accordion-details">
                    <span class="tag is-warning is-light mr-2">
                      <span class="icon is-small"><i class="fas fa-laptop"></i></span>
                      <span>{person.asset_count} {selectedAssetType?.Name || 'assets'}</span>
                    </span>
                    <span class="tag is-light mr-2">
                      <span class="icon is-small"><i class="fas fa-envelope"></i></span>
                      <span>{person.email || 'No email'}</span>
                    </span>
                    <span class="tag is-light mr-2">
                      <span class="icon is-small"><i class="fas fa-phone"></i></span>
                      <span>{person.phone || 'No phone'}</span>
                    </span>
                    {#each attributes as attr}
                      <span class="tag is-info is-light mr-1">
                        {attr.Name}:
                        {#if isAttributeBoolean(attr.ID)}
                          {#if getAttributeBoolValue(person, attr.ID) === true}
                            <span class="icon is-small has-text-success ml-1"><i class="fas fa-check"></i></span>
                          {:else if getAttributeBoolValue(person, attr.ID) === false}
                            <span class="icon is-small has-text-danger ml-1"><i class="fas fa-times"></i></span>
                          {:else}
                            -
                          {/if}
                        {:else}
                          {getAttributeValue(person, attr.ID)}
                        {/if}
                      </span>
                    {/each}
                  </div>
                </div>

                {#if expandedPersonId === person.id}
                  <div class="accordion-content">
                    <div class="box">
                      <h6 class="title is-6 mb-3">Assigned {selectedAssetType?.Name || 'Assets'}</h6>
                      {#if loadingAssets[person.id]}
                        <p class="has-text-grey"><i class="fas fa-spinner fa-spin"></i> Loading...</p>
                      {:else if !personAssets[person.id] || personAssets[person.id].length === 0}
                        <p class="has-text-grey is-italic">No assets currently assigned</p>
                      {:else}
                        <table class="table is-fullwidth is-striped">
                          <thead>
                            <tr>
                              <th>Asset</th>
                              <th>Type</th>
                              <th>Model</th>
                              <th>Serial Number</th>
                              <th>Assigned Since</th>
                              {#each properties as prop}
                                <th>{prop.Name}</th>
                              {/each}
                            </tr>
                          </thead>
                          <tbody>
                            {#each personAssets[person.id] as assignment}
                              <tr>
                                <td>
                                  <a href="#/assets?edit={assignment.AssetID}" class="has-text-link">
                                    {assignment.AssetName}
                                  </a>
                                </td>
                                <td>{assignment.AssetTypeName || '-'}</td>
                                <td>{assignment.AssetModel || '-'}</td>
                                <td>{assignment.AssetSerialNumber || '-'}</td>
                                <td>{formatDate(assignment.EffectiveFrom)}</td>
                                {#each properties as prop}
                                  <td>{getPropertyValue(assignment, prop.ID)}</td>
                                {/each}
                              </tr>
                            {/each}
                          </tbody>
                        </table>
                      {/if}
                    </div>
                  </div>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/if}
  </section>
</div>
