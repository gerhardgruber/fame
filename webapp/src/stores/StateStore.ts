import {action, observable} from 'mobx';
import Api from '../core/Api';
import NotificationStore from './NotificationStore';

class StateStore {
    private static instance: StateStore;

    public notifications: NotificationStore;

    constructor() {
        this.notifications = new NotificationStore();
    }

    public static getInstance() : StateStore {
        if ( StateStore.instance == null ) {
            StateStore.instance = new StateStore();
        }

        return StateStore.instance;
    }

    public setStatePromise(instance, content) : Promise<any> {
        return new Promise(function(resolve, reject) {
            instance.setState(content, resolve);
        });
    }

    @action public startLoading() {
        /* TODO: Set state, that application has started loading... */
    }

    @action public endLoading() {
        /* TODO: Set state, that application has ended loading... */
    }

    @action public progressAfterLoading() {
        /* TODO: Set state, that application is now afterloading... */
    }
}

export default StateStore;