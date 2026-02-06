<script>
  import FormField from './FormField.svelte';

  export let dataType = 'string';
  export let value = '';
  export let label = '';
  export let name = '';
  export let required = false;
  export let disabled = false;
  export let enumOptions = '[]'; // JSON string of options

  $: parsedEnumOptions = (() => {
    try {
      const opts = JSON.parse(enumOptions || '[]');
      return opts.map(o => ({ value: o, label: o }));
    } catch {
      return [];
    }
  })();

  $: fieldType = (() => {
    switch (dataType) {
      case 'boolean': return 'checkbox';
      case 'int': return 'number';
      case 'decimal': return 'number';
      case 'date': return 'date';
      case 'datetime': return 'datetime-local';
      case 'enum': return 'select';
      default: return 'text';
    }
  })();
</script>

{#if dataType === 'boolean'}
  <FormField
    type="checkbox"
    {label}
    {name}
    {disabled}
    bind:value
    placeholder={label}
  />
{:else if dataType === 'enum'}
  <FormField
    type="select"
    {label}
    {name}
    {required}
    {disabled}
    bind:value
    options={parsedEnumOptions}
  />
{:else if dataType === 'int'}
  <FormField
    type="number"
    {label}
    {name}
    {required}
    {disabled}
    bind:value
  />
{:else if dataType === 'decimal'}
  <div class="field">
    {#if label}
      <label class="label" for={name}>
        {label}
        {#if required}<span class="has-text-danger">*</span>{/if}
      </label>
    {/if}
    <div class="control">
      <input
        class="input"
        type="number"
        step="0.01"
        {name}
        id={name}
        {required}
        {disabled}
        bind:value
      />
    </div>
  </div>
{:else if dataType === 'date'}
  <div class="field">
    {#if label}
      <label class="label" for={name}>
        {label}
        {#if required}<span class="has-text-danger">*</span>{/if}
      </label>
    {/if}
    <div class="control">
      <input
        class="input"
        type="date"
        {name}
        id={name}
        {required}
        {disabled}
        bind:value
      />
    </div>
  </div>
{:else if dataType === 'datetime'}
  <div class="field">
    {#if label}
      <label class="label" for={name}>
        {label}
        {#if required}<span class="has-text-danger">*</span>{/if}
      </label>
    {/if}
    <div class="control">
      <input
        class="input"
        type="datetime-local"
        {name}
        id={name}
        {required}
        {disabled}
        bind:value
      />
    </div>
  </div>
{:else}
  <FormField
    type="text"
    {label}
    {name}
    {required}
    {disabled}
    bind:value
  />
{/if}
