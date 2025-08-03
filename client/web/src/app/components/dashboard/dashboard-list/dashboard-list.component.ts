import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {CommonModule} from '@angular/common';
import {MatSnackBar} from '@angular/material/snack-bar';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {Dashboard} from '../../../services/dashboard.service';
import {APIService} from '../../../services/api.service';
import {AppListComponent, ListProps} from '../../common/list/list.component';

@Component({
  selector: 'app-dashboard-list',
  templateUrl: './dashboard-list.component.html',
  styleUrl: './dashboard-list.component.scss',
  imports: [CommonModule, MatProgressSpinnerModule, AppListComponent],
  standalone: true,
})
export class AppDashboardListComponent implements OnInit {
  isLoading: boolean = false;
  dashboards: Dashboard[] = [];
  properties: ListProps = {
    title: 'Dashboards',
    description: 'Manage your dashboards to keep your data stories connected and up to date.',
    imageUrl: 'assets/card_dashboard.svg',
    emptyTitle: 'No Dashboards Yet?',
    emptyDescription: 'Create your first dashboard to begin visualizing your data and unlock powerful insights at a glance.',
    emptyImageUrl: 'assets/empty_dashboard.svg',
    addNewLink: '/dashboards/new',
    addNewTitle: 'Add another dashboard to organize new perspectives and monitor your data more effectively.',
    addNewButton: 'DASHBOARD',
  }

  constructor(
    private api: APIService,
    private snack: MatSnackBar,
    private router: Router,
  ) {
  }

  ngOnInit() {
    this.isLoading = true;
    this.api.dashboards().all().subscribe({
      next: dashboards => {
        this.dashboards = dashboards;
        this.isLoading = false;
      },
      error: err => {
        console.error(err);
        this.snack.open('Unable to fetch Dashboards!', 'close', {duration: 3000});
        this.isLoading = false;
      }
    })
  }

  onDashboardDelete(id: number) {
    this.api.dashboards().delete(id).subscribe({
      next: _ => {
        this.dashboards = this.dashboards.filter(d => d.id !== id);
      },
      error: err => {
        this.snack.open('Unable to delete Dashboard!', 'close', {duration: 3000});
        console.error(err);
      }
    })
  }

  onDashboardClick(id: number) {
    this.router.navigate([`/dashboards/${id}`]).catch(console.error);
  }
}
