import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';
import {HttpClient} from '@angular/common/http';
import {GridsterItem} from 'angular-gridster2';

export interface EChart extends GridsterItem {
  chartId: number;
}

export interface Dashboard {
  id?: number
  name?: string
  grid?: EChart[]
}

export interface DashboardGrid {
  x: number
  y: number
  rows: number
  cols: number
  chartId: number
}

export interface CreateDashboardReq {
  name: string
  grid: DashboardGrid[]
}

@Injectable({
  providedIn: 'root'
})
export class DashboardService {
  private base = '/api/dashboards';

  constructor(private http: HttpClient) {
  }

  all(): Observable<Dashboard[]> {
    return this.http.get<any>(`${this.base}`);
  }

  id(id: number): Observable<Dashboard> {
    return this.http.get<any>(`${this.base}/${id}`);
  }

  create(req: CreateDashboardReq): Observable<Dashboard> {
    return this.http.post<any>(`${this.base}`, req);
  }

  delete(id: number): Observable<any> {
    return this.http.delete(`${this.base}/${id}`);
  }
}
