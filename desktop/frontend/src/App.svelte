<script>
  import { onMount } from 'svelte';
  import Router, { location, push } from 'svelte-spa-router';
  import { routes } from './routes.js';
  import { auth, config, notifications, initConfig } from './stores.js';
  import Navbar from '../../../shared/components/Navbar.svelte';
  import Sidebar from '../../../shared/components/Sidebar.svelte';
  import Notification from '../../../shared/components/Notification.svelte';
  import Loading from '../../../shared/components/Loading.svelte';

  let loading = true;

  $: isAuthenticated = $auth.isAuthenticated;
  $: isConfigured = $config.configured;
  $: user = $auth.user;
  $: currentPath = $location;

  onMount(async () => {
    await initConfig();
    loading = false;
    
    // Redirect to config if not configured
    if (!$config.configured) {
      push('/config');
    }
  });

  function handleLogout() {
    auth.logout();
    push('/login');
  }

  // Redirect logic
  $: if (!loading) {
    if (!isConfigured && currentPath !== '/config') {
      push('/config');
    } else if (isConfigured && !isAuthenticated && currentPath !== '/login' && currentPath !== '/config') {
      push('/login');
    }
  }
</script>

<Notification notifications={$notifications} onRemove={(id) => notifications.remove(id)} />

{#if loading}
  <Loading fullPage text="Loading..." />
{:else if !isConfigured}
  <Router {routes} />
{:else if isAuthenticated}
  <Navbar title="Asset Manager" {user} onLogout={handleLogout} />
  <div class="main-content">
    <Sidebar {currentPath} />
    <main class="content-area">
      <Router {routes} />
    </main>
  </div>
{:else}
  <Router {routes} />
{/if}
