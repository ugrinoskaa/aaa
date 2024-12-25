import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {CommonModule} from '@angular/common';
import {MatSnackBar} from '@angular/material/snack-bar';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {Source} from '../../../services/source.service';
import {APIService} from '../../../services/api.service';
import {AppListComponent, ListProps} from '../../common/list/list.component';

@Component({
  selector: 'app-source-list',
  templateUrl: './source-list.component.html',
  styleUrl: './source-list.component.scss',
  imports: [CommonModule, MatProgressSpinnerModule, AppListComponent],
  standalone: true,
})
export class AppSourceListComponent implements OnInit {
  isLoading: boolean = false;
  sources: Source[] = [];
  properties: ListProps = {
    title: 'Sources',
    description: 'Manage your sources to keep data connections organized and up to date.',
    imageUrl: 'assets/',
    emptyTitle: 'No Sources Yet?',
    emptyDescription: 'Create your first source to connect external data and start building your datasets.',
    emptyImageUrl: 'assets/empty_source.svg',
    addNewLink: '/sources/new',
    addNewTitle: 'Add another source to continue integrating external data and expanding your data landscape.',
    addNewButton: 'SOURCE',
  }

  constructor(
    private api: APIService,
    private snack: MatSnackBar,
    private router: Router,
  ) {
  }

  ngOnInit() {
    this.isLoading = true;
    this.api.sources().all().subscribe({
      next: sources => {
        this.sources = sources;
        this.isLoading = false;
      },
      error: err => {
        console.error(err);
        this.snack.open('Unable to fetch Sources!', 'close', {duration: 3000});
        this.isLoading = false;
      }
    })
  }

  onSourceDelete(id: number) {
    this.api.sources().delete(id).subscribe({
      next: _ => {
        this.sources = this.sources.filter(d => d.id !== id);
      },
      error: err => {
        this.snack.open('Unable to delete Source!', 'close', {duration: 3000});
        console.error(err);
      }
    })
  }

  onSourceClick(id: number) {
    this.router.navigate([`/sources/${id}`]).catch(console.error)
  }
}
