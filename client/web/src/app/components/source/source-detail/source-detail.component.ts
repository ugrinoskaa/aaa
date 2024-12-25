import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {CommonModule} from '@angular/common';
import {Source} from '../../../services/source.service';
import {APIService} from '../../../services/api.service';
import {MatIconModule} from '@angular/material/icon';
import {MatCardModule} from '@angular/material/card';

@Component({
  selector: 'app-source-detail',
  templateUrl: './source-detail.component.html',
  styleUrl: './source-detail.component.scss',
  imports: [CommonModule, MatIconModule, MatCardModule],
  standalone: true,
})
export class AppSourceDetailComponent implements OnInit {
  source: Source | undefined;

  constructor(
    private api: APIService,
    private route: ActivatedRoute,
  ) {
  }

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.api.sources().id(+id).subscribe({
        next: source => {
          this.source = source;
        },
        error: err => {
          console.error(err);
        }
      });
    }
  }
}
