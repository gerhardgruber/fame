class User {
  public ID : number;

  public EMail : string;

  public Name : string;

  public FirstName : string;

  public LastName : string;

  public Lang : string;

  public PW : string;

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
}

export default User;