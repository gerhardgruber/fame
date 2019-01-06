class Position {
    public lat: number;
    public lng: number;
    public createdAt: string;

    constructor(data: any) {
        this.lat = Number(data.Latitude);
        this.lng = Number(data.Longitude);
        this.createdAt = data.CreatedAt;
    }
}

export default Position;