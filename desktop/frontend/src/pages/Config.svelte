<script>
  import { push } from 'svelte-spa-router';
  import { config, notifications, saveConfig, getApi, auth } from '../stores.js';
  import FormField from '../../../../shared/components/FormField.svelte';
  import Button from '../../../../shared/components/Button.svelte';
  import Card from '../../../../shared/components/Card.svelte';

  let apiUrl = $config.apiUrl || 'http://localhost:8084';
  let username = '';
  let password = '';
  let remember = true;
  let loading = false;
  let error = '';
  let step = $config.configured ? 'login' : 'config';

  async function handleTestConnection() {
    loading = true;
    error = '';

    try {
      // Try to reach the API
      const response = await fetch(`${apiUrl}/api/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ Username: 'test', Password: 'test', Remember: false })
      });

      // Even a 401 means the API is reachable
      if (response.status === 401 || response.ok) {
        notifications.success('API connection successful');
        step = 'login';
      } else {
        error = 'Could not connect to API';
      }
    } catch (err) {
      error = 'Could not connect to API. Please check the URL.';
    } finally {
      loading = false;
    }
  }

  async function handleLogin() {
    if (!username || !password) {
      error = 'Please enter username and password';
      return;
    }

    loading = true;
    error = '';

    try {
      const response = await fetch(`${apiUrl}/api/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ Username: username, Password: password, Remember: remember })
      });

      if (!response.ok) {
        const data = await response.json();
        error = data.Error || 'Login failed';
        return;
      }

      const data = await response.json();
      
      // Save config with token
      await saveConfig(apiUrl, data.Token);
      
      // Update auth store
      auth.login(data.Token, data.User);
      
      notifications.success('Connected successfully');
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
        <div class="column is-5">
          <Card title="Asset Manager - Setup">
            {#if step === 'config'}
              <p class="mb-4">Configure the API server connection</p>
              
              {#if error}
                <div class="notification is-danger is-light">{error}</div>
              {/if}

              <FormField
                label="API Server URL"
                name="apiUrl"
                bind:value={apiUrl}
                placeholder="http://localhost:8084"
                help="The URL of the Asset Manager API server"
                required
              />

              <Button color="primary" fullwidth {loading} on:click={handleTestConnection}>
                Test Connection
              </Button>
            {:else}
              <p class="mb-4">Sign in to <strong>{apiUrl}</strong></p>
              
              {#if error}
                <div class="notification is-danger is-light">{error}</div>
              {/if}

              <form on:submit|preventDefault={handleLogin}>
                <FormField
                  label="Username"
                  name="username"
                  bind:value={username}
                  required
                />
                <FormField
                  label="Password"
                  type="password"
                  name="password"
                  bind:value={password}
                  required
                />
                <FormField
                  type="checkbox"
                  name="remember"
                  bind:value={remember}
                  placeholder="Remember me"
                />

                <div class="buttons">
                  <Button type="submit" color="primary" {loading}>
                    Sign In
                  </Button>
                  <Button on:click={() => step = 'config'}>
                    Change Server
                  </Button>
                </div>
              </form>
            {/if}
          </Card>
        </div>
      </div>
    </div>
  </div>
</section>
