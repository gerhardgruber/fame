import User from "./User";
import { map } from 'lodash';
import Api from "../core/Api";
import { observable } from 'mobx';

export default class UserStore {
  private static instance: UserStore;

  @observable public users: User[];

  constructor() {
    UserStore.instance = this;
  }

  public static getInstance(): UserStore {
    if ( UserStore.instance ) {
      return UserStore.instance;
    }

    UserStore.instance = new UserStore();
    return UserStore.instance;
  }

  public loadUsers(): Promise<User[]> {
    return Api.GET( '/users', {} ).then( ( response ) => {
      this.users = map( response.data.data.Users, ( u: any ) => {
        return new User( u );
      } );

      return this.users;
    } );
  }
}