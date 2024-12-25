import {Routes} from '@angular/router';
import {HomeComponent} from "../../../../../../Desktop/AAA/client/web/src/app/components/home/home.component";
import {AppSourceListComponent} from './components/source/source-list/source-list.component';
import {AppSourceAddComponent} from './components/source/source-add/source-add.component';
import {AppSourceDetailComponent} from './components/source/source-detail/source-detail.component';

export const routes: Routes = [
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
