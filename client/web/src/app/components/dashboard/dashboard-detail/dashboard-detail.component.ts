import {Component, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Dashboard} from '../../../services/dashboard.service';
import {MatSnackBar} from '@angular/material/snack-bar';
import {CommonModule} from '@angular/common';
import {MatGridListModule} from '@angular/material/grid-list';
import {GridsterComponent, GridsterConfig, GridsterModule} from 'angular-gridster2';
import {AppEChartPreviewComponent} from '../../echart/echart-preview.component';
import {APIService} from '../../../services/api.service';
import {AppToolbarComponent} from '../../toolbar/toolbar.component';

@Component({
  selector: 'app-dashboard-detail',
  standalone: true,
  imports: [
    CommonModule,
    GridsterModule,
    MatGridListModule,
    AppToolbarComponent,
    AppEChartPreviewComponent
  ],
  templateUrl: './dashboard-detail.component.html',
  styleUrl: './dashboard-detail.component.scss',
})
export class AppDashboardDetailComponent implements OnInit {
  @ViewChild(GridsterComponent) gridster?: GridsterComponent;

  dashboard: Dashboard | null = null;
  config: GridsterConfig = {
    gridType: 'verticalFixed',
    compactType: 'none',
  };

  constructor(
    private api: APIService,
    private route: ActivatedRoute,
    private snack: MatSnackBar
  ) {
  }

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.api.dashboards().id(+id).subscribe({
        next: (dashboard) => {
          this.dashboard = dashboard;
        },
        error: (err) => {
          this.snack.open('Failed to load dashboard', 'close', {duration: 3000});
          console.error(err);
        }
      });
    }
  }
}
