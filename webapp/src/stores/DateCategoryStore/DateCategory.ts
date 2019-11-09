import {isNil, map} from 'lodash';
import User from '../User';
import { Address } from '../Address';
import UiStore from '../UiStore';
import { observable } from 'mobx';

export class DateCategory {
  public ID : number;

  public Name : string;

  constructor( data ) {
    this.setData(data);
  }

  public setData(data) {
    if (!isNil(data.ID)) {
      this.ID = data.ID;
    }
    this.Name = data.Name;

    return this;
  }

  public getData( ) : object {
    return {
      Name: this.Name
    };
  }
}