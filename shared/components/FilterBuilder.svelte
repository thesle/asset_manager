<script>
  import { onMount } from 'svelte';
  import { api } from '../../web/src/stores.js';
  
  export let entityType = 'asset'; // 'asset' or 'person'
  export let filters = [];
  export let properties = [];
  export let attributes = [];

  let assetTypes = [];

  $: availableFields = getAvailableFields(entityType, properties, attributes);
  
  onMount(async () => {
    try {
      assetTypes = await api.getAssetTypes() || [];
    } catch (err) {
      console.error('Failed to load asset types:', err);
    }
  });

  function getAvailableFields(type, props, attrs) {
    const baseFields = type === 'asset' 
      ? [
          { value: 'Name', label: 'Asset Name', type: 'text' },
          { value: 'AssetTypeName', label: 'Asset Type', type: 'text' },
          { value: 'Model', label: 'Model', type: 'text' },
          { value: 'SerialNumber', label: 'Serial Number', type: 'text' },
          { value: 'OrderNo', label: 'Order No', type: 'text' },
          { value: 'LicenseNumber', label: 'License Number', type: 'text' },
          { value: 'Notes', label: 'Notes', type: 'text' },
          { value: 'PurchasedAt', label: 'Purchased At', type: 'date' },
          { value: 'CurrentAssignee', label: 'Current Assignee', type: 'text' },
        ]
      : [
          { value: 'PersonName', label: 'Name', type: 'text' },
          { value: 'PersonEmail', label: 'Email', type: 'text' },
          { value: 'PersonPhone', label: 'Phone', type: 'text' },
        ];

    // Add properties or attributes
    const customFields = type === 'asset'
      ? (props || []).map(p => ({
          value: `prop_${p.Name}`,
          label: `Property: ${p.Name}`,
          type: p.DataType || 'text'
        }))
      : (attrs || []).map(a => ({
          value: `attr_${a.Name}`,
          label: `Attribute: ${a.Name}`,
          type: a.DataType || 'text'
        }));

    return [...baseFields, ...customFields];
  }

  function getOperatorsForType(fieldType) {
    switch (fieldType) {
      case 'number':
        return ['=', '!=', '>', '<', '>=', '<=', 'IS NULL', 'IS NOT NULL'];
      case 'date':
        return ['=', '!=', '>', '<', '>=', '<=', 'IS NULL', 'IS NOT NULL'];
      case 'boolean':
        return ['=', 'IS NULL', 'IS NOT NULL'];
      case 'text':
      default:
        return ['=', '!=', 'LIKE', 'NOT LIKE', 'IS NULL', 'IS NOT NULL'];
    }
  }

  function addFilter() {
    filters = [...filters, {
      field: availableFields[0]?.value || '',
      fieldType: availableFields[0]?.type || 'text',
      operator: '=',
      value: '',
      logicOperator: 'AND'
    }];
  }

  function removeFilter(index) {
    filters = filters.filter((_, i) => i !== index);
  }

  function onFieldChange(index) {
    const field = availableFields.find(f => f.value === filters[index].field);
    if (field) {
      filters[index].fieldType = field.type;
      const operators = getOperatorsForType(field.type);
      if (!operators.includes(filters[index].operator)) {
        filters[index].operator = operators[0];
      }
    }
  }

  function getInputType(fieldType) {
    switch (fieldType) {
      case 'number':
        return 'number';
      case 'date':
        return 'date';
      case 'boolean':
        return 'checkbox';
      default:
        return 'text';
    }
  }
</script>

<div class="filter-builder">
  <div class="level mb-3">
    <div class="level-left">
      <div class="level-item">
        <h5 class="title is-5">Filters</h5>
      </div>
    </div>
    <div class="level-right">
      <div class="level-item">
        <button class="button is-small is-primary" on:click={addFilter}>
          <span class="icon"><i class="fas fa-plus"></i></span>
          <span>Add Filter</span>
        </button>
      </div>
    </div>
  </div>

  {#if filters.length === 0}
    <div class="notification is-info is-light">
      <p>No filters added. Click "Add Filter" to start building your query.</p>
    </div>
  {:else}
    <div class="filters-list">
      {#each filters as filter, index}
        <div class="filter-row box p-3 mb-2">
          <div class="columns is-vcentered is-mobile">
            <div class="column is-narrow">
              <span class="tag is-light">{index + 1}</span>
            </div>
            
            <div class="column">
              <div class="select is-fullwidth is-small">
                <select bind:value={filter.field} on:change={() => onFieldChange(index)}>
                  {#each availableFields as field}
                    <option value={field.value}>{field.label}</option>
                  {/each}
                </select>
              </div>
            </div>

            <div class="column is-narrow">
              <div class="select is-small">
                <select bind:value={filter.operator}>
                  {#each getOperatorsForType(filter.fieldType) as op}
                    <option value={op}>{op}</option>
                  {/each}
                </select>
              </div>
            </div>

            {#if !['IS NULL', 'IS NOT NULL'].includes(filter.operator)}
              <div class="column">
                {#if filter.field === 'AssetTypeName'}
                  <div class="select is-small is-fullwidth">
                    <select bind:value={filter.value}>
                      <option value="">Select Asset Type</option>
                      {#each assetTypes as assetType}
                        <option value={assetType.Name}>{assetType.Name}</option>
                      {/each}
                    </select>
                  </div>
                {:else if filter.fieldType === 'boolean'}
                  <div class="select is-small is-fullwidth">
                    <select bind:value={filter.value}>
                      <option value="true">True</option>
                      <option value="false">False</option>
                    </select>
                  </div>
                {:else if filter.fieldType === 'number'}
                  <input
                    class="input is-small"
                    type="number"
                    bind:value={filter.value}
                    placeholder="Value"
                  />
                {:else if filter.fieldType === 'date'}
                  <input
                    class="input is-small"
                    type="date"
                    bind:value={filter.value}
                    placeholder="Value"
                  />
                {:else}
                  <input
                    class="input is-small"
                    type="text"
                    bind:value={filter.value}
                    placeholder="Value"
                  />
                {/if}
              </div>
            {/if}

            <div class="column is-narrow">
              <button 
                class="button is-small is-danger is-light" 
                on:click={() => removeFilter(index)}
                title="Remove filter"
              >
                <span class="icon"><i class="fas fa-times"></i></span>
              </button>
            </div>
          </div>

          {#if index < filters.length - 1}
            <div class="logic-operator has-text-centered mt-2">
              <div class="buttons is-centered">
                <button 
                  class="button is-small"
                  class:is-primary={filter.logicOperator === 'AND'}
                  class:is-outlined={filter.logicOperator !== 'AND'}
                  on:click={() => filter.logicOperator = 'AND'}
                >
                  AND
                </button>
                <button 
                  class="button is-small"
                  class:is-primary={filter.logicOperator === 'OR'}
                  class:is-outlined={filter.logicOperator !== 'OR'}
                  on:click={() => filter.logicOperator = 'OR'}
                >
                  OR
                </button>
              </div>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .filter-builder {
    margin-bottom: 1rem;
  }

  .filter-row {
    background-color: #fafafa;
  }

  .filters-list {
    max-height: 400px;
    overflow-y: auto;
  }
</style>
