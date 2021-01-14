import {isNil} from 'lodash';
import { observable } from 'mobx';
import Api from '../core/Api';
import { PauseAction, PauseType } from './PauseAction';

export enum RightType {
  STANDARD = 0,
  ADMIN = 1,
}
class User {
  public ID : number;

  public EMail : string;

  public Name : string;

  public FirstName : string;

  public LastName : string;

  public Lang : string;

  public PW : string;

  public RightType: RightType;

  @observable public TrainingPause: PauseAction;

  @observable public OperationPause: PauseAction;

  // TODO: Set Company
  constructor( data ) {
    this.setData(data);
  }

  public setData(data) {
    if (!isNil(data.ID)) {
      this.ID = data.ID;
    }
    this.EMail = data.EMail;
    this.Name = data.Name;
    this.FirstName = data.FirstName;
    this.LastName = data.LastName;
    this.Lang = data.Lang;
    this.PW = data.PW;
    this.RightType = data.RightType;
    if (!isNil(data.TrainingPause)) {
      this.TrainingPause = new PauseAction(data.TrainingPause);
    }
    if (!isNil(data.OperationPause)) {
      this.OperationPause = new PauseAction(data.OperationPause);
    }

    return this;
  }

  public getData( ) : object {
    return {
      Name: this.Name,
      EMail: this.EMail,
      FirstName: this.FirstName,
      LastName: this.LastName,
      Lang: this.Lang,
      PW: this.PW,
      RightType: this.RightType
    };
  }

  public getFullName(): string {
    if ( this.FirstName && this.LastName) {
      return `${this.FirstName} ${this.LastName}`;
    } else if ( this.FirstName ) {
      return this.FirstName;
    } else {
      return this.LastName;
    }
  }

  public startPause(type: PauseType, startTime: Date) {
    Api.POST(`/users/${this.ID}/start_pause`, {
      Type: type,
      StartTime: startTime
    }).then((response) => {
      if (type === PauseType.TrainingPause) {
        this.TrainingPause = new PauseAction(response.data.data.PauseAction);
      } else if (type === PauseType.OperationPause) {
        this.OperationPause = new PauseAction(response.data.data.PauseAction);
      }
    });
  }

  public stopPause(type: PauseType, endTime: Date) {
    Api.POST(`/users/${this.ID}/stop_pause`, {
      Type: type,
      EndTime: endTime
    }).then((response) => {
      if (type === PauseType.TrainingPause) {
        this.TrainingPause = new PauseAction(response.data.data.PauseAction);
      } else if (type === PauseType.OperationPause) {
        this.OperationPause = new PauseAction(response.data.data.PauseAction);
      }
    });
  }
}

export default User;