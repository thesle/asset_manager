<script>
  export let title = 'Asset Manager';
  export let user = null;
  export let onLogout = () => {};

  let menuActive = false;
</script>

<nav class="navbar is-primary" role="navigation" aria-label="main navigation">
  <div class="navbar-brand">
    <a class="navbar-item" href="/">
      <strong>{title}</strong>
    </a>

    <button
      class="navbar-burger"
      class:is-active={menuActive}
      aria-label="menu"
      aria-expanded={menuActive}
      on:click={() => menuActive = !menuActive}
    >
      <span aria-hidden="true"></span>
      <span aria-hidden="true"></span>
      <span aria-hidden="true"></span>
    </button>
  </div>

  <div class="navbar-menu" class:is-active={menuActive}>
    <div class="navbar-start">
      <slot name="nav-start"></slot>
    </div>

    <div class="navbar-end">
      <slot name="nav-end"></slot>
      
      {#if user}
        <div class="navbar-item has-dropdown is-hoverable">
          <a class="navbar-link">
            <span class="icon">
              <i class="fas fa-user"></i>
            </span>
            <span>{user.Username}</span>
          </a>
          <div class="navbar-dropdown is-right">
            <a class="navbar-item" href="/profile">
              <span class="icon"><i class="fas fa-user-cog"></i></span>
              <span>Profile</span>
            </a>
            <hr class="navbar-divider">
            <a class="navbar-item" href="#logout" on:click|preventDefault={onLogout}>
              <span class="icon"><i class="fas fa-sign-out-alt"></i></span>
              <span>Logout</span>
            </a>
          </div>
        </div>
      {/if}
    </div>
  </div>
</nav>
