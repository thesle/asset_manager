<script>
  import { onMount } from "svelte";
  import { api, notifications } from "../../stores.js";
  import Card from "../../../../shared/components/Card.svelte";
  import Loading from "../../../../shared/components/Loading.svelte";

  let persons = [];
  let attributes = [];
  let properties = [];
  let loading = true;
  let expandedPersonId = null;
  let personAssets = {};
  let loadingAssets = {};
  let searchTerm = "";
  let showDeleted = false;

  $: filteredPersons = searchTerm.trim()
    ? persons.filter((person) => {
        const term = searchTerm.toLowerCase();
        return (
          (person.Name || "").toLowerCase().includes(term) ||
          (person.Email || "").toLowerCase().includes(term) ||
          (person.Phone || "").toLowerCase().includes(term) ||
          (person.attributes || []).some((a) => (a.Value || "").toLowerCase().includes(term))
        );
      })
    : persons;

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    loading = true;
    expandedPersonId = null;
    personAssets = {};
    try {
      const [personsResult, attrsResult, propsResult] = await Promise.all([
        api.getPersons(showDeleted),
        api.getAttributes(),
        api.getProperties(),
      ]);
      persons = (personsResult || []).filter((p) => p.Name !== "Unassigned");
      attributes = attrsResult || [];
      properties = propsResult || [];

      // Load attributes for each person
      for (const person of persons) {
        try {
          const personAttrs = await api.getPersonAttributes(person.ID);
          person.attributes = personAttrs || [];
        } catch {
          person.attributes = [];
        }
      }
    } catch (err) {
      notifications.error("Failed to load data");
    } finally {
      loading = false;
    }
  }

  async function toggleDeleted() {
    showDeleted = !showDeleted;
    searchTerm = "";
    await loadData();
  }

  async function toggleExpand(personId) {
    if (expandedPersonId === personId) {
      expandedPersonId = null;
      return;
    }

    expandedPersonId = personId;

    // Load current assets if not already loaded
    if (!personAssets[personId]) {
      loadingAssets[personId] = true;
      loadingAssets = loadingAssets;
      try {
        const assignments = await api.getCurrentPersonAssignments(personId);
        const assets = assignments || [];

        // Load properties for each assigned asset
        const assetsWithProps = await Promise.all(
          assets.map(async (assignment) => {
            try {
              const assetProps = await api.getAssetProperties(assignment.AssetID);
              return { ...assignment, properties: assetProps || [] };
            } catch {
              return { ...assignment, properties: [] };
            }
          }),
        );

        personAssets[personId] = assetsWithProps;
        personAssets = personAssets; // Trigger reactivity
      } catch {
        personAssets[personId] = [];
        personAssets = personAssets; // Trigger reactivity
      } finally {
        loadingAssets[personId] = false;
        loadingAssets = loadingAssets;
      }
    }
  }

  function getAttributeValue(person, attrId) {
    const attr = person.attributes?.find((a) => a.AttributeID === attrId);
    if (!attr) return "-";

    // Find the attribute definition to check type
    const attrDef = attributes.find((a) => a.ID === attrId);
    if (attrDef?.DataType === "boolean") {
      return attr.Value === "true" ? "Yes" : "No";
    }
    return attr.Value || "-";
  }

  function getPropertyValue(assignment, propId) {
    const prop = assignment.properties?.find((p) => p.PropertyID === propId);
    return prop?.Value || "-";
  }

  function formatDate(dateStr) {
    if (!dateStr) return "-";
    return new Date(dateStr).toLocaleDateString();
  }
</script>

<h1 class="title">Person Listing Report</h1>

