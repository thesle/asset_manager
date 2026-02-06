<script>
  import Router, { location } from 'svelte-spa-router';
  import { routes } from './routes.js';
  import { auth, notifications } from './stores.js';
  import Navbar from '../../shared/components/Navbar.svelte';
  import Sidebar from '../../shared/components/Sidebar.svelte';
  import Notification from '../../shared/components/Notification.svelte';

  $: isAuthenticated = $auth.isAuthenticated;
  $: user = $auth.user;
  $: currentPath = $location;

  function handleLogout() {
    auth.logout();
    window.location.hash = '#/login';
  }

  // Redirect to login if not authenticated
  $: if (!isAuthenticated && currentPath !== '/login') {
    window.location.hash = '#/login';
  }
</script>

<Notification notifications={$notifications} onRemove={(id) => notifications.remove(id)} />

{#if isAuthenticated}
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
