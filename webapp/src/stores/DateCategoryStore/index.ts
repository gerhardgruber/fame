import { DateCategory } from "./DateCategory";
import { map } from 'lodash';
import Api from "../../core/Api";
import { observable } from 'mobx';
import User from "../User";
import UiStore from "../UiStore";

export default class DateCategoryStore {
  private static instance: DateCategoryStore;

  @observable public dateCategories: DateCategory[];

  constructor() {
    DateCategoryStore.instance = this;
  }

  public static getInstance(): DateCategoryStore {
    if ( DateCategoryStore.instance ) {
      return DateCategoryStore.instance;
    }

    DateCategoryStore.instance = new DateCategoryStore();
    return DateCategoryStore.instance;
  }

  public loadDateCategory(id: number): Promise<DateCategory> {
    return Api.GET( `/date_categories/${id}`, {} ).then( ( response ) => {
      const dateCategory = new DateCategory( response.data.data.DateCategory );

      return dateCategory;
    } );
  }

  public loadDateCategories(): Promise<DateCategory[]> {
    return Api.GET( '/date_categories', {} ).then( ( response ) => {
      this.dateCategories = map( response.data.data.DateCategories, ( o: any ) => {
        return new DateCategory( o );
      } );

      return this.dateCategories;
    } );
  }

  public saveDateCategory(dateCategory: DateCategory): Promise<void> {
    return Api.POST(`/date_categories/${dateCategory.ID}`, dateCategory.getData());
  }

  public createDateCategory(dateCategory: DateCategory): Promise<void> {
    return Api.POST(`/date_categories`, dateCategory.getData());
  }

  public deleteDateCategory(dateCategory: DateCategory): Promise<void> {
    return Api.POST(`/date_categories/${dateCategory.ID}/delete`, {});
  }
}