import {Component, OnInit} from '@angular/core';
import {FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators} from '@angular/forms';
import {MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {MatInputModule} from '@angular/material/input';
import {MatSelectModule} from '@angular/material/select';
import {CommonModule} from '@angular/common';
import {MatCardModule} from "@angular/material/card";
import {MatTableModule} from "@angular/material/table";
import {Source} from '../../../services/source.service';
import {APIService} from '../../../services/api.service';
import {Column} from '../../../services/dataset.service';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {AppToolbarComponent} from '../../toolbar/toolbar.component';
import {MatSnackBar, MatSnackBarModule} from '@angular/material/snack-bar';
import {ActivatedRoute, Router} from '@angular/router';
import {MatGridListModule} from '@angular/material/grid-list';

@Component({
  selector: 'app-dataset-add',
  templateUrl: './dataset-add.component.html',
  styleUrl: './dataset-add.component.scss',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatButtonModule,
    MatIconModule,
    MatInputModule,
    MatSelectModule,
    ReactiveFormsModule,
    MatCardModule,
    MatTableModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
    AppToolbarComponent,
    MatGridListModule,
  ],
})
export class AppDatasetAddComponent implements OnInit {
  form: FormGroup;
  sources: Source[] = [];
  selected?: Source;

  schemas: (string | undefined)[] = [];
  tables: (string | undefined)[] = [];
  columns: (Column | undefined)[] = [{name: 'username', type: 'string'}, {name: 'created_at', type: 'datetime'}];

  isLoading = false;

  constructor(
    readonly api: APIService,
    private router: Router,
    private route: ActivatedRoute,
    private snack: MatSnackBar,
    private fb: FormBuilder,
  ) {
    this.form = this.fb.group({
      name: ['', Validators.required],
      source: ['', Validators.required],
      schema: ['', Validators.required],
      table: ['', Validators.required],
    });

    this.form.get('source')?.valueChanges.subscribe(id => {
      if (id) {
        this.loadSource(id)
      } else {
        this.schemas = [];
        this.selected = undefined;
      }
    });

    this.form.get('schema')?.valueChanges.subscribe(schema => {
      if (this.selected && schema) {
        this.tables = this.selected.datasets?.filter(d => d.schema === schema).map(d => d.table) || [];
      } else {
        this.tables = [];
      }
    });

    this.form.get('table')?.valueChanges.subscribe(table => {
      if (this.selected && table) {
        this.columns = this.selected.datasets
          ?.filter(d => d.schema === this.form.get('schema')?.value && d.table === table)
          .flatMap(d => d.columns)
          .map(column => {
            const parts = column!.split("::")
            return ({name: parts[0], type: parts[1]})
          }) || [];
      } else {
        this.columns = [];
      }
    });
  }

  ngOnInit(): void {
    this.api.sources().all().subscribe({
      next: sources => {
        this.sources = sources;
        this.route.queryParams.subscribe(params => {
          this.form.patchValue({source: params['source_id']});
        });
      },
      error: err => {
        console.error("failed to fetch sources", err)
      }
    })
  }

  loadSource(id: number): void {
    this.isLoading = true;
    this.api.sources().id(id).subscribe({
      next: source => {
        this.selected = source;
        this.schemas = Array.from(new Set(source.datasets?.map(d => d.schema) || []));
        this.isLoading = false;
      },
      error: err => {
        console.error('failed to get source by id', err);
        this.isLoading = false;
      }
    });
  }

  onSubmit(): void {
    if (this.form.valid) {
      const values = this.form.value
      this.api.datasets().create({
        name: values.name,
        sourceId: values.source,
        sourceSchema: values.schema,
        sourceTable: values.table,
      }).subscribe({
        next: dataset => {
          this.snack.open('Dataset successfully created!', 'close', {duration: 3000});
          this.router.navigate(['/charts/new'], {
            queryParams: {
              dataset_id: dataset.id,
            }
          });
        },
        error: err => {
          console.error("Unable to save the Dataset!", err)
        }
      })
    }
  }
}
