import { observable } from 'mobx'
import AbstractNotification from './AbstractNotification'

class ErrorNotification extends AbstractNotification {
    public getType() : string {
        return "error";
    }
}

export default ErrorNotification;