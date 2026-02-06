<script>
  export let columns = [];
  export let data = [];
  export let loading = false;
  export let emptyMessage = 'No data available';
  export let onRowClick = null;
  export let sortable = true;

  let sortColumn = null;
  let sortDirection = 'asc';

  function handleSort(column) {
    if (!sortable || !column.sortable) return;
    
    if (sortColumn === column.key) {
      sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
    } else {
      sortColumn = column.key;
      sortDirection = 'asc';
    }
  }

  $: sortedData = sortColumn
    ? [...data].sort((a, b) => {
        const aVal = a[sortColumn];
        const bVal = b[sortColumn];
        const modifier = sortDirection === 'asc' ? 1 : -1;
        
        if (aVal == null) return 1;
        if (bVal == null) return -1;
        if (typeof aVal === 'string') {
          return aVal.localeCompare(bVal) * modifier;
        }
        return (aVal - bVal) * modifier;
      })
    : data;
</script>

<div class="table-container">
  <table class="table is-fullwidth is-hoverable is-striped">
    <thead>
      <tr>
        {#each columns as column}
          <th 
            class:is-clickable={sortable && column.sortable !== false}
            on:click={() => handleSort(column)}
          >
            {column.label}
            {#if sortColumn === column.key}
              <span class="icon is-small">
                {#if sortDirection === 'asc'}
                  <i class="fas fa-sort-up"></i>
                {:else}
                  <i class="fas fa-sort-down"></i>
                {/if}
              </span>
            {/if}
          </th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#if loading}
        <tr>
          <td colspan={columns.length} class="has-text-centered">
            <span class="icon">
              <i class="fas fa-spinner fa-spin"></i>
            </span>
            Loading...
          </td>
        </tr>
      {:else if sortedData.length === 0}
        <tr>
          <td colspan={columns.length} class="has-text-centered has-text-grey">
            {emptyMessage}
          </td>
        </tr>
      {:else}
        {#each sortedData as row, index}
          <tr 
            class:is-clickable={onRowClick}
            on:click={() => onRowClick && onRowClick(row)}
          >
            {#each columns as column}
              <td>
                {#if column.render}
                  {@html column.render(row[column.key], row, index)}
                {:else if column.component}
                  <svelte:component this={column.component} value={row[column.key]} {row} />
                {:else}
                  {row[column.key] ?? ''}
                {/if}
              </td>
            {/each}
          </tr>
        {/each}
      {/if}
    </tbody>
  </table>
</div>

<style>
  .is-clickable {
    cursor: pointer;
  }
  .is-clickable:hover {
    background-color: #f5f5f5;
  }
  th.is-clickable:hover {
    background-color: #e8e8e8;
  }
</style>
