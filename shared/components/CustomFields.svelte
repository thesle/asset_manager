<script>
  import FormField from './FormField.svelte';
  
  // definitions: array of property/attribute definitions with ID, Name, DataType, EnumOptions
  // values: object mapping property/attribute ID to current value
  export let definitions = [];
  export let values = {};
  
  function handleChange(id, value) {
    values[id] = value;
    values = values; // Trigger reactivity
  }
  
  function getFieldType(dataType) {
    switch (dataType) {
      case 'int':
      case 'decimal':
        return 'number';
      case 'boolean':
        return 'checkbox';
      case 'date':
        return 'date';
      case 'datetime':
        return 'datetime-local';
      case 'enum':
        return 'select';
      default:
        return 'text';
    }
  }
  
  function getEnumOptions(enumOptionsJson) {
    if (!enumOptionsJson) return [];
    try {
      const options = JSON.parse(enumOptionsJson);
      return options.map(opt => ({ value: opt, label: opt }));
    } catch {
      return [];
    }
  }
  
  function getStep(dataType) {
    if (dataType === 'decimal') return '0.01';
    if (dataType === 'int') return '1';
    return undefined;
  }
</script>

{#if definitions.length > 0}
  <div class="custom-fields">
    {#each definitions as def}
      {#if def.DataType === 'boolean'}
        <div class="field">
          <label class="checkbox">
            <input 
              type="checkbox" 
              checked={values[def.ID] === 'true' || values[def.ID] === true}
              on:change={(e) => handleChange(def.ID, e.target.checked ? 'true' : 'false')}
            />
            {def.Name}
          </label>
        </div>
      {:else if def.DataType === 'enum'}
        <FormField
          label={def.Name}
          type="select"
          name={`custom_${def.ID}`}
          value={values[def.ID] || ''}
          options={getEnumOptions(def.EnumOptions)}
          on:input={(e) => handleChange(def.ID, e.target.value)}
        />
      {:else}
        <FormField
          label={def.Name}
          type={getFieldType(def.DataType)}
          name={`custom_${def.ID}`}
          value={values[def.ID] || ''}
          step={getStep(def.DataType)}
          on:input={(e) => handleChange(def.ID, e.target.value)}
        />
      {/if}
    {/each}
  </div>
{:else}
  <p class="has-text-grey-light is-italic">No custom fields defined</p>
{/if}

<style>
  .custom-fields {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
</style>
