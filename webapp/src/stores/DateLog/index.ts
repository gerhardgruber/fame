import { observable } from "mobx";
import Api from "../../core/Api";

export class DateLog {
  public ID: number;

  @observable public UserID: number;
  @observable public DateID: number;

  @observable public FromTime: Date;
  @observable public UntilTime: Date;
  @observable public Present: boolean;

  @observable public Comment: string;

  constructor(data: any) {
    this.ID = data.ID;

    this.UserID = data.UserID;
    this.DateID = data.DateID;

    this.FromTime = new Date( data.FromTime );
    this.UntilTime = new Date( data.UntilTime );
    this.Present = data.Present;

    this.Comment = data.Comment;
  }

  public create() {
    Api.POST( "/date_logs", {
      UserID: this.UserID,
      DateID: this.DateID,

      FromTime: this.FromTime,
      UntilTime: this.UntilTime,
      Present: this.Present,

      Comment: this.Comment
    } ).then( ( response ) => {
      this.ID = response.data.data.DateLog.ID;
    } );
  }

  public save() {
    Api.POST( `/date_logs/${this.ID}`, {
      UserID: this.UserID,
      DateID: this.DateID,

      FromTime: this.FromTime,
      UntilTime: this.UntilTime,
      Present: this.Present,

      Comment: this.Comment
    } );
  }
}