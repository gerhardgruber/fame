import {isNil, map} from 'lodash';
import User from '../User';
import { Address } from '../Address';
import UiStore from '../UiStore';
import { observable } from 'mobx';
import { DateCategory } from '../DateCategoryStore/DateCategory';
import find from 'lodash/find';
import { DateLog } from '../DateLog';

export interface IFeedbackStat {
  Feedback: number;
  count: number;
}

export class DateFeedback {
  User: User;

  @observable public Feedback: number;

  @observable public UpdatedAt: Date;

  constructor(u: User, feedback: number, updatedAt: string) {
    this.User = u;
    this.Feedback = feedback;
    if (updatedAt) {
      this.UpdatedAt = new Date(updatedAt);
    }
  }
}

export class DateModel {
  public ID : number;

  public Title : string;

  public Description : string;

  public LocationAddress: Address;
  public Location : string;

  public CreatedBy: User;

  public CategoryID: number;
  public Category: DateCategory;

  public StartTime: Date;
  public EndTime: Date;

  public Closed: boolean;

  @observable public Feedbacks: DateFeedback[];

  private dateFeedbacks: any;

  public DateLogs: DateLog[];
  public DateLogsByUserID: Record<number, DateLog>;

  constructor( data ) {
    this.setData(data);
    this.Feedbacks = [];
  }

  public get orderedFeedbacks(): DateFeedback[] {
    return this.Feedbacks.sort((a, b) => {
      if (a.Feedback < b.Feedback) {
        return -1;
      } else if ( a.Feedback > b.Feedback) {
        return 1;
      } else {
        return 0;
      }
    })
  }

  public get orderedFeedbacksWithHeaders(): Array<DateFeedback | IFeedbackStat> {
    const groupedFeedbacks: Record<number,DateFeedback[]> = {};
    this.Feedbacks.forEach( (fb) => {
      if (!groupedFeedbacks[fb.Feedback]) {
        groupedFeedbacks[fb.Feedback] = [];
      }
      groupedFeedbacks[fb.Feedback].push(fb);
    });

    const ret: Array<DateFeedback | IFeedbackStat> = [];

    Object.keys(groupedFeedbacks).sort().forEach((type) => {
      const fbType = Number(type)
      const feedbacks = groupedFeedbacks[fbType].sort((a, b) => {
        return a.User.LastName.localeCompare(b.User.LastName)
      })
      ret.push({
        Feedback: fbType,
        count: feedbacks.length
      }, ...feedbacks);
    })

    return ret;
  }

  public getMyFeedback(): DateFeedback {
    const uiStore = UiStore.getInstance();

    return find(this.dateFeedbacks, (fb) => {
      return fb.UserID === uiStore.currentUser.ID;
    });
  }

  public calculateFeedbacks(feedbacks: any[], users: User[]) {
    const uiStore = UiStore.getInstance();

    const _feedbacks: Record<number,DateFeedback> = {};

    users.forEach( (u) => {
      _feedbacks[u.ID] = new DateFeedback(
        u,
        uiStore.dateFeedbackTypes["Unknown"],
        null
      );
    });

    feedbacks.forEach( (f) => {
      /* We have to check if _feedbacks[f.UserID] exists, because if a user was deleted, it will not exist */
      if (_feedbacks[f.UserID]) {
        _feedbacks[f.UserID].Feedback = f.Feedback;
        _feedbacks[f.UserID].UpdatedAt = new Date(f.UpdatedAt);
      }
    });

    this.Feedbacks = map(_feedbacks);
  }

  public setData(data) {
    if (!isNil(data.ID)) {
      this.ID = data.ID;
    }
    this.Title = data.Title;
    this.Description = data.Description;
    this.StartTime = new Date(data.StartTime);
    this.EndTime = new Date(data.EndTime);
    if (typeof data.Location === "string") {
      this.Location = data.Location
    } else if (!isNil(data.Location)) {
      this.LocationAddress = new Address(data.Location);
      this.Location = this.LocationAddress.toString();
    }
    this.CreatedBy = data.CreatedBy;
    this.CategoryID = data.CategoryID;
    if (data.Category) {
      this.Category = new DateCategory(data.Category);
    }
    if (data.DateFeedbacks) {
      this.dateFeedbacks = data.DateFeedbacks;
    } else {
      this.dateFeedbacks = [];
    }

    this.Closed = data.Closed;

    if (data.DateLogs) {
      this.DateLogsByUserID = {};
      this.DateLogs = map( data.DateLogs, ( dl ) => {
        const dateLog = new DateLog( dl )

        this.DateLogsByUserID[ dateLog.UserID ] = dateLog;

        return dateLog;
      } );
    } else {
      this.DateLogs = [];
      this.DateLogsByUserID = {};
    }

    return this;
  }

  public getData( ) : object {
    return {
      Title: this.Title,
      Description: this.Description,
      StartTime: this.StartTime,
      EndTime: this.EndTime,
      Location: this.Location,
      CategoryID: this.CategoryID,
      Closed: this.Closed
    };
  }

  public saveDateLog( dateLog: DateLog ) {
    if ( !dateLog.ID ) {
      this.DateLogsByUserID[ dateLog.UserID ] = dateLog;
      this.DateLogs.push( dateLog );
      dateLog.create();

    } else {
      dateLog.save();
    }
  }
}