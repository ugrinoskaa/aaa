import {Injectable} from '@angular/core';
import {SourceService} from './source.service';
import {DatasetService} from "./dataset.service";

@Injectable({
  providedIn: 'root'
})
export class APIService {

  constructor(
    private _sources: SourceService,
    private _datasets: DatasetService
  ) {
  }

  sources(): SourceService {
    return this._sources;
  }

  datasets(): DatasetService {
    return this._datasets;
  }
}
