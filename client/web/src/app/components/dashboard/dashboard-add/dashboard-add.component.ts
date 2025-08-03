import {Component, OnInit} from '@angular/core';
import {MatCardModule} from '@angular/material/card';
import {MatButtonModule} from '@angular/material/button';
import {MatGridListModule} from '@angular/material/grid-list';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatInputModule} from '@angular/material/input';
import {AppToolbarComponent} from '../../toolbar/toolbar.component';
import {FormsModule} from '@angular/forms';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatSnackBar, MatSnackBarModule} from '@angular/material/snack-bar';
import {GridsterConfig, GridsterModule} from 'angular-gridster2';
import {CdkDragDrop, DragDropModule} from '@angular/cdk/drag-drop';
import {AppEChartPreviewComponent} from '../../echart/echart-preview.component';
import {Chart} from '../../../services/chart.service';
import {APIService} from '../../../services/api.service';
import {DashboardGrid, EChart} from '../../../services/dashboard.service';
import {CommonModule} from '@angular/common';
import {MatIconModule} from '@angular/material/icon';
import {Router} from '@angular/router';

@Component({
  selector: 'app-dashboard-add',
  templateUrl: './dashboard-add.component.html',
  styleUrl: './dashboard-add.component.scss',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatInputModule,
    MatButtonModule,
    MatToolbarModule,
    MatGridListModule,
    MatFormFieldModule,
    MatSnackBarModule,
    DragDropModule,
    GridsterModule,
    MatIconModule,
    AppToolbarComponent,
    AppEChartPreviewComponent,
  ],
})
export class AppDashboardAddComponent implements OnInit {
  dashboard: EChart[] = [];
  charts: Chart[] = [];
  name: string = '';

  config: GridsterConfig = {
    gridType: 'scrollVertical',
    compactType: 'none',
    defaultItemCols: 2,
    defaultItemRows: 2,
    fixedRowHeight: 100,
    minCols: 6,
    minRows: 6,
    maxRows: 100,
    scrollSensitivity: 10,
    scrollSpeed: 20,
    resizable: {enabled: true},
    draggable: {enabled: true},
  };

  constructor(
    private api: APIService,
    private snack: MatSnackBar,
    private router: Router,
  ) {
  }

  ngOnInit() {
    this.api.charts().all().subscribe({
      next: charts => {
        this.charts = charts;
      }
    });
  }

  drop(event: CdkDragDrop<Chart[]>) {
    const chart = this.charts[event.previousIndex];
    this.dashboard.push({cols: 2, rows: 2, x: 0, y: 0, label: chart.name, chartId: chart.id!});
  }

  onChartDelete(item: EChart) {
    this.dashboard = this.dashboard.filter(i => i !== item);
  }

  onDashboardSave() {
    if (this.dashboard.length && !!this.name) {
      const name = this.name;
      const grid = this.dashboard.map((d): DashboardGrid => {
        return {x: d.x, y: d.y, rows: d.rows, cols: d.cols, chartId: d.chartId}
      })

      this.api.dashboards().create({name, grid}).subscribe({
        next: dashboard => {
          this.snack.open('Dashboard successfully saved!', 'close', {duration: 3000});
          this.router.navigate([`/dashboards/${dashboard.id}`]);
        },
        error: err => {
          this.snack.open('Unable to save Dashboard!', 'close', {duration: 3000});
          console.error(err);
        }
      })
    }
  }
}
