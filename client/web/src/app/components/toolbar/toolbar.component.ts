import {Component, Input} from '@angular/core';
import {MatIconModule} from "@angular/material/icon";
import {MatToolbarModule} from "@angular/material/toolbar";

@Component({
  selector: 'app-toolbar',
  templateUrl: './toolbar.component.html',
  styleUrl: './toolbar.component.scss',
  standalone: true,
  imports: [
    MatIconModule,
    MatToolbarModule
  ],
})
export class AppToolbarComponent {
  @Input() title: string = "";
  @Input() subtitle: string = "";
}
