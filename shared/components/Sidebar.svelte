<script>
  export let items = [];
  export let currentPath = '/';

  // Default menu items if none provided
  const defaultItems = [
    { path: '/', label: 'Dashboard', icon: 'fas fa-tachometer-alt' },
    { path: '/assets', label: 'Assets', icon: 'fas fa-boxes' },
    { path: '/persons', label: 'Persons', icon: 'fas fa-users' },
    { path: '/assignments', label: 'Assignments', icon: 'fas fa-exchange-alt' },
    { 
      label: 'Reports', 
      icon: 'fas fa-chart-bar',
      children: [
        { path: '/reports/assets', label: 'Asset Listing', icon: 'fas fa-boxes' },
        { path: '/reports/persons', label: 'Person Listing', icon: 'fas fa-users' },
      ]
    },
    { 
      label: 'Configuration', 
      icon: 'fas fa-cog',
      children: [
        { path: '/config/asset-types', label: 'Asset Types', icon: 'fas fa-tags' },
        { path: '/config/properties', label: 'Properties', icon: 'fas fa-list' },
        { path: '/config/attributes', label: 'Attributes', icon: 'fas fa-sliders-h' },
        { path: '/config/users', label: 'Users', icon: 'fas fa-user-cog' },
      ]
    },
  ];

  $: menuItems = items.length > 0 ? items : defaultItems;
</script>

<aside class="menu sidebar">
  <ul class="menu-list">
    {#each menuItems as item}
      {#if item.children}
        <li>
          <p class="menu-label">
            {#if item.icon}<span class="icon"><i class={item.icon}></i></span>{/if}
            {item.label}
          </p>
          <ul>
            {#each item.children as child}
              <li>
                <a href="#/{child.path.replace(/^\//, '')}" class:is-active={currentPath === child.path}>
                  {#if child.icon}<span class="icon"><i class={child.icon}></i></span>{/if}
                  {child.label}
                </a>
              </li>
            {/each}
          </ul>
        </li>
      {:else}
        <li>
          <a href="#/{item.path.replace(/^\//, '')}" class:is-active={currentPath === item.path}>
            {#if item.icon}<span class="icon"><i class={item.icon}></i></span>{/if}
            {item.label}
          </a>
        </li>
      {/if}
    {/each}
  </ul>
</aside>

<style>
  .sidebar {
    padding: 1rem;
    background-color: #f5f5f5;
    min-height: calc(100vh - 52px);
  }
  .menu-list a {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .menu-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
</style>
