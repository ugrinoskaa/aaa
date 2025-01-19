import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

export interface Dataset {
  id?: number
  name?: string
  sourceId?: number
  schema?: string
  table?: string
  columns?: string[]
  dimensions?: string[];
  metrics?: string[]
}

export interface Column {
  name: string
  type: string
}

export interface CreateDatasetReq {
  name: string
  sourceId: string
  sourceSchema: string
  sourceTable: string
}

@Injectable({
  providedIn: 'root'
})
export class DatasetService {
  private base = '/api/datasets';

  constructor(private http: HttpClient) {
  }

  all(): Observable<Dataset[]> {
    return this.http.get<any>(`${this.base}`);
  }

  id(id: number): Observable<Dataset> {
    return this.http.get<any>(`${this.base}/${id}`);
  }

  create(req: CreateDatasetReq): Observable<Dataset> {
    return this.http.post<any>(`${this.base}`, req);
  }

  delete(id: number): Observable<any> {
    return this.http.delete(`${this.base}/${id}`);
  }
}
