import {Component, OnInit} from '@angular/core';
import {EChartsOption} from 'echarts';
import {
  AbstractControl,
  FormArray,
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  ValidationErrors,
  ValidatorFn,
  Validators
} from '@angular/forms';
import {Router} from '@angular/router';
import {MatSnackBar, MatSnackBarModule} from '@angular/material/snack-bar';
import {MatDividerModule} from '@angular/material/divider';
import {MatCardModule} from '@angular/material/card';
import {MatButtonModule} from '@angular/material/button';
import {MatGridListModule} from '@angular/material/grid-list';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatDatepickerModule} from '@angular/material/datepicker';
import {MatInputModule} from '@angular/material/input';
import {MatSelectModule} from '@angular/material/select';
import {NgxEchartsModule} from 'ngx-echarts';
import {MatToolbarModule} from '@angular/material/toolbar';
import {CommonModule} from '@angular/common';
import {APIService} from '../../../services/api.service';
import {Dataset} from '../../../services/dataset.service';
import {ChartSchemaRules, Filter, ValidateChartReq} from '../../../services/chart.service';
import {AppToolbarComponent} from '../../toolbar/toolbar.component';
import {MatNativeDateModule} from '@angular/material/core';
import {MatIconModule} from '@angular/material/icon';

@Component({
  selector: 'app-chart-add',
  templateUrl: './chart-add.component.html',
  styleUrl: './chart-add.component.scss',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatIconModule,
    MatCardModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatToolbarModule,
    MatDividerModule,
    MatSnackBarModule,
    MatGridListModule,
    MatFormFieldModule,
    MatDatepickerModule,
    MatNativeDateModule,
    ReactiveFormsModule,
    NgxEchartsModule,
    AppToolbarComponent,
  ],
})
export class AppChartAddComponent implements OnInit {
  datasets: Dataset[] = [];
  chartTypes: string[] = [];
  chartType?: string;
  chartTypeSchema?: ChartSchemaRules;

  dimensions: string[] = [];
  metrics: string[] = [];

  data: { example: boolean, options: EChartsOption } = {example: true, options: {}}
  preview = this.generatePreviewChart()

  chartForm = new FormGroup({
    name: new FormControl('', [Validators.required]),
    dataset: new FormControl('', [Validators.required]),
    dimensions: new FormControl<string[]>([]),
    metrics: new FormControl<string[]>([]),
    filters: new FormArray<FormControl<Filter>>([])
  });

  get chartName() {
    return this.chartForm.get('name');
  }

  get chartDataset() {
    return this.chartForm.get('dataset');
  }

  get chartDimensions() {
    return this.chartForm.get('dimensions');
  }

  get chartMetrics() {
    return this.chartForm.get('metrics');
  }

  get chartFilters() {
    return this.chartForm.get('filters') as FormArray
  }

  constructor(
    private router: Router,
    readonly api: APIService,
    private snack: MatSnackBar,
  ) {
    this.update('disable')
    this.subscribe()

    this.chartForm.controls.dataset.valueChanges.subscribe(datasetId => {
      this.update('reset')

      if (!!datasetId) {
        this.api.datasets().id(parseInt(datasetId, 10)).subscribe({
          next: dataset => {
            this.dimensions = dataset.dimensions || [];
            this.metrics = dataset.metrics || [];
            this.update('enable')
          }
        })
      } else {
        this.update('disable')
      }
    });
  }

  ngOnInit() {
    this.api.charts().types().subscribe({
      next: ctypes => {
        this.chartTypes = ctypes;
        this.onSelectChartType(ctypes[0]);
      },
      error: err => {
        console.error(err)
      }
    });

    this.api.datasets().all().subscribe({
      next: datasets => {
        this.datasets = datasets;
      },
      error: err => {
        console.error(err);
      }
    });
  }

  onSelectChartType(chartType: string) {
    this.api.charts().type(chartType).subscribe({
      next: chart => {
        this.chartType = chartType;
        this.chartTypeSchema = chart.schema
        this.updateControls()

        if (this.isPreviewReady()) {
          this.preview();
        } else {
          this.preview(chart.example);
        }
      },
      error: err => {
        console.error(err);
      }
    });
  }

