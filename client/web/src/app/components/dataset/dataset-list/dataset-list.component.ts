import {Component, OnInit} from '@angular/core';
import {Dataset} from '../../../services/dataset.service';
import {APIService} from '../../../services/api.service';
import {MatSnackBar} from '@angular/material/snack-bar';
import {Router} from '@angular/router';
import {CommonModule} from '@angular/common';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {AppListComponent, ListProps} from '../../common/list/list.component';

@Component({
  selector: 'app-dataset-list',
  templateUrl: './dataset-list.component.html',
  styleUrl: './dataset-list.component.scss',
  imports: [CommonModule, MatProgressSpinnerModule, AppListComponent],
  standalone: true,
})
export class AppDatasetListComponent implements OnInit {
  isLoading: boolean = false;
  datasets: Dataset[] = [];
  properties: ListProps = {
    title: 'Datasets',
    description: 'Manage your datasets to stay organized and ready for data analysis and exploration.',
    imageUrl: 'assets/card_dataset.svg',
    emptyTitle: 'No Datasets Yet?',
    emptyDescription: 'Create your first dataset to start exploring your data and uncover valuable insights.',
    emptyImageUrl: 'assets/empty_dataset.svg',
    addNewLink: '/datasets/new',
    addNewTitle: 'Add another dataset to continue expanding your data exploration and analysis.',
    addNewButton: 'DATASET',
  }

  constructor(
    private api: APIService,
    private snack: MatSnackBar,
    private router: Router,
  ) {
  }

  ngOnInit() {
    this.isLoading = true
    this.api.datasets().all().subscribe({
      next: datasets => {
        this.datasets = datasets;
        this.isLoading = false
      },
      error: err => {
        console.error(err);
        this.snack.open('Unable to fetch Datasets!', 'close', {duration: 3000});
        this.isLoading = false
      }
    })
  }

  onDatasetDelete(id: number) {
    this.api.datasets().delete(id).subscribe({
      next: _ => {
        this.datasets = this.datasets.filter(d => d.id !== id);
      },
      error: err => {
        this.snack.open('Unable to delete Dataset!', 'close', {duration: 3000});
        console.error(err);
      }
    })
  }

  onDatasetClick(id: number) {
    this.router.navigate([`/datasets/${id}`]).catch(console.error);
  }
}
