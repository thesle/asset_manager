<script>
  import { auth, api, notifications } from '../stores.js';
  import Card from '../../../shared/components/Card.svelte';
  import FormField from '../../../shared/components/FormField.svelte';
  import Button from '../../../shared/components/Button.svelte';

  $: user = $auth.user;

  let passwordForm = { currentPassword: '', newPassword: '', confirmPassword: '' };
  let saving = false;
  let error = '';

  async function handleChangePassword() {
    error = '';
    
    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      error = 'New passwords do not match';
      return;
    }

    if (passwordForm.newPassword.length < 6) {
      error = 'Password must be at least 6 characters';
      return;
    }

    saving = true;
    try {
      await api.changePassword(passwordForm.currentPassword, passwordForm.newPassword);
      notifications.success('Password changed successfully');
      passwordForm = { currentPassword: '', newPassword: '', confirmPassword: '' };
    } catch (err) {
      error = err.message;
    } finally {
      saving = false;
    }
  }
</script>

<h1 class="title">Profile</h1>

<div class="columns">
  <div class="column is-6">
    <Card title="Account Information">
      <table class="table is-fullwidth">
        <tbody>
          <tr><th>Username</th><td>{user?.Username}</td></tr>
          <tr><th>Email</th><td>{user?.Email}</td></tr>
        </tbody>
      </table>
    </Card>
  </div>

  <div class="column is-6">
    <Card title="Change Password">
      {#if error}
        <div class="notification is-danger is-light">{error}</div>
      {/if}
      
      <form on:submit|preventDefault={handleChangePassword}>
        <FormField
          label="Current Password"
          type="password"
          name="currentPassword"
          bind:value={passwordForm.currentPassword}
          required
        />
        <FormField
          label="New Password"
          type="password"
          name="newPassword"
          bind:value={passwordForm.newPassword}
          required
        />
        <FormField
          label="Confirm New Password"
          type="password"
          name="confirmPassword"
          bind:value={passwordForm.confirmPassword}
          required
        />
        <Button type="submit" color="primary" loading={saving}>
          Change Password
        </Button>
      </form>
    </Card>
  </div>
</div>
