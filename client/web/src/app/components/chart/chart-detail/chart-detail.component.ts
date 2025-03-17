import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {CommonModule} from '@angular/common';
import {MatIconModule} from '@angular/material/icon';
import {NgxEchartsModule} from 'ngx-echarts';
import {APIService} from '../../../services/api.service';
import {Chart} from '../../../services/chart.service';
import {EChartsOption} from 'echarts';

@Component({
  selector: 'app-chart-detail',
  templateUrl: './chart-detail.component.html',
  styleUrl: './chart-detail.component.scss',
  imports: [CommonModule, MatIconModule, NgxEchartsModule],
  standalone: true,
})
export class AppChartDetailComponent implements OnInit {
  chart: Chart | undefined;
  options: EChartsOption = {};

  constructor(
    private api: APIService,
    private route: ActivatedRoute,
  ) {
  }

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.api.charts().id(+id).subscribe({
        next: chart => {
          this.chart = chart;
          this.api.charts().run(chart.id!).subscribe({
            next: rsp => {
              this.options = rsp.options;
            }
          });
        },
        error: err => {
          console.error(err)
        }
      });
    }
  }

  getFilter(filter: string, part: number) {
    return filter.split("/")[part]
  }
}
