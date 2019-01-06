import {observable} from "mobx";
import AbstractNotification from "../core/AbstractNotification";
import ErrorNotification from "../core/ErrorNotification";
import InfoNotification from "../core/InfoNotification";
import {notification} from "antd";
import UiStore from "./UiStore";

const antdNotification = notification;

class NotificationStore {
    @observable public entities: AbstractNotification[];

    constructor() {
        this.entities = [];
    }

    public error( message: string ) {
        this.addNotification( new ErrorNotification( message ) );
    }

    public info( message: string ) {
        this.addNotification( new InfoNotification( message ) );
    }

    private addNotification( notification: AbstractNotification ) {
        antdNotification[ notification.getType() ]( {
            message: UiStore.getInstance( ).T( "ERROR" ),
            description: notification.message
        })
        this.entities.push( notification );
    }
}

export default NotificationStore;