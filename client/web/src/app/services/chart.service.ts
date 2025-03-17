import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {EChartsOption} from 'echarts';

export interface Chart {
  id?: number
  datasetId?: string
  name?: string
  type?: string
  dimensions?: string[]
  metrics?: string[]
  filters?: string[]
}

export interface ChartSchema {
  type: string
  schema: ChartSchemaRules
  example: any
}

export interface ChartSchemaRules {
  dimensions: FieldRule
  metrics: FieldRule
  filters?: FieldRule
}

export interface FieldRule {
  min: number
  max: number
  values?: string[]
}

export interface Filter {
  dimension: string;
  operator: string;
  value: string;
}

export interface CreateChartReq {
  datasetId: number
  name: string
  type: string
  dimensions: string[]
  metrics: string[]
  filters?: string[]
}

export interface ValidateChartReq {
  datasetId: number
  name: string
  type: string
  dimensions: string[]
  metrics: string[]
  filters: string[]
}

export interface ValidateChartRsp {
  valid: boolean
  options: EChartsOption
}

@Injectable({
  providedIn: 'root'
})
export class ChartService {
  private base = '/api/charts';
  private baseTypes = '/api/chart-types';

  constructor(private http: HttpClient) {
  }

  all(): Observable<Chart[]> {
    return this.http.get<any>(`${this.base}`);
  }

  id(id: number): Observable<Chart> {
    return this.http.get<any>(`${this.base}/${id}`);
  }

  create(req: CreateChartReq): Observable<Chart> {
    return this.http.post<any>(`${this.base}`, req);
  }

  delete(id: number): Observable<any> {
    return this.http.delete(`${this.base}/${id}`);
  }

  types(): Observable<string[]> {
    return this.http.get<any>(`${this.baseTypes}`);
  }

  type(type: string): Observable<ChartSchema> {
    return this.http.get<any>(`${this.baseTypes}/${type}`);
  }

  validate(req: ValidateChartReq): Observable<ValidateChartRsp> {
    return this.http.post<any>(`${this.base}/validate`, req);
  }

  run(id: number): Observable<ValidateChartRsp> {
    return this.http.post<any>(`${this.base}/${id}/data`, null);
  }
}