  generatePreviewChart() {
    let cached: ValidateChartReq | undefined

    return (example?: EChartsOption) => {
      if (example) {
        this.data = {
          example: true,
          options: example,
        }

        return;
      }

      const form = this.chartForm.value
      const req: ValidateChartReq = {
        name: form.name || '',
        type: this.chartType || '',
        datasetId: parseInt(form.dataset || '', 10),
        dimensions: form.dimensions || [],
        metrics: form.metrics || [],
        filters: (form.filters || []).map(f => `${f.dimension}/${f.operator}/${f.value}`),
      }

      if (req === cached) {
        return
      }

      this.api.charts().validate(req).subscribe({
        next: rsp => {
          if (rsp.valid) {
            this.data = {example: false, options: rsp.options}
          } else {
            console.error(rsp);
          }

          cached = req
        },
        error: err => {
          console.error(err);
        },

      })
    }
  }

  isPreviewReady() {
    return (
      this.chartDataset?.valid &&
      this.chartDimensions?.valid &&
      this.chartMetrics?.valid &&
      this.chartFilters?.valid
    );
  }

  onSubmit() {
    if (!this.chartForm.valid) {
      return
    }

    this.api.charts().create({
      type: this.chartType!,
      name: this.chartForm.value.name!,
      datasetId: parseInt(this.chartForm.value.dataset!, 10),
      dimensions: this.chartForm.value.dimensions!,
      metrics: this.chartForm.value.metrics!,
      filters: (this.chartForm.value.filters || []).map(f => `${f.dimension}/${f.operator}/${f.value}`)
    }).subscribe({
      next: chart => {
        this.snack.open('Chart successfully saved!', 'close', {duration: 3000});
        this.router.navigate(['/dashboards/new'], {
          queryParams: {
            chart_id: chart.id,
          }
        });
      },
      error: err => {
        this.snack.open('Unable to save Chart!', 'close', {duration: 3000});
        console.error(err);
      }
    });
  }

  private maxSelections(max: number): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const value = control.value as string[] | null;
      if (value && value.length > max) {
        return {maxSelections: {max, actual: value.length}};
      }
      return null;
    };
  }

  private updateControls() {
    const dimensions = this.controlDimensions()
    const metrics = this.controlMetrics()

    this.chartForm.controls.dimensions.setValidators(dimensions.validators)
    this.chartForm.controls.metrics.setValidators(metrics.validators)

    this.chartForm.patchValue({dimensions: dimensions.value, metrics: metrics.value})

    for (let message of [dimensions.snack, metrics.snack]) {
      if (!!message) {
        this.snack.open(message, 'close', {duration: 3000});
      }
    }
  }

  private controlDimensions() {
    const constrains = this.chartTypeSchema?.dimensions;
    const required = constrains?.min != 0
    const maximum = this.chartTypeSchema?.dimensions?.max || 1;
    const validators = required ? [Validators.required, this.maxSelections(maximum)] : [this.maxSelections(maximum)];
    let value = this.chartForm?.controls?.dimensions?.value || [];
    let snack: string | undefined = undefined

    if (value.length > maximum) {
      value = value.slice(0, maximum);
      snack = `Dimensions reduced to ${maximum} for this chart`
    }

    return {validators, value, snack}
  }

  private controlMetrics() {
    const constrains = this.chartTypeSchema?.metrics;
    const required = constrains?.min != 0
    const maximum = this.chartTypeSchema?.metrics?.max || 1;
    const validators = required ? [Validators.required, this.maxSelections(maximum)] : [this.maxSelections(maximum)];
    let value = this.chartForm?.controls?.metrics?.value || [];
    let snack: string | undefined = undefined

    if (value.length > maximum) {
      value = value.slice(0, maximum);
      snack = `Metrics reduced to ${maximum} for this chart`
    }

    return {validators, value, snack}
  }

  private update(state: 'reset' | 'disable' | 'enable') {
    const controls = [
      this.chartForm.controls.dimensions,
      this.chartForm.controls.metrics,
      this.chartForm.controls.filters,
    ]

    switch (state) {
      case "reset":
        this.dimensions = [];
        this.metrics = [];
        controls.forEach((control: AbstractControl) => {
          control.reset()
        })
        break;
      case "disable":
        controls.forEach((control: AbstractControl) => {
          control.disable()
        })
        break;
      case "enable":
        controls.forEach((control: AbstractControl) => {
          control.enable()
        })
        break;
    }
  }

  private subscribe() {
    this.chartForm.valueChanges.subscribe(_ => {
      if (this.isPreviewReady()) {
        this.preview();
      }
    })
  }

  noCompareFunction(): number {
    return 0;
  }

  addFilter() {
    this.chartFilters.push(
      new FormGroup({
        dimension: new FormControl('', [Validators.required]),
        operator: new FormControl('', [Validators.required]),
        value: new FormControl('', [Validators.required])
      })
    );
  }

  removeFilter(index: number) {
    this.chartFilters.removeAt(index);
  }
}
