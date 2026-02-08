<script>
  export let label = '';
  export let type = 'text';
  export let value = '';
  export let placeholder = '';
  export let required = false;
  export let disabled = false;
  export let error = '';
  export let help = '';
  export let name = '';
  export let options = []; // For select type
  export let rows = 3; // For textarea
</script>

<div class="field">
  {#if label}
    <label class="label" for={name}>
      {label}
      {#if required}<span class="has-text-danger">*</span>{/if}
    </label>
  {/if}
  
  <div class="control" class:has-icons-left={$$slots.iconLeft} class:has-icons-right={$$slots.iconRight || error}>
    {#if type === 'textarea'}
      <textarea
        class="textarea"
        class:is-danger={error}
        {name}
        id={name}
        {placeholder}
        {required}
        {disabled}
        {rows}
        bind:value
        on:input
        on:change
        on:blur
      ></textarea>
    {:else if type === 'select'}
      <div class="select is-fullwidth" class:is-danger={error}>
        <select
          {name}
          id={name}
          {required}
          {disabled}
          bind:value
          on:change
        >
          <option value="">{placeholder || 'Select...'}</option>
          {#each options as option}
            <option value={option.value}>{option.label}</option>
          {/each}
        </select>
      </div>
    {:else if type === 'checkbox'}
      <label class="checkbox">
        <input
          type="checkbox"
          {name}
          id={name}
          {disabled}
          bind:checked={value}
          on:change
        />
        {placeholder}
      </label>
    {:else if type === 'password'}
      <input
        class="input"
        class:is-danger={error}
        type="password"
        {name}
        id={name}
        {placeholder}
        {required}
        {disabled}
        bind:value
        on:input
        on:change
        on:blur
      />
    {:else if type === 'email'}
      <input
        class="input"
        class:is-danger={error}
        type="email"
        {name}
        id={name}
        {placeholder}
        {required}
        {disabled}
        bind:value
        on:input
        on:change
        on:blur
      />
    {:else if type === 'number'}
      <input
        class="input"
        class:is-danger={error}
        type="number"
        {name}
        id={name}
        {placeholder}
        {required}
        {disabled}
        bind:value
        on:input
        on:change
        on:blur
      />
    {:else if type === 'date'}
      <input
        class="input"
        class:is-danger={error}
        type="date"
        {name}
        id={name}
        {required}
        {disabled}
        bind:value
        on:input
        on:change
        on:blur
      />
    {:else}
      <input
        class="input"
        class:is-danger={error}
        type="text"
        {name}
        id={name}
        {placeholder}
        {required}
        {disabled}
        bind:value
        on:input
        on:change
        on:blur
      />
    {/if}
    
    {#if $$slots.iconLeft}
      <span class="icon is-small is-left">
        <slot name="iconLeft"></slot>
      </span>
    {/if}
    
    {#if error}
      <span class="icon is-small is-right has-text-danger">
        <i class="fas fa-exclamation-triangle"></i>
      </span>
    {:else if $$slots.iconRight}
      <span class="icon is-small is-right">
        <slot name="iconRight"></slot>
      </span>
    {/if}
  </div>
  
  {#if error}
    <p class="help is-danger">{error}</p>
  {:else if help}
    <p class="help">{help}</p>
  {/if}
</div>
