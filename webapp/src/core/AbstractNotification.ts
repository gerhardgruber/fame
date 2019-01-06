import { observable } from 'mobx'

abstract class AbstractNotification {
    @observable public message: string;

    constructor( message: string ) {
        this.message = message;
    }

    public abstract getType() : string;
}

export default AbstractNotification;