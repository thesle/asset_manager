<script>
  export let currentPage = 1;
  export let totalPages = 1;
  export let onPageChange = () => {};

  $: pages = (() => {
    const result = [];
    const delta = 2;
    const left = currentPage - delta;
    const right = currentPage + delta + 1;
    let l;

    for (let i = 1; i <= totalPages; i++) {
      if (i === 1 || i === totalPages || (i >= left && i < right)) {
        if (l && i - l !== 1) {
          result.push({ type: 'ellipsis' });
        }
        result.push({ type: 'page', value: i });
        l = i;
      }
    }
    return result;
  })();

  function goToPage(page) {
    if (page >= 1 && page <= totalPages && page !== currentPage) {
      onPageChange(page);
    }
  }
</script>

{#if totalPages > 1}
  <nav class="pagination is-centered" aria-label="pagination">
    <button
      class="pagination-previous"
      disabled={currentPage === 1}
      on:click={() => goToPage(currentPage - 1)}
    >
      Previous
    </button>
    <button
      class="pagination-next"
      disabled={currentPage === totalPages}
      on:click={() => goToPage(currentPage + 1)}
    >
      Next
    </button>
    <ul class="pagination-list">
      {#each pages as page}
        {#if page.type === 'ellipsis'}
          <li><span class="pagination-ellipsis">&hellip;</span></li>
        {:else}
          <li>
            <button
              class="pagination-link"
              class:is-current={page.value === currentPage}
              aria-label="Go to page {page.value}"
              on:click={() => goToPage(page.value)}
            >
              {page.value}
            </button>
          </li>
        {/if}
      {/each}
    </ul>
  </nav>
{/if}
