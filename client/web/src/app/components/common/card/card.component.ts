import {Component, EventEmitter, Input, Output} from '@angular/core';
import {MatCardModule} from '@angular/material/card';
import {MatIconModule} from '@angular/material/icon';
import {MatButtonModule} from '@angular/material/button';

@Component({
  selector: 'app-card',
  standalone: true,
  imports: [
    MatCardModule,
    MatIconModule,
    MatButtonModule,
  ],
  templateUrl: './card.component.html',
  styleUrl: './card.component.scss',
})
export class AppCardComponent {
  @Input() id?: number;
  @Input() name?: string;
  @Input() imageUrl!: string;
  @Output() onCardDelete = new EventEmitter<number>();
  @Output() onCardClick = new EventEmitter<number>();

  isFlipped = false;

  flipCard() {
    this.isFlipped = !this.isFlipped;
  }

  onDelete() {
    this.onCardDelete.emit(this.id);
    this.isFlipped = false;
  }

  onClick() {
    if (!this.isFlipped) {
      this.onCardClick.emit(this.id);
    }
  }
}
