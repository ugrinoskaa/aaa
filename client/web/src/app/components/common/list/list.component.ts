import {Component, EventEmitter, Input, Output} from '@angular/core';
import {CommonModule} from '@angular/common';
import {MatIconModule} from '@angular/material/icon';
import {MatButtonModule} from '@angular/material/button';
import {MatCardModule} from '@angular/material/card';
import {RouterLink} from '@angular/router';
import {AppCardComponent} from '../card/card.component';

export interface ListData {
  id?: number;
  name?: string;
  type?: string;
}

export interface ListProps {
  title: string;
  description: string;
  imageUrl: string;
  emptyTitle: string;
  emptyDescription: string;
  emptyImageUrl: string;
  addNewLink: string;
  addNewTitle: string;
  addNewButton: string;
}

@Component({
  selector: 'app-list',
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
  standalone: true,
  imports: [
    RouterLink,
    CommonModule,
    MatCardModule,
    MatIconModule,
    MatButtonModule,
    AppCardComponent,
  ],
})
export class AppListComponent {
  @Input() props!: ListProps;
  @Input() data!: ListData[];

  @Output() onClick = new EventEmitter<number>();
  @Output() onDelete = new EventEmitter<number>();

  onCardClick(id: number) {
    this.onClick.emit(id);
  }

  onCardDelete(id: number) {
    this.onDelete.emit(id);
  }
}
