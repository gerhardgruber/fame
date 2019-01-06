import { observable } from 'mobx'
import AbstractNotification from './AbstractNotification'

class InfoNotification extends AbstractNotification {
    public getType() : string {
        return "info";
    }
}

export default InfoNotification;