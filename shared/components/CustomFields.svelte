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
    {#each definitions as def (def.ID)}
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
        <div class="field">
          <label class="label" for={`custom_enum_${def.ID}`}>{def.Name}</label>
          <div class="control">
            <div class="select is-fullwidth">
              <select
                id={`custom_enum_${def.ID}`}
                value={values[def.ID] || ''}
                on:change={(e) => handleChange(def.ID, e.target.value)}
              >
                <option value="">Select...</option>
                {#each getEnumOptions(def.EnumOptions) as opt}
                  <option value={opt.value}>{opt.label}</option>
                {/each}
              </select>
            </div>
          </div>
        </div>
      {:else}
        <FormField
          label={def.Name}
          type={getFieldType(def.DataType)}
          name={`custom_${def.ID}`}
          value={values[def.ID] || ''}
          step={getStep(def.DataType)}
          on:input={(e) => handleChange(def.ID, e.target.value)}
          on:change={(e) => handleChange(def.ID, e.target.value)}
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
