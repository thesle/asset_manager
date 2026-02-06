/**
 * Creates an API client for the Asset Manager API
 * @param {string} baseUrl - The base URL of the API
 * @param {function} getToken - Function that returns the current auth token
 * @param {function} onUnauthorized - Callback when 401 is received
 */
export function createApiClient(baseUrl, getToken, onUnauthorized) {
  async function request(method, path, data = null) {
    const headers = {
      'Content-Type': 'application/json',
    };

    const token = getToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    const options = {
      method,
      headers,
    };

    if (data && (method === 'POST' || method === 'PUT' || method === 'PATCH')) {
      options.body = JSON.stringify(data);
    }

    const response = await fetch(`${baseUrl}${path}`, options);

    if (response.status === 401) {
      if (onUnauthorized) {
        onUnauthorized();
      }
      throw new Error('Unauthorized');
    }

    if (!response.ok) {
      const error = await response.json().catch(() => ({ Error: 'Request failed' }));
      throw new Error(error.Error || 'Request failed');
    }

    if (response.status === 204) {
      return null;
    }

    return response.json();
  }

  return {
    // Auth
    login: (username, password, remember) => 
      request('POST', '/api/auth/login', { Username: username, Password: password, Remember: remember }),
    me: () => request('GET', '/api/auth/me'),
    changePassword: (currentPassword, newPassword) => 
      request('POST', '/api/auth/change-password', { CurrentPassword: currentPassword, NewPassword: newPassword }),

    // Users
    getUsers: () => request('GET', '/api/users'),
    getUser: (id) => request('GET', `/api/users/${id}`),
    createUser: (data) => request('POST', '/api/users', data),
    updateUser: (id, data) => request('PUT', `/api/users/${id}`, data),
    resetUserPassword: (id, password) => request('POST', `/api/users/${id}/reset-password`, { Password: password }),
    deleteUser: (id) => request('DELETE', `/api/users/${id}`),

    // Asset Types
    getAssetTypes: () => request('GET', '/api/asset-types'),
    getAssetType: (id) => request('GET', `/api/asset-types/${id}`),
    createAssetType: (data) => request('POST', '/api/asset-types', data),
    updateAssetType: (id, data) => request('PUT', `/api/asset-types/${id}`, data),
    deleteAssetType: (id) => request('DELETE', `/api/asset-types/${id}`),

    // Assets
    getAssets: () => request('GET', '/api/assets'),
    getAssetsWithAssignments: () => request('GET', '/api/assets/with-assignments'),
    getAsset: (id) => request('GET', `/api/assets/${id}`),
    getAssetsByType: (typeId) => request('GET', `/api/assets/by-type/${typeId}`),
    searchAssets: (query) => request('GET', `/api/assets/search?q=${encodeURIComponent(query)}`),
    createAsset: (data) => request('POST', '/api/assets', data),
    updateAsset: (id, data) => request('PUT', `/api/assets/${id}`, data),
    deleteAsset: (id) => request('DELETE', `/api/assets/${id}`),
    getAssetProperties: (id) => request('GET', `/api/assets/${id}/properties`),
    setAssetProperty: (id, data) => request('POST', `/api/assets/${id}/properties`, data),
    deleteAssetProperty: (id, propId) => request('DELETE', `/api/assets/${id}/properties/${propId}`),

    // Properties
    getProperties: () => request('GET', '/api/properties'),
    getProperty: (id) => request('GET', `/api/properties/${id}`),
    createProperty: (data) => request('POST', '/api/properties', data),
    updateProperty: (id, data) => request('PUT', `/api/properties/${id}`, data),
    deleteProperty: (id) => request('DELETE', `/api/properties/${id}`),

    // Persons
    getPersons: () => request('GET', '/api/persons'),
    getPerson: (id) => request('GET', `/api/persons/${id}`),
    searchPersons: (query) => request('GET', `/api/persons/search?q=${encodeURIComponent(query)}`),
    createPerson: (data) => request('POST', '/api/persons', data),
    updatePerson: (id, data) => request('PUT', `/api/persons/${id}`, data),
    deletePerson: (id) => request('DELETE', `/api/persons/${id}`),
    getPersonAttributes: (id) => request('GET', `/api/persons/${id}/attributes`),
    setPersonAttribute: (id, data) => request('POST', `/api/persons/${id}/attributes`, data),
    deletePersonAttribute: (id, attrId) => request('DELETE', `/api/persons/${id}/attributes/${attrId}`),

    // Attributes
    getAttributes: () => request('GET', '/api/attributes'),
    getAttribute: (id) => request('GET', `/api/attributes/${id}`),
    createAttribute: (data) => request('POST', '/api/attributes', data),
    updateAttribute: (id, data) => request('PUT', `/api/attributes/${id}`, data),
    deleteAttribute: (id) => request('DELETE', `/api/attributes/${id}`),

    // Assignments
    getAssetAssignments: (assetId) => request('GET', `/api/assignments/asset/${assetId}`),
    getCurrentAssetAssignment: (assetId) => request('GET', `/api/assignments/asset/${assetId}/current`),
    getPersonAssignments: (personId) => request('GET', `/api/assignments/person/${personId}`),
    getCurrentPersonAssignments: (personId) => request('GET', `/api/assignments/person/${personId}/current`),
    createAssignment: (data) => request('POST', '/api/assignments', data),
    assignAsset: (assetId, personId, notes, effectiveDate) => 
      request('POST', '/api/assignments/assign', { AssetID: assetId, PersonID: personId, Notes: notes, EffectiveDate: effectiveDate }),
    unassignAsset: (assetId) => request('POST', `/api/assignments/unassign/${assetId}`),
    updateAssignment: (id, data) => request('PUT', `/api/assignments/${id}`, data),
    endAssignment: (id, endDate) => request('POST', `/api/assignments/${id}/end`, { EndDate: endDate }),
    deleteAssignment: (id) => request('DELETE', `/api/assignments/${id}`),
  };
}
