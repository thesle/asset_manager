<script>
  import { push } from 'svelte-spa-router';
  import { auth, config, notifications, getApi, saveConfig } from '../stores.js';
  import FormField from '../../../../shared/components/FormField.svelte';
  import Button from '../../../../shared/components/Button.svelte';

  let username = '';
  let password = '';
  let remember = true;
  let loading = false;
  let error = '';

  async function handleLogin() {
    if (!username || !password) {
      error = 'Please enter username and password';
      return;
    }

    loading = true;
    error = '';

    try {
      const api = getApi();
      const response = await api.login(username, password, remember);
      
      // Save token to config
      await saveConfig($config.apiUrl, response.Token);
      
      auth.login(response.Token, response.User);
      notifications.success('Login successful');
      push('/');
    } catch (err) {
      error = err.message || 'Login failed';
    } finally {
      loading = false;
    }
  }
</script>

<section class="hero is-primary is-fullheight">
  <div class="hero-body">
    <div class="container">
      <div class="columns is-centered">
        <div class="column is-4">
          <div class="box">
            <h1 class="title has-text-centered">Asset Manager</h1>
            <h2 class="subtitle has-text-centered has-text-grey">Sign in to continue</h2>

            {#if error}
              <div class="notification is-danger is-light">
                {error}
              </div>
            {/if}

            <form on:submit|preventDefault={handleLogin}>
              <FormField
                label="Username"
                name="username"
                bind:value={username}
                placeholder="Enter your username"
                required
              />

              <FormField
                label="Password"
                type="password"
                name="password"
                bind:value={password}
                placeholder="Enter your password"
                required
              />

              <FormField
                type="checkbox"
                name="remember"
                bind:value={remember}
                placeholder="Remember me"
              />

              <Button type="submit" color="primary" fullwidth {loading}>
                Sign In
              </Button>
            </form>

            <hr>
            <p class="has-text-centered">
              <a href="#/config" class="has-text-grey">Change server settings</a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</section>
