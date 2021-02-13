import User from "./User";
import { map } from 'lodash';
import Api from "../core/Api";
import { observable } from 'mobx';
import {isNil} from 'lodash';

export default class UserStore {
  private static instance: UserStore;

  @observable public users: User[];

  @observable public stati: any[];

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
      const userData = response.data.data.User
      userData.TrainingPause = response.data.data.TrainingPause;
      userData.OperationPause = response.data.data.OperationPause;
      return new User( userData );
    } );
  }

  public loadUsers(): Promise<User[]> {
    return Api.GET( '/users', {} ).then( ( response ) => {
      this.stati = response.data.data.Stati;
      this.users = map( response.data.data.Users, ( u: any ) => {
        return new User( u );
      } );

      return this.users;
    } );
  }

  public saveUser(user: User): Promise<void> {
    return Api.POST(`/users/${user.ID}`, user.getData());
  }

  public createUser(user: User): Promise<void> {
    return Api.POST(`/users`, user.getData());
  }

  public deleteUser(user: User): Promise<void> {
    return Api.POST(`/users/${user.ID}/delete`, {});
  }
}