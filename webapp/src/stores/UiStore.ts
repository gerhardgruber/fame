import {action, observable, runInAction, computed, observe, intercept} from 'mobx';
import Api from '../core/Api';
import User, { RightType } from './User';
import DateTimeComponent from '../components/DateTimeComponent';
import moment from 'moment';

class UiStore {
    private static instance: UiStore;

    private translationData: Map<string, string>

    @observable public loggedIn: boolean;

    @observable public loginTime: Date;

    @observable public dateTypes: Record<string,number>;

    @observable public dateFeedbackTypes: Record<string,number>;

    @observable public bootstrapped: boolean;

    @observable public currentUser : User;

    public formItemLayout = {
        labelCol: {
            xs: { span: 24 },
            sm: { span: 5 },
        },
        wrapperCol: {
            xs: { span: 24 },
            sm: { span: 16 },
        },
    };

    constructor() {
        this.loggedIn = false;
        this.bootstrapped = false;

        this.bootstrap();
        Api.callback.translate = this.T.bind( this );
    }

    public static getInstance() : UiStore {
        if ( UiStore.instance == null ) {
            UiStore.instance = new UiStore();
        }

        return UiStore.instance;
    }

    public bootstrap() {
        this.loadStatics().then( () => {
            this.bootstrapped = true;
        } );
    }

    private loadStatics() : Promise<any> {
        return Api.GET( "/statics", {} ).then( ( response ) => {
            this.translationData = response.data.data.translation_data || new Map();

            if (response.data.data.logged_in) {
                this.loginTime = new Date(response.data.data.login_time);

                Api.GET( "/users/" + response.data.data.UserID, {} ).then( ( response ) => {
                    this.currentUser = new User( response.data.data.User )
                } )
            }

            this.loggedIn = response.data.data.logged_in;

            this.dateTypes = response.data.data.date_types;
            this.dateFeedbackTypes = response.data.data.date_feedback_types;
        } );
    }

    public T(translationName: string) : string {
        var translation = this.translationData[translationName];
        if (translation !== undefined) {
            return translation;
        }

        console.log("ERROR: Translation for " + translationName + " not found!");
        return translationName;
    }

    public DateTime(datetime: any): JSX.Element {
        return DateTimeComponent.makeComponent(datetime, this.T('DATE_TIME_FORMAT'));
    }

    public DateTimeString(datetime: any) {
        return moment(datetime).format(this.T('DATE_TIME_FORMAT'));
    }

    public setSession( sessionKey : string ) {
        window.sessionStorage.setItem( "session", sessionKey);
        return this.loadStatics();
    }

    public login( username : string, password : string ) : Promise<any> {
        return Api.POST( "/authentication/login", { Name: username, PW: password } ).then( ( response ) => {
            this.currentUser = new User(response.data.data.user);
            return this.setSession( response.data.data.session );
        } )
    }

    public logout() : Promise<any> {
        this.loggedIn = false;

        return Api.POST( '/authentication/logout', {} ).then( ( response ) => {
            window.sessionStorage.removeItem( "session" );
            this.loadStatics();
        } );
    }

    public register( u : User ) : PromiseLike<User> {
        return Api.POST( '/authentication/register', u.getData() ).then( ( response ) => {
            window.sessionStorage.setItem( "session", response.data.data.session );
            this.currentUser = new User(response.data.data.user);
            return this.loadStatics().then(() => u.setData(response.data.data.user));
        } );
    }

    public updateUser(u: User): PromiseLike<User> {
        return Api.POST( `/users/${u.ID}`, u.getData()).then((response) => {
            u.setData(response.data.data.user);
            return this.loadStatics().then(() => u);
        });
    }

    @computed get username() : string {
        if ( this.currentUser ) {
            return this.currentUser.Name
        }
        return ""
    }

    public changePassword(oldPassword: string, newPassword: string): Promise<void> {
        return Api.POST( `/users/${this.currentUser.ID}/password`, {
            OldPassword: oldPassword,
            NewPassword: newPassword,
        });
    }

    public isAdmin(): boolean {
        if (this.currentUser) {
            return this.currentUser.RightType === RightType.ADMIN
        }

        return false;
    }
}

export default UiStore;