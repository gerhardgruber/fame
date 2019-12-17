import {isNil, map} from 'lodash';
import User from '../User';
import { Address } from '../Address';
import UiStore from '../UiStore';
import { observable } from 'mobx';
import { DateCategory } from '../DateCategoryStore/DateCategory';

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

  @observable public Feedbacks: DateFeedback[];

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
        return a.User.LastName.localeCompare(a.User.LastName)
      }
    })
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
      _feedbacks[f.UserID].Feedback = f.Feedback;
      _feedbacks[f.UserID].UpdatedAt = new Date(f.UpdatedAt);
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

    return this;
  }

  public getData( ) : object {
    return {
      Title: this.Title,
      Description: this.Description,
      StartTime: this.StartTime,
      EndTime: this.EndTime,
      Location: this.Location,
      CategoryID: this.CategoryID
    };
  }
}