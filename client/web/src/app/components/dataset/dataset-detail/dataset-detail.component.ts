import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {CommonModule} from '@angular/common';
import {MatIconModule} from '@angular/material/icon';
import {MatCardModule} from '@angular/material/card';
import {MatTableModule} from '@angular/material/table';
import {APIService} from '../../../services/api.service';
import {Column, Dataset} from '../../../services/dataset.service';

@Component({
  selector: 'app-dataset-detail',
  templateUrl: './dataset-detail.component.html',
  styleUrl: './dataset-detail.component.scss',
  imports: [CommonModule, MatIconModule, MatCardModule, MatTableModule],
  standalone: true,
})
export class AppDatasetDetailComponent implements OnInit {
  dataset: Dataset | undefined;
  columns: (Column | undefined)[] = [];

  constructor(
    private api: APIService,
    private route: ActivatedRoute,
  ) {
  }

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.api.datasets().id(+id).subscribe({
        next: dataset => {
          this.dataset = dataset;
          this.columns = (dataset.columns || []).map(c => {
            const parts = c.split("::")
            return {name: parts[0], type: parts[1]};
          })
        },
        error: err => {
          console.error(err)
          this.columns = []
        }
      });
    }
  }
}
