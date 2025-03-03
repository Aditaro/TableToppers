import { Injectable } from '@angular/core';
import {HttpClient, HttpStatusCode} from '@angular/common/http';
import { Observable } from 'rxjs';

import { environment } from 'src/environments/environment';
import {NewTable, Table} from '../models/table.model';


@Injectable({
  providedIn: 'root'
})
export class TablesService {
  private apiUrl = environment.apiBaseUrl;

  // private mockTables: Table[] = [
  //   {
  //     id: '1',
  //     restaurantId: '1',
  //     name: 'Table 1',
  //     minCapacity: 2,
  //     maxCapacity: 4,
  //     status: 'available',
  //     x: 10,
  //     y: 10
  //   },
  //   {
  //     id: '2',
  //     restaurantId: '1',
  //     name: 'Table 2',
  //     minCapacity: 4,
  //     maxCapacity: 6,
  //     status: 'available',
  //     x: 100,
  //     y: 100
  //   },
  //   {
  //     id: '3',
  //     restaurantId: '1',
  //     name: 'Table 3',
  //     minCapacity: 6,
  //     maxCapacity: 8,
  //     status: 'available',
  //     x: 200,
  //     y: 200
  //   }
  // ];

  constructor(private http: HttpClient) {}

  getTables(restaurantId: string): Observable<Table[]> {
    // return this.mockTables;
    return this.http.get<Table[]>(`${this.apiUrl}/restaurants/${restaurantId}/tables`);
  }

  updateTable(restaurantId: string, tableId: string, payload: Partial<Table>): Observable<Table> {
    return this.http.put<Table>(`${this.apiUrl}/restaurants/${restaurantId}/tables/${tableId}`, payload);
  }

  addTable(restaurantId: string, payload: NewTable): Observable<Table> {
    return this.http.post<Table>(`${this.apiUrl}/restaurants/${restaurantId}/tables`, payload);
  }

  deleteTable(restaurantId: string, tableId: string): Observable<any> {
    return this.http.delete(`${this.apiUrl}/restaurants/${restaurantId}/tables/${tableId}`);
  }
}
