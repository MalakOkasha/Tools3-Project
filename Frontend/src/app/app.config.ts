
// import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
// import { provideRouter } from '@angular/router';
// import { routes } from './app.routes'; // This should now work correctly

// export const appConfig: ApplicationConfig = {
//   providers: [
//     provideZoneChangeDetection({ eventCoalescing: true }),
//     provideRouter(routes)
//   ]
// };



import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';
import { routes } from './app.routes';
import {authInterceptor} from './interceptors/auth.interceptor';
import {provideHttpClient, withInterceptors} from '@angular/common/http';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes),
    provideHttpClient(
      withInterceptors([authInterceptor])
    )
  ]
};
