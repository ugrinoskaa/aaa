import {Injectable} from '@angular/core';
import {SourceService} from './source.service';
import {DatasetService} from "./dataset.service";
import {ChartService} from './chart.service';
import {DashboardService} from './dashboard.service';

@Injectable({
  providedIn: 'root'
})
export class APIService {

  constructor(
    private _sources: SourceService,
    private _datasets: DatasetService,
    private _charts: ChartService,
    private _dashboards: DashboardService,
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

  dashboards(): DashboardService {
    return this._dashboards;
  }
}
