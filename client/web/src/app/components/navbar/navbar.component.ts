import {Component} from '@angular/core';
import {RouterLink, RouterLinkActive, RouterModule} from '@angular/router';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatButtonModule} from "@angular/material/button";
import {MatIconModule} from '@angular/material/icon';
import {MatSidenavModule} from '@angular/material/sidenav';
import {MatListModule} from '@angular/material/list';
import {CommonModule} from '@angular/common';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.scss',
  standalone: true,
  imports: [
    CommonModule,
    MatToolbarModule,
    MatSidenavModule,
    MatButtonModule,
    MatIconModule,
    MatListModule,
    RouterModule,
    RouterLink,
    RouterLinkActive,
  ],
})
export class AppNavbarComponent {
  isSidebarOpen = true;

  toggleSidebar(): void {
    this.isSidebarOpen = !this.isSidebarOpen;
    setTimeout(() => {
      window.dispatchEvent(new Event('resize'));
    }, 200);
  }
}
