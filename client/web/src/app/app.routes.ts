import {Routes} from '@angular/router';
import {HomeComponent} from "../../../../../../Desktop/AAA/client/web/src/app/components/home/home.component";
import {AppSourceListComponent} from './components/source/source-list/source-list.component';
import {AppSourceAddComponent} from './components/source/source-add/source-add.component';
import {AppSourceDetailComponent} from './components/source/source-detail/source-detail.component';
import {AppDatasetDetailComponent} from './components/dataset/dataset-detail/dataset-detail.component';
import {AppDatasetAddComponent} from './components/dataset/dataset-add/dataset-add.component';
import {AppDatasetListComponent} from './components/dataset/dataset-list/dataset-list.component';
import {AppChartListComponent} from './components/chart/chart-list/chart-list.component';
import {AppChartAddComponent} from './components/chart/chart-add/chart-add.component';
import {AppChartDetailComponent} from './components/chart/chart-detail/chart-detail.component';
import {AppDashboardDetailComponent} from './components/dashboard/dashboard-detail/dashboard-detail.component';
import {AppDashboardAddComponent} from './components/dashboard/dashboard-add/dashboard-add.component';
import {AppDashboardListComponent} from './components/dashboard/dashboard-list/dashboard-list.component';

export const routes: Routes = [
  {
    path: 'dashboards',
    title: 'A&A | Dashboards',
    component: AppDashboardListComponent
  },
  {
    path: 'dashboards/new',
    title: 'A&A | Add Dashboard',
    component: AppDashboardAddComponent
  },
  {
    path: 'dashboards/:id',
    title: 'A&A | Dashboard',
    component: AppDashboardDetailComponent
  },
  {
    path: 'charts',
    title: 'A&A | Charts',
    component: AppChartListComponent
  },
  {
    path: 'charts/new',
    title: 'A&A | Add Chart',
    component: AppChartAddComponent
  },
  {
    path: 'charts/:id',
    title: 'A&A | Chart',
    component: AppChartDetailComponent
  },
  {
    path: 'datasets',
    title: 'A&A | Datasets',
    component: AppDatasetListComponent
  },
  {
    path: 'datasets/new',
    title: 'A&A | Add Dataset',
    component: AppDatasetAddComponent
  },
  {
    path: 'datasets/:id',
    title: 'A&A | Dataset',
    component: AppDatasetDetailComponent
  },
  {
    path: 'sources',
    title: 'A&A | Sources',
    component: AppSourceListComponent,
  },
  {
    path: 'sources/new',
    title: 'A&A | Add Source',
    component: AppSourceAddComponent,
  },
  {
    path: 'sources/:id',
    title: 'A&A | Source',
    component: AppSourceDetailComponent,
  },
  {
    path: '',
    title: 'A&A | Home',
    component: HomeComponent,
    pathMatch: 'full',
  },
  {
    path: '**',
    component: HomeComponent
  }
];
