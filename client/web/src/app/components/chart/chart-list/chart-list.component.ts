import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {CommonModule} from '@angular/common';
import {MatSnackBar} from '@angular/material/snack-bar';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {Chart} from '../../../services/chart.service';
import {APIService} from '../../../services/api.service';
import {AppListComponent, ListProps} from '../../common/list/list.component';

@Component({
  selector: 'app-chart-list',
  templateUrl: './chart-list.component.html',
  styleUrl: './chart-list.component.scss',
  imports: [CommonModule, MatProgressSpinnerModule, AppListComponent],
  standalone: true,
})
export class AppChartListComponent implements OnInit {
  isLoading: boolean = false;
  charts: Chart[] = [];
  properties: ListProps = {
    title: 'Charts',
    description: 'Manage your charts to fine-tune how your data is presented and interpreted.',
    imageUrl: 'assets/',
    emptyTitle: 'No Charts Yet?',
    emptyDescription: 'Create your first chart to visually explore your data and better understand patterns and trends.',
    emptyImageUrl: 'assets/empty_chart.svg',
    addNewLink: '/charts/new',
    addNewTitle: 'Add another chart to further visualize your data and uncover meaningful patterns.',
    addNewButton: 'CHART',
  }

  constructor(
    private api: APIService,
    private snack: MatSnackBar,
    private router: Router,
  ) {
  }

  ngOnInit() {
    this.isLoading = true;
    this.api.charts().all().subscribe({
      next: charts => {
        this.charts = charts;
        this.isLoading = false;
      },
      error: err => {
        console.error(err);
        this.snack.open('Unable to fetch Charts!', 'close', {duration: 3000});
        this.isLoading = false;
      }
    })
  }

  onChartDelete(id: number) {
    this.api.charts().delete(id).subscribe({
      next: _ => {
        this.charts = this.charts.filter(d => d.id !== id);
      },
      error: err => {
        this.snack.open('Unable to delete Chart!', 'close', {duration: 3000});
        console.error(err);
      }
    })
  }

  onChartClick(id: number) {
    this.router.navigate([`/charts/${id}`]).catch(console.error);
  }
}
