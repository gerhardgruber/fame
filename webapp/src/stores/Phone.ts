import { observable } from 'mobx'
import Api from '../core/Api'
import User from './User'

class Phone {
  public ID : number;
  public PhoneNumber : string;
  public Device : string;
  public Passcode : string;
  public User : User;
  public UserID : number;

  constructor( data ) {
    this.setData( data )
  }

  public setData( data ) : Phone {
    if ( data.ID ) {
      this.ID = data.ID;
    }

    if ( data.PhoneNumber ) {
      this.PhoneNumber = data.PhoneNumber;
    }

    if ( data.Device ) {
      this.Device = data.Device;
    }

    if ( data.Passcode ) {
      this.Passcode = data.PublicKey;
    }

    if ( data.User ) {
      this.User = new User( data.User );
    }

    if ( data.UserID ) {
      this.UserID = data.UserID;
    }

    return this;
  }

  public getData( ) : object {
    return {
      ID: this.ID,
      PhoneNumber: this.PhoneNumber,
      Device: this.Device,
      UserID: this.UserID
    };
  }
}

export default Phone;