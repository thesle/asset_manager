<script>
  export let title = '';
  export let subtitle = '';
</script>

<div class="card">
  {#if title || $$slots.header}
    <header class="card-header">
      {#if $$slots.header}
        <slot name="header"></slot>
      {:else}
        <p class="card-header-title">
          {title}
          {#if subtitle}
            <span class="has-text-grey-light ml-2 has-text-weight-normal is-size-7">{subtitle}</span>
          {/if}
        </p>
      {/if}
      {#if $$slots.headerIcon}
        <span class="card-header-icon">
          <slot name="headerIcon"></slot>
        </span>
      {/if}
    </header>
  {/if}
  
  <div class="card-content">
    <slot></slot>
  </div>
  
  {#if $$slots.footer}
    <footer class="card-footer">
      <slot name="footer"></slot>
    </footer>
  {/if}
</div>
