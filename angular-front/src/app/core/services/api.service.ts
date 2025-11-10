import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { environment } from '../../../environments/environment';
import {
  ChatResponse,
  ProductDetailsResponse,
  SessionResponse,
  PreferencesResponse,
  ActiveSessionResponse,
  SavedSearchBackend,
  UserPreferences,
} from '../../shared/types';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = environment.apiUrl;

  private getHeaders(accessToken?: string): HttpHeaders {
    let headers = new HttpHeaders({
      'Content-Type': 'application/json',
    });

    if (accessToken) {
      headers = headers.set('Authorization', `Bearer ${accessToken}`);
    }

    return headers;
  }

  private handleError(error: any): Observable<never> {
    console.error('API Error:', error);
    const message = error.error?.error || error.message || 'An error occurred';
    return throwError(() => new Error(message));
  }

  // Chat endpoints
  sendMessage(
    message: string,
    sessionId: string,
    accessToken?: string,
    country?: string,
    language?: string,
    currency?: string,
    newSearch?: boolean
  ): Observable<ChatResponse> {
    const body: any = {
      message,
      session_id: sessionId,
    };

    if (country) body.country = country;
    if (language) body.language = language;
    if (currency) body.currency = currency;
    if (newSearch !== undefined) body.new_search = newSearch;

    return this.http
      .post<ChatResponse>(`${this.baseUrl}/api/chat`, body, {
        headers: this.getHeaders(accessToken),
      })
      .pipe(catchError(this.handleError));
  }

  getSession(sessionId: string, accessToken?: string): Observable<SessionResponse> {
    return this.http
      .get<SessionResponse>(`${this.baseUrl}/api/session/${sessionId}`, {
        headers: this.getHeaders(accessToken),
      })
      .pipe(catchError(this.handleError));
  }

  // Product endpoints
  getProductDetails(
    pageToken: string,
    accessToken?: string
  ): Observable<ProductDetailsResponse> {
    return this.http
      .post<ProductDetailsResponse>(
        `${this.baseUrl}/api/product-details`,
        { page_token: pageToken },
        {
          headers: this.getHeaders(accessToken),
        }
      )
      .pipe(catchError(this.handleError));
  }

  // Auth endpoints
  checkAuth(): Observable<{ authenticated: boolean; user?: any }> {
    return this.http
      .get<{ authenticated: boolean; user?: any }>(`${this.baseUrl}/api/auth/check`, {
        withCredentials: true,
      })
      .pipe(catchError(this.handleError));
  }

  logout(): Observable<{ message: string }> {
    return this.http
      .post<{ message: string }>(
        `${this.baseUrl}/api/auth/logout`,
        {},
        {
          withCredentials: true,
        }
      )
      .pipe(catchError(this.handleError));
  }

  // User preferences endpoints
  getPreferences(accessToken: string): Observable<PreferencesResponse> {
    return this.http
      .get<PreferencesResponse>(`${this.baseUrl}/api/user/preferences`, {
        headers: this.getHeaders(accessToken),
      })
      .pipe(catchError(this.handleError));
  }

  updatePreferences(
    preferences: Partial<UserPreferences>,
    accessToken: string
  ): Observable<{ message: string }> {
    return this.http
      .put<{ message: string }>(
        `${this.baseUrl}/api/user/preferences`,
        preferences,
        {
          headers: this.getHeaders(accessToken),
        }
      )
      .pipe(catchError(this.handleError));
  }

  saveSearch(
    savedSearch: SavedSearchBackend,
    accessToken: string
  ): Observable<{ message: string }> {
    return this.http
      .post<{ message: string }>(
        `${this.baseUrl}/api/user/saved-search`,
        { saved_search: savedSearch },
        {
          headers: this.getHeaders(accessToken),
        }
      )
      .pipe(catchError(this.handleError));
  }

  clearSavedSearch(accessToken: string): Observable<{ message: string }> {
    return this.http
      .delete<{ message: string }>(`${this.baseUrl}/api/user/saved-search`, {
        headers: this.getHeaders(accessToken),
      })
      .pipe(catchError(this.handleError));
  }

  // Active session endpoints
  getActiveSession(accessToken: string): Observable<ActiveSessionResponse> {
    return this.http
      .get<ActiveSessionResponse>(`${this.baseUrl}/api/user/active-session`, {
        headers: this.getHeaders(accessToken),
      })
      .pipe(catchError(this.handleError));
  }

  setActiveSession(
    sessionId: string,
    accessToken: string
  ): Observable<{ message: string }> {
    return this.http
      .post<{ message: string }>(
        `${this.baseUrl}/api/user/active-session`,
        { session_id: sessionId },
        {
          headers: this.getHeaders(accessToken),
        }
      )
      .pipe(catchError(this.handleError));
  }

  clearActiveSession(accessToken: string): Observable<{ message: string }> {
    return this.http
      .delete<{ message: string }>(`${this.baseUrl}/api/user/active-session`, {
        headers: this.getHeaders(accessToken),
      })
      .pipe(catchError(this.handleError));
  }

  // Health check
  healthCheck(): Observable<any> {
    return this.http
      .get(`${this.baseUrl}/health`)
      .pipe(catchError(this.handleError));
  }

  // Statistics endpoints
  getStats(): Observable<any> {
    return this.http
      .get(`${this.baseUrl}/api/stats/all`)
      .pipe(catchError(this.handleError));
  }
}
