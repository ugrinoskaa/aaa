import {Component} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {AppNavbarComponent} from './components/navbar/navbar.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    AppNavbarComponent
  ],
  template: `<app-navbar></app-navbar>`,
})
export class AppComponent {
}
