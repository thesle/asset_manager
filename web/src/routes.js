import Dashboard from './pages/Dashboard.svelte';
import Login from './pages/Login.svelte';
import Assets from './pages/Assets.svelte';
import AssetDetail from './pages/AssetDetail.svelte';
import Persons from './pages/Persons.svelte';
import PersonDetail from './pages/PersonDetail.svelte';
import Assignments from './pages/Assignments.svelte';
import AssetTypes from './pages/config/AssetTypes.svelte';
import Properties from './pages/config/Properties.svelte';
import Attributes from './pages/config/Attributes.svelte';
import Users from './pages/config/Users.svelte';
import Profile from './pages/Profile.svelte';
import AssetListing from './pages/reports/AssetListing.svelte';
import PersonListing from './pages/reports/PersonListing.svelte';
import CustomReport from './pages/reports/CustomReport.svelte';
import MultipleAssets from './pages/reports/MultipleAssets.svelte';

export const routes = {
  '/': Dashboard,
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
  '/reports/assets': AssetListing,
  '/reports/persons': PersonListing,
  '/reports/custom': CustomReport,
  '/reports/multiple-assets': MultipleAssets,
  '*': Dashboard,
};
