<script>
  export let value = '';
  export let placeholder = 'Search...';
  export let onSearch = () => {};
  export let debounceMs = 300;

  let timeout;

  function handleInput(event) {
    value = event.target.value;
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      onSearch(value);
    }, debounceMs);
  }

  function handleSubmit(event) {
    event.preventDefault();
    clearTimeout(timeout);
    onSearch(value);
  }

  function handleClear() {
    value = '';
    onSearch('');
  }
</script>

<form on:submit={handleSubmit}>
  <div class="field has-addons">
    <div class="control has-icons-left has-icons-right is-expanded">
      <input
        class="input"
        type="text"
        {placeholder}
        {value}
        on:input={handleInput}
      />
      <span class="icon is-small is-left">
        <i class="fas fa-search"></i>
      </span>
      {#if value}
        <span class="icon is-small is-right is-clickable" on:click={handleClear} on:keypress={handleClear} role="button" tabindex="0">
          <i class="fas fa-times"></i>
        </span>
      {/if}
    </div>
    <div class="control">
      <button class="button is-primary" type="submit">
        Search
      </button>
    </div>
  </div>
</form>

<style>
  .is-clickable {
    cursor: pointer;
    pointer-events: auto !important;
  }
</style>
