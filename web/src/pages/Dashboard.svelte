<script>
  import { onMount } from 'svelte';
  import { api, notifications } from '../stores.js';
  import Card from '../../../shared/components/Card.svelte';
  import Loading from '../../../shared/components/Loading.svelte';

  let loading = true;
  let stats = {
    totalAssets: 0,
    totalPersons: 0,
    assetTypes: 0,
    assignedAssets: 0,
    recentAssignments: []
  };

  onMount(async () => {
    try {
      const [assets, persons, assetTypes] = await Promise.all([
        api.getAssetsWithAssignments(),
        api.getPersons(),
        api.getAssetTypes()
      ]);

      // Handle null/undefined responses
      const assetList = assets || [];
      const personList = persons || [];
      const typeList = assetTypes || [];

      stats = {
        totalAssets: assetList.length,
        totalPersons: personList.length,
        assetTypes: typeList.length,
        assignedAssets: assetList.filter(a => a.CurrentAssignee && a.CurrentAssignee !== '' && a.CurrentAssignee !== 'Unassigned').length
      };
    } catch (err) {
      notifications.error('Failed to load dashboard data');
    } finally {
      loading = false;
    }
  });
</script>

<h1 class="title">Dashboard</h1>

{#if loading}
  <Loading />
{:else}
  <div class="columns">
    <div class="column is-3">
      <Card title="Total Assets">
        <div class="has-text-centered">
          <p class="title is-1 has-text-primary">{stats.totalAssets}</p>
          <p class="subtitle is-6">Assets tracked</p>
        </div>
      </Card>
    </div>
    <div class="column is-3">
      <Card title="Assigned">
        <div class="has-text-centered">
          <p class="title is-1 has-text-success">{stats.assignedAssets}</p>
          <p class="subtitle is-6">Currently assigned</p>
        </div>
      </Card>
    </div>
    <div class="column is-3">
      <Card title="Persons">
        <div class="has-text-centered">
          <p class="title is-1 has-text-info">{stats.totalPersons}</p>
          <p class="subtitle is-6">People registered</p>
        </div>
      </Card>
    </div>
    <div class="column is-3">
      <Card title="Asset Types">
        <div class="has-text-centered">
          <p class="title is-1 has-text-warning">{stats.assetTypes}</p>
          <p class="subtitle is-6">Categories</p>
        </div>
      </Card>
    </div>
  </div>

  <div class="columns">
    <div class="column">
      <Card title="Quick Actions">
        <div class="buttons">
          <a href="#/assets" class="button is-primary">
            <span class="icon"><i class="fas fa-boxes"></i></span>
            <span>View Assets</span>
          </a>
          <a href="#/persons" class="button is-info">
            <span class="icon"><i class="fas fa-users"></i></span>
            <span>View Persons</span>
          </a>
          <a href="#/assignments" class="button is-success">
            <span class="icon"><i class="fas fa-exchange-alt"></i></span>
            <span>Manage Assignments</span>
          </a>
        </div>
      </Card>
    </div>
  </div>
{/if}
