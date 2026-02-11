<script>
  import { onMount } from 'svelte';
  import { api, notifications } from '../../stores.js';
  import Card from '../../../../shared/components/Card.svelte';
  import Loading from '../../../../shared/components/Loading.svelte';
  import FilterBuilder from '../../../../shared/components/FilterBuilder.svelte';
  import Button from '../../../../shared/components/Button.svelte';
  import { exportToCSV, generateFilename } from '../../../../shared/utils/csvExport.js';

  let entityType = 'asset';
  let filters = [];
  let properties = [];
  let attributes = [];
  let results = [];
  let loading = false;
  let searching = false;
  let hasSearched = false;

  onMount(async () => {
    await loadMetadata();
  });

  async function loadMetadata() {
    loading = true;
    try {
      const [propsResult, attrsResult] = await Promise.all([
        api.getProperties(),
        api.getAttributes()
      ]);
      properties = propsResult || [];
      attributes = attrsResult || [];
    } catch (err) {
      notifications.error('Failed to load metadata: ' + err.message);
    } finally {
      loading = false;
    }
  }

  function onEntityTypeChange() {
    // Clear filters when switching entity type
    filters = [];
    results = [];
    hasSearched = false;
  }

  async function runReport() {
    if (filters.length === 0) {
      notifications.warning('Please add at least one filter');
      return;
    }

    // Validate filters
    for (const filter of filters) {
      if (!['IS NULL', 'IS NOT NULL'].includes(filter.operator) && !filter.value) {
        notifications.warning('Please fill in all filter values');
        return;
      }
    }

    searching = true;
    hasSearched = false;
    try {
      // Convert filters to API format
      const apiFilters = filters.map(f => ({
        Field: f.field,
        Operator: f.operator,
        Value: f.value,
        LogicOperator: f.logicOperator || 'AND'
      }));

      const response = await api.executeCustomReport({
        EntityType: entityType,
        Filters: apiFilters
      });

      results = response || [];
      hasSearched = true;
      
      if (results.length === 0) {
        notifications.info('No results found');
      } else {
        notifications.success(`Found ${results.length} result${results.length !== 1 ? 's' : ''}`);
      }
    } catch (err) {
      notifications.error('Search failed: ' + err.message);
      results = [];
    } finally {
      searching = false;
    }
  }

  function clearFilters() {
    filters = [];
    results = [];
    hasSearched = false;
  }

  function exportResults() {
    if (results.length === 0) {
      notifications.warning('No results to export');
      return;
    }

    const filename = generateFilename(`custom_report_${entityType}`, 'csv');
    exportToCSV(results, filename);
    notifications.success('Report exported successfully');
  }

  function formatValue(value) {
    if (value === null || value === undefined) {
      return '-';
    }
    if (typeof value === 'boolean') {
      return value ? 'Yes' : 'No';
    }
    if (value instanceof Date) {
      return value.toLocaleDateString();
    }
    if (typeof value === 'object') {
      return JSON.stringify(value);
    }
    return String(value);
  }

  function getDisplayColumns(data) {
    if (!data || data.length === 0) return [];
    
    // Get all unique keys from all results
    const allKeys = new Set();
    data.forEach(row => {
      Object.keys(row).forEach(key => {
        // Skip internal fields
        if (!key.startsWith('_') && key !== 'deleted_at') {
          allKeys.add(key);
        }
      });
    });
    
    const keys = Array.from(allKeys);
    
    // Sort with base fields first, then properties/attributes
    return keys.sort((a, b) => {
      const aIsProp = a.startsWith('prop_') || a.startsWith('attr_');
      const bIsProp = b.startsWith('prop_') || b.startsWith('attr_');
      
      // Base fields come before properties/attributes
      if (!aIsProp && bIsProp) return -1;
      if (aIsProp && !bIsProp) return 1;
      
      // Within same category, sort alphabetically
      return a.localeCompare(b);
    });
  }

  $: displayColumns = getDisplayColumns(results);
  $: columnHeaders = displayColumns.map(col => {
    // Format column names for display
    return col
      .replace(/_/g, ' ')
      .replace(/\b\w/g, l => l.toUpperCase())
      .replace(/^Prop /, 'Property: ')
      .replace(/^Attr /, 'Attribute: ');
  });
</script>

<h1 class="title">Custom Report</h1>

{#if loading}
  <Loading />
{:else}
  <Card>
    <div class="content">
      <h5 class="title is-5 mb-4">Select Entity Type</h5>
      <div class="field">
        <div class="control">
          <label class="radio">
            <input 
              type="radio" 
              bind:group={entityType} 
              value="asset"
              on:change={onEntityTypeChange}
            />
            Assets
          </label>
          <label class="radio ml-4">
            <input 
              type="radio" 
              bind:group={entityType} 
              value="person"
              on:change={onEntityTypeChange}
            />
            Persons
          </label>
        </div>
      </div>

      <hr>

      <FilterBuilder 
        {entityType}
        bind:filters
        {properties}
        {attributes}
      />

      <div class="buttons mt-4">
        <Button 
          color="primary" 
          loading={searching}
          disabled={filters.length === 0}
          on:click={runReport}
        >
          <span class="icon"><i class="fas fa-search"></i></span>
          <span>Run Report</span>
        </Button>
        
        <Button 
          color="light"
          disabled={filters.length === 0}
          on:click={clearFilters}
        >
          <span class="icon"><i class="fas fa-eraser"></i></span>
          <span>Clear</span>
        </Button>

        {#if hasSearched && results.length > 0}
          <Button 
            color="info"
            on:click={exportResults}
          >
            <span class="icon"><i class="fas fa-download"></i></span>
            <span>Export CSV</span>
          </Button>
        {/if}
      </div>
    </div>
  </Card>

  {#if searching}
    <Card>
      <Loading text="Searching..." />
    </Card>
  {:else if hasSearched}
    <Card>
      <div class="level mb-3">
        <div class="level-left">
          <div class="level-item">
            <h5 class="title is-5">Results ({results.length})</h5>
          </div>
        </div>
      </div>

      {#if results.length === 0}
        <div class="notification is-warning is-light">
          <p>No results found matching your criteria.</p>
        </div>
      {:else}
        <div class="table-container">
          <table class="table is-fullwidth is-striped is-hoverable">
            <thead>
              <tr>
                {#each columnHeaders as header}
                  <th>{header}</th>
                {/each}
              </tr>
            </thead>
            <tbody>
              {#each results as row}
                <tr>
                  {#each displayColumns as col}
                    <td>{formatValue(row[col])}</td>
                  {/each}
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </Card>
  {/if}
{/if}
