import {Injectable} from '@angular/core';
import {SourceService} from './source.service';

@Injectable({
  providedIn: 'root'
})
export class APIService {

  constructor(
    private _sources: SourceService
  ) {
  }

  sources(): SourceService {
    return this._sources;
  }
}
