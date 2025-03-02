import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterModule } from '@angular/router';
import { TranslocoRootModule } from './transloco-root.module';
import { MatToolbarModule } from '@angular/material/toolbar';

@Component({
    selector: 'app-root',
    imports: [RouterOutlet, RouterLink, TranslocoRootModule, MatToolbarModule, RouterModule],
    templateUrl: './app.component.html',
    styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'Table Topper!';
}
