import {AfterViewInit, Component, ElementRef, EventEmitter, Input, OnDestroy, Output, ViewChild} from '@angular/core';
import {MatCardModule} from '@angular/material/card';
import {MatIconModule} from '@angular/material/icon';
import {MatButtonModule} from '@angular/material/button';
import {MatToolbarModule} from '@angular/material/toolbar';
import {NgxEchartsModule} from 'ngx-echarts';
import {EChartsOption} from 'echarts';
import {APIService} from '../../services/api.service';
import {CommonModule} from '@angular/common';

@Component({
  selector: 'app-echart-preview',
  templateUrl: './echart-preview.component.html',
  styleUrl: './echart-preview.component.scss',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatIconModule,
    MatButtonModule,
    MatToolbarModule,
    NgxEchartsModule,
  ],
})
export class AppEChartPreviewComponent implements AfterViewInit, OnDestroy {
  @Input() chartId: number | undefined;
  @Input() readonly: boolean = false;
  @Output() deleteChart: EventEmitter<number> = new EventEmitter<number>();

  @ViewChild('wrapper') wrapper!: ElementRef;
  observer!: ResizeObserver;
  width: string = "0px";
  height: string = "0px";
  label: string = '';

  options: EChartsOption = {};

  constructor(readonly api: APIService) {
  }

  ngAfterViewInit() {
    if (this.chartId) {
      this.api.charts().run(this.chartId).subscribe({
        next: chart => {
          this.options = chart.options;
        }
      });
    }

    this.observer = new ResizeObserver(() => {
      this.width = (this.wrapper.nativeElement.clientWidth) + "px";
      this.height = (this.wrapper.nativeElement.clientHeight - 8) + "px";
    });

    this.observer.observe(this.wrapper.nativeElement);
  }

  ngOnDestroy() {
    if (this.observer) {
      this.observer.disconnect();
    }
  }

  onDelete() {
    if (this.chartId) {
      this.deleteChart.emit(this.chartId);
    }
  }
}
