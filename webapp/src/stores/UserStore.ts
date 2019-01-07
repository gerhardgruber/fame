import User from "./User";
import { map } from 'lodash';
import Api from "../core/Api";
import { observable } from 'mobx';
import {isNil} from 'lodash';

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

  public loadUser(id: number): Promise<User> {
    return Api.GET( `/users/${id}`, {} ).then( ( response ) => {
      return new User( response.data.data.User );
    } );
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