import {Injectable} from '@angular/core';
import {SourceService} from './source.service';
import {DatasetService} from "./dataset.service";
import {ChartService} from './chart.service';

@Injectable({
  providedIn: 'root'
})
export class APIService {

  constructor(
    private _sources: SourceService,
    private _datasets: DatasetService,
    private _charts: ChartService,
  ) {
  }

  sources(): SourceService {
    return this._sources;
  }

  datasets(): DatasetService {
    return this._datasets;
  }

  charts(): ChartService {
    return this._charts;
  }
}
