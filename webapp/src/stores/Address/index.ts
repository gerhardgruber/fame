export class Address {
  public Street: string;

  public Number: string;

  public Postcode: string;

  public City: string;

  public Country: string;

  constructor(data: any) {
    this.Street = data.Street;
    this.Number = data.Number;
    this.Postcode = data.Postcode;
    this.City = data.City;
    this.Country = data.Country;
  }

  public toString(): string {
    return [
      this.Street,
      this.Number,
      this.Postcode,
      this.City,
      this.Country
    ].join( " " ).replace( / +/g, " " );
  }
}