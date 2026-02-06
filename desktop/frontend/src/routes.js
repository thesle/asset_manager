import Config from './pages/Config.svelte';
import Login from './pages/Login.svelte';

// Import pages from web (reuse them)
import Dashboard from '../../../web/src/pages/Dashboard.svelte';
import Assets from '../../../web/src/pages/Assets.svelte';
import AssetDetail from '../../../web/src/pages/AssetDetail.svelte';
import Persons from '../../../web/src/pages/Persons.svelte';
import PersonDetail from '../../../web/src/pages/PersonDetail.svelte';
import Assignments from '../../../web/src/pages/Assignments.svelte';
import AssetTypes from '../../../web/src/pages/config/AssetTypes.svelte';
import Properties from '../../../web/src/pages/config/Properties.svelte';
import Attributes from '../../../web/src/pages/config/Attributes.svelte';
import Users from '../../../web/src/pages/config/Users.svelte';
import Profile from '../../../web/src/pages/Profile.svelte';

export const routes = {
  '/': Dashboard,
  '/config': Config,
  '/login': Login,
  '/assets': Assets,
  '/assets/:id': AssetDetail,
  '/persons': Persons,
  '/persons/:id': PersonDetail,
  '/assignments': Assignments,
  '/config/asset-types': AssetTypes,
  '/config/properties': Properties,
  '/config/attributes': Attributes,
  '/config/users': Users,
  '/profile': Profile,
  '*': Dashboard,
};
