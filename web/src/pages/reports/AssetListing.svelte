<script>
  import { onMount } from 'svelte';
  import { api, notifications } from '../../stores.js';
  import Card from '../../../../shared/components/Card.svelte';
  import Loading from '../../../../shared/components/Loading.svelte';

  let assets = [];
  let properties = [];
  let loading = true;
  let expandedAssetId = null;
  let assignmentHistory = {};
  let loadingHistory = {};

  // For edit modal integration
  let showEditModal = false;
  let editingAsset = null;

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    loading = true;
    try {
      const [assetsResult, propsResult] = await Promise.all([
        api.getAssetsWithAssignments(),
        api.getProperties()
      ]);
      const rawAssets = assetsResult || [];
      properties = propsResult || [];
      
      // Load properties for each asset
      assets = await Promise.all(rawAssets.map(async (asset) => {
        try {
          const assetProps = await api.getAssetProperties(asset.ID);
          return { ...asset, properties: assetProps || [] };
        } catch {
          return { ...asset, properties: [] };
        }
      }));
    } catch (err) {
      notifications.error('Failed to load data');
    } finally {
      loading = false;
    }
  }

  async function toggleExpand(assetId) {
    if (expandedAssetId === assetId) {
      expandedAssetId = null;
      return;
    }
    
    expandedAssetId = assetId;
    
    // Load assignment history if not already loaded
    if (!assignmentHistory[assetId]) {
      loadingHistory[assetId] = true;
      try {
        const history = await api.getAssetAssignments(assetId);
        assignmentHistory[assetId] = history || [];
      } catch {
        assignmentHistory[assetId] = [];
      } finally {
        loadingHistory[assetId] = false;
        loadingHistory = loadingHistory;
      }
    }
  }

  function getPropertyValue(asset, propId) {
    const prop = asset.properties?.find(p => p.PropertyID === propId);
    return prop?.Value || '-';
  }

  function formatDate(dateStr) {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleDateString();
  }

  function openEditModal(asset) {
    // Navigate to assets page with edit - using hash routing
    window.location.hash = `#/assets?edit=${asset.ID}`;
  }
</script>

<h1 class="title">Asset Listing Report</h1>

<Card>
  {#if loading}
    <Loading />
  {:else if assets.length === 0}
    <p class="has-text-grey">No assets found</p>
  {:else}
    <div class="table-container">
      <table class="table is-fullwidth is-hoverable">
        <thead>
          <tr>
            <th style="width: 30px;"></th>
            <th>Name</th>
            <th>Type</th>
            <th>Model</th>
            <th>Serial Number</th>
            <th>Assigned To</th>
            {#each properties as prop}
              <th>{prop.Name}</th>
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each assets as asset}
            <tr 
              class="is-clickable" 
              class:is-selected={expandedAssetId === asset.ID}
              on:click={() => toggleExpand(asset.ID)}
            >
              <td>
                <span class="icon is-small">
                  <i class="fas" class:fa-chevron-down={expandedAssetId === asset.ID} class:fa-chevron-right={expandedAssetId !== asset.ID}></i>
                </span>
              </td>
              <td>
                <a 
                  href="#/assets?edit={asset.ID}" 
                  on:click|stopPropagation
                  class="has-text-link"
                >
                  {asset.Name}
                </a>
              </td>
              <td>{asset.AssetTypeName || '-'}</td>
              <td>{asset.Model || '-'}</td>
              <td>{asset.SerialNumber || '-'}</td>
              <td>
                <span class:has-text-grey={!asset.CurrentAssignee || asset.CurrentAssignee === 'Unassigned'}>
                  {asset.CurrentAssignee || 'Unassigned'}
                </span>
              </td>
              {#each properties as prop}
                <td>{getPropertyValue(asset, prop.ID)}</td>
              {/each}
            </tr>
            {#if expandedAssetId === asset.ID}
              <tr>
                <td colspan={6 + properties.length} class="accordion-content">
                  <div class="box ml-5">
                    <h6 class="title is-6 mb-3">Assignment History</h6>
                    {#if loadingHistory[asset.ID]}
                      <p class="has-text-grey"><i class="fas fa-spinner fa-spin"></i> Loading...</p>
                    {:else if !assignmentHistory[asset.ID] || assignmentHistory[asset.ID].length === 0}
                      <p class="has-text-grey is-italic">No assignment history</p>
                    {:else}
                      <table class="table is-fullwidth is-narrow is-striped">
                        <thead>
                          <tr>
                            <th>Person</th>
                            <th>From</th>
                            <th>To</th>
                            <th>Notes</th>
                          </tr>
                        </thead>
                        <tbody>
                          {#each assignmentHistory[asset.ID] as assignment}
                            <tr>
                              <td>{assignment.PersonName}</td>
                              <td>{formatDate(assignment.EffectiveFrom)}</td>
                              <td>{assignment.EffectiveTo ? formatDate(assignment.EffectiveTo) : 'Current'}</td>
                              <td>{assignment.Notes || '-'}</td>
                            </tr>
                          {/each}
                        </tbody>
                      </table>
                    {/if}
                  </div>
                </td>
              </tr>
            {/if}
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Card>

<style>
  .is-clickable {
    cursor: pointer;
  }
  .accordion-content {
    background-color: #fafafa;
    padding: 0 !important;
  }
  .accordion-content .box {
    margin: 1rem;
    background-color: white;
  }
  tr.is-selected {
    background-color: #e8f4fc !important;
  }
  tr.is-selected td {
    border-color: #dbdbdb;
  }
</style>
