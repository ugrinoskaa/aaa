import {Component} from '@angular/core';
import {FormBuilder, FormGroup, ReactiveFormsModule, Validators} from '@angular/forms';
import {MatStepperModule} from '@angular/material/stepper';
import {MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {CommonModule} from '@angular/common';
import {MatCardModule} from '@angular/material/card';
import {MatInputModule} from '@angular/material/input';
import {Router} from '@angular/router';
import {MatSnackBar, MatSnackBarModule} from '@angular/material/snack-bar';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {AppToolbarComponent} from '../../toolbar/toolbar.component';
import {APIService} from '../../../services/api.service';
import {Dataset, SourceType} from '../../../services/source.service';
import {MatGridListModule} from '@angular/material/grid-list';
import {MatListModule} from '@angular/material/list';

@Component({
  selector: 'app-source-add',
  templateUrl: './source-add.component.html',
  styleUrl: './source-add.component.scss',
  standalone: true,
  imports: [
    CommonModule,
    MatGridListModule,
    MatStepperModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatInputModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
    ReactiveFormsModule,
    MatListModule,
    AppToolbarComponent,
  ],
})
export class AppSourceAddComponent {
  public SourceType = SourceType;

  firstStep: FormGroup;
  secondStep: FormGroup;

  files: File[] = [];
  discovered: Dataset[] = [];
  isLoading: boolean = false;

  constructor(
    private fb: FormBuilder,
    private router: Router,
    private snack: MatSnackBar,
    private api: APIService
  ) {
    this.firstStep = this.fb.group({
      sourceType: ['', Validators.required]
    });

    this.secondStep = this.fb.group({
      name: ['', Validators.required],
      uri: [''],
      file: ['']
    });

    this.firstStep.get('sourceType')?.valueChanges.subscribe(value => {
      if (value === SourceType.POSTGRES) {
        this.secondStep.controls['uri'].setValidators([Validators.required]);
        this.secondStep.controls['file'].clearValidators();
      } else if (value === SourceType.CSV) {
        this.secondStep.controls['uri'].clearValidators();
        this.secondStep.controls['file'].setValidators([Validators.required]);
      }

      this.secondStep.updateValueAndValidity();
    });
  }

  onSourceSelect(source: SourceType) {
    this.firstStep.get('sourceType')?.setValue(source);
  }

  onConnect(stepper: any) {
    this.isLoading = true;
    const type = this.firstStep.get('sourceType')?.value
    const uri = this.secondStep.get('uri')?.value;

    if (type === SourceType.POSTGRES) {
      this.api.sources().connect(uri).subscribe({
        next: (source) => {
          this.snack.open('Connection Successful!', 'close', {duration: 3000});
          this.discovered = source.datasets!;
          this.isLoading = false;
          stepper.next();
        },
        error: (err) => {
          console.error(err);
          this.snack.open('Connection Failed!', 'close', {duration: 3000});
          this.isLoading = false;
        }
      })
    }

    if (type === SourceType.CSV) {
      this.api.sources().upload(this.files).subscribe({
        next: (source) => {
          this.snack.open('Upload Successful!', 'close', {duration: 3000});
          this.discovered = source.datasets!;
          this.isLoading = false;
          stepper.next();
        },
        error: (err) => {
          console.error(err);
          this.snack.open('Upload Failed!', 'close', {duration: 3000});
          this.isLoading = false;
        }
      })
    }
  }

  onSourceSave() {
    this.isLoading = true;
    const type = this.firstStep.get('sourceType')?.value
    const name = this.secondStep.get('name')?.value;
    const uri = this.secondStep.get('uri')?.value;
    const resource = uri || this.discovered[0].schema

    this.api.sources().create(name, type, resource).subscribe({
      next: (source) => {
        this.snack.open('New source created successfully', 'close', {duration: 3000});
        this.isLoading = false;
        this.router.navigate(['/datasets/new'], {
          queryParams: {
            source_id: source.id,
          }
        });
      },
      error: (err) => {
        this.snack.open('Source creation failed!', 'close', {duration: 3000});
        this.isLoading = false;
        console.error(err);
      }
    })
  }

  onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files) {
      this.files = Array.from(input.files);
    }
  }

  get discoveredLength() {
    if (this.firstStep.get('sourceType')?.value === SourceType.CSV) {
      return this.files.length;
    }
    return this.discovered.length;
  }

  get discoveredSources() {
    return this.discovered;
  }
}
