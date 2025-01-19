import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Dataset} from './dataset.service';

export enum SourceType {
  POSTGRES = 'postgres',
  CSV = 'csv',
}

export interface Source {
  id?: number
  name?: string
  type?: SourceType
  datasets?: Dataset[]
}

@Injectable({
  providedIn: 'root'
})
export class SourceService {
  private base = '/api/sources';

  constructor(private http: HttpClient) {
  }

  all(): Observable<Source[]> {
    return this.http.get<any>(`${this.base}`);
  }

  id(id: number): Observable<Source> {
    return this.http.get<any>(`${this.base}/${id}`);
  }

  create(name: string, type: string, resource: string): Observable<Source> {
    return this.http.post<any>(`${this.base}`, {name, type, resource});
  }

  connect(uri: string): Observable<Source> {
    return this.http.post<any>(`${this.base}/discovery`, {type: SourceType.POSTGRES, uri});
  }

  upload(files: File[]): Observable<Source> {
    const form = new FormData();
    form.append("type", SourceType.CSV);

    for (const file of files) {
      form.append('files', file, file.name);
    }

    return this.http.post<any>(`${this.base}/discovery`, form);
  }

  delete(id: number): Observable<any> {
    return this.http.delete(`${this.base}/${id}`);
  }
}
