import { bootstrapApplication } from '@angular/platform-browser';
import {BrowserAnimationsModule, provideAnimations} from '@angular/platform-browser/animations';
import { appConfig } from './app/app.config';
import { provideHttpClient } from '@angular/common/http'; // Import provideHttpClient
import { AppComponent } from './app/app.component';
import {importProvidersFrom} from '@angular/core';
import {provideRouter} from '@angular/router';
import {routes} from './app/app.routes';

bootstrapApplication(AppComponent, {
  ...appConfig,
  providers: [
    ...(appConfig.providers || []),
    provideAnimations(),
    provideHttpClient(),
    provideRouter(routes),
    importProvidersFrom(BrowserAnimationsModule)
  ]
}).catch((err) => console.error(err));
