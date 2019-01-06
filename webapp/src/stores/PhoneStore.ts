import {action, observable, runInAction, computed} from 'mobx';
import Api from '../core/Api';
import Phone from './Phone';

class PhoneStore {
    private static instance: PhoneStore;

    @observable public currentPhone : Phone;

    public phones : Phone[]

    constructor() {
    }

    public static getInstance() : PhoneStore {
        if ( PhoneStore.instance == null ) {
            PhoneStore.instance = new PhoneStore();
        }

        return PhoneStore.instance;
    }

    public loadPhone( phoneID : number ) : PromiseLike<Phone> {
        return Api.GET( "/mobile_phones/" + phoneID, {} ).then( ( response ) => {
            const phone = new Phone( response.data.data.MobilePhone )
            this.currentPhone = phone
            return phone
        } )
    }

    public newPhone( phone : Phone ) : Promise<any> {
        return Api.POST( "/mobile_phones", phone.getData() ).then( ( response ) => {
            this.currentPhone = new Phone( response.data.data.MobilePhone )
        } )
    }

    public savePhone( phone : Phone ) : Promise<any> {
        let url = '/mobile_phones'
        if ( phone.ID ) {
            url += '/' + phone.ID
        }

        return Api.POST( url, phone.getData( ) ).then( ( response ) => {
            phone.setData( response.data.data.MobilePhone )
        } )
    }

    public loadAll( ) : PromiseLike<Array<Phone>> {
        return Api.GET( '/mobile_phones', {} ).then( ( response ) => {
            const phones = response.data.data.rows.map((phoneData) => new Phone(phoneData));
            this.phones = phones;
            return phones
        } )
    }
}

export default PhoneStore;