<Card>
  <div class="report-controls mb-4">
    <div class="field is-grouped">
      <div class="control has-icons-left is-expanded">
        <input class="input" type="text" placeholder="Search across all fields..." bind:value={searchTerm} />
        <span class="icon is-left">
          <i class="fas fa-search"></i>
        </span>
      </div>
      <div class="control">
        <button class="button" class:is-danger={showDeleted} class:is-outlined={!showDeleted} on:click={toggleDeleted}>
          <span class="icon is-small">
            <i class="fas" class:fa-trash={!showDeleted} class:fa-undo={showDeleted}></i>
          </span>
          <span>{showDeleted ? "Show Active" : "Show Deleted"}</span>
        </button>
      </div>
    </div>
    {#if showDeleted}
      <p class="help is-danger">Showing deleted records</p>
    {/if}
  </div>

  {#if loading}
    <Loading />
  {:else if persons.length === 0}
    <p class="has-text-grey">{showDeleted ? "No deleted persons found" : "No persons found"}</p>
  {:else if filteredPersons.length === 0}
    <p class="has-text-grey">No matching persons found</p>
  {:else}
    <div class="accordion-list">
      {#each filteredPersons as person}
        <div class="accordion-item" class:is-expanded={expandedPersonId === person.ID}>
          <div
            class="accordion-header is-clickable"
            on:click={() => toggleExpand(person.ID)}
            on:keydown={(e) => e.key === "Enter" && toggleExpand(person.ID)}
            role="button"
            tabindex="0"
          >
            <div class="accordion-icon">
              <span class="icon">
                <i
                  class="fas"
                  class:fa-chevron-down={expandedPersonId === person.ID}
                  class:fa-chevron-right={expandedPersonId !== person.ID}
                ></i>
              </span>
            </div>
            <div class="accordion-title">
              <a
                href="#/persons?edit={person.ID}"
                on:click|stopPropagation
                class="has-text-link has-text-weight-semibold"
              >
                {person.Name}
              </a>
            </div>
            <div class="accordion-details">
              <span class="tag is-light mr-2">
                <span class="icon is-small"><i class="fas fa-envelope"></i></span>
                <span>{person.Email || "No email"}</span>
              </span>
              <span class="tag is-light mr-2">
                <span class="icon is-small"><i class="fas fa-phone"></i></span>
                <span>{person.Phone || "No phone"}</span>
              </span>
              {#each attributes as attr}
                <span class="tag is-info is-light mr-1">
                  {attr.Name}: {getAttributeValue(person, attr.ID)}
                </span>
              {/each}
            </div>
          </div>

          {#if expandedPersonId === person.ID}
            <div class="accordion-content">
              <div class="box">
                <h6 class="title is-6 mb-3">Assigned Assets</h6>
                {#if loadingAssets[person.ID]}
                  <p class="has-text-grey"><i class="fas fa-spinner fa-spin"></i> Loading...</p>
                {:else if !personAssets[person.ID] || personAssets[person.ID].length === 0}
                  <p class="has-text-grey is-italic">No assets currently assigned</p>
                {:else}
                  <table class="table is-fullwidth is-striped">
                    <thead>
                      <tr>
                        <th>Asset</th>
                        <th>Type</th>
                        <th>Model</th>
                        <th>Serial Number</th>
                        <th>Assigned Since</th>
                        {#each properties as prop}
                          <th>{prop.Name}</th>
                        {/each}
                      </tr>
                    </thead>
                    <tbody>
                      {#each personAssets[person.ID] as assignment}
                        <tr>
                          <td>
                            <a href="#/assets?edit={assignment.AssetID}" class="has-text-link">
                              {assignment.AssetName}
                            </a>
                          </td>
                          <td>{assignment.AssetTypeName || "-"}</td>
                          <td>{assignment.AssetModel || "-"}</td>
                          <td>{assignment.AssetSerialNumber || "-"}</td>
                          <td>{formatDate(assignment.EffectiveFrom)}</td>
                          {#each properties as prop}
                            <td>{getPropertyValue(assignment, prop.ID)}</td>
                          {/each}
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                {/if}
              </div>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</Card>

<style>
  .report-controls {
    padding-bottom: 0.5rem;
    border-bottom: 1px solid #f0f0f0;
  }

  .accordion-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .accordion-item {
    border: 1px solid #dbdbdb;
    border-radius: 4px;
    overflow: hidden;
  }

  .accordion-item.is-expanded {
    border-color: #3273dc;
  }

  .accordion-header {
    display: flex;
    align-items: center;
    padding: 0.75rem 1rem;
    background-color: #f5f5f5;
    gap: 0.75rem;
  }

  .accordion-header:hover {
    background-color: #eeeeee;
  }

  .is-clickable {
    cursor: pointer;
  }

  .accordion-icon {
    flex-shrink: 0;
  }

  .accordion-title {
    flex-shrink: 0;
    min-width: 200px;
  }

  .accordion-details {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
    align-items: center;
  }

  .accordion-content {
    padding: 0;
    background-color: #fafafa;
  }

  .accordion-content .box {
    margin: 1rem;
    background-color: white;
  }
</style>
