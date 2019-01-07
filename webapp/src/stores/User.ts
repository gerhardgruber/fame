export enum RightType {
  STANDARD = 0,
  ADMIN = 1,
}
class User {
  public ID : number;

  public EMail : string;

  public Name : string;

  public FirstName : string;

  public LastName : string;

  public Lang : string;

  public PW : string;

  public RightType: RightType;

  // TODO: Set Company
  constructor( data ) {
    this.setData(data);
  }

  public setData(data) {
    this.ID = data.ID;
    this.EMail = data.EMail;
    this.Name = data.Name;
    this.FirstName = data.FirstName;
    this.LastName = data.LastName;
    this.Lang = data.Lang;
    this.PW = data.PW;
    this.RightType = data.RightType;

    return this;
  }

  public getData( ) : object {
    return {
      Name: this.Name,
      EMail: this.EMail,
      FirstName: this.FirstName,
      LastName: this.LastName,
      Lang: this.Lang,
      PW: this.PW
    };
  }

  public getFullName(): string {
    if ( this.FirstName && this.LastName) {
      return `${this.FirstName} ${this.LastName}`;
    } else if ( this.FirstName ) {
      return this.FirstName;
    } else {
      return this.LastName;
    }
  }
}

export default User;