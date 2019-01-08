import {isNil} from 'lodash';

class Operation {
  public ID : number;

  public Title : string;

  public FirstName : string;

  public LastName : string;

  // TODO: Set Company
  constructor( data ) {
    this.setData(data);
  }

  public setData(data) {
    if (!isNil(data.ID)) {
      this.ID = data.ID;
    }
    this.Title = data.Title;
    this.FirstName = data.FirstName;
    this.LastName = data.LastName;

    return this;
  }

  public getData( ) : object {
    return {
      Title: this.Title,
      FirstName: this.FirstName,
      LastName: this.LastName,
    };
  }
}

export default Operation;