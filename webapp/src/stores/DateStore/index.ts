import { DateModel, DateFeedback } from "./Date";
import { map } from 'lodash';
import Api from "../../core/Api";
import { observable } from 'mobx';
import User from "../User";
import UiStore from "../UiStore";

export default class DateStore {
  private static instance: DateStore;

  @observable public dates: DateModel[];

  constructor() {
    DateStore.instance = this;
  }

  public static getInstance(): DateStore {
    if ( DateStore.instance ) {
      return DateStore.instance;
    }

    DateStore.instance = new DateStore();
    return DateStore.instance;
  }

  public loadDate(id: number): Promise<DateModel> {
    return Api.GET( `/dates/${id}`, {} ).then( ( response ) => {
      const date = new DateModel( response.data.data.Date );

      const users = map(response.data.data.Users, (u) => {
        return new User(u);
      });

      date.calculateFeedbacks(response.data.data.Date.DateFeedbacks || [], users);

      return date;
    } );
  }

  public loadDates(loadPastDates: boolean): Promise<DateModel[]> {
    return Api.GET( '/dates', {
      loadPastDates
    } ).then( ( response ) => {
      this.dates = map( response.data.data.Dates, ( o: any ) => {
        return new DateModel( o );
      } );

      return this.dates;
    } );
  }

  public saveDate(date: DateModel): Promise<void> {
    return Api.POST(`/dates/${date.ID}`, date.getData());
  }

  public createDate(date: DateModel): Promise<void> {
    return Api.POST(`/dates`, date.getData());
  }

  public deleteDate(date: DateModel): Promise<void> {
    return Api.POST(`/dates/${date.ID}/delete`, {});
  }

  public sendFeedback(dateID: number, feedback: DateFeedback): Promise<DateFeedback> {
    return Api.POST( `/dates/${dateID}/feedback`, {
      UserID: UiStore.getInstance().currentUser.ID,
      Feedback: feedback.Feedback
    }).then( ( response ) => {
      return new DateFeedback(feedback.User, feedback.Feedback, response.data.data.DateFeedback.UpdatedAt);
    } );
  }
}