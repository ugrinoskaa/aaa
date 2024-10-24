import {Routes} from '@angular/router';
import {HomeComponent} from "../../../../../../Desktop/AAA/client/web/src/app/components/home/home.component";

export const routes: Routes = [
  {
    path: '**',
    component: HomeComponent
  }
];
