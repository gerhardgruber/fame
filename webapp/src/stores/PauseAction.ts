export enum PauseType {
  TrainingPause = 0,
  OperationPause = 1
}

export class PauseAction {
  public ID: number;

  public UserID: number;

  public Type: PauseType;

  public StartTime: Date;

  public EndTime: Date;

  constructor(data: Record<string, any>) {
    this.setData(data);
  }

  public setData(data: Record<string, any>) {
    this.ID = data.ID;
    this.UserID = data.UserID;
    this.Type = data.Type;
    this.StartTime = data.StartTime
    this.EndTime = data.EndTime;
  }

  public getData(): Record<string,any> {
    return {
      ID: this.ID,
      UserID: this.UserID,
      Type: this.Type,
      StartTime: this.StartTime,
      EndTime: this.EndTime
    };
  }
}