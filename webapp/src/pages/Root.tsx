import * as React from 'react';
import { Route, Router } from 'react-router';
import { observer, Provider } from 'mobx-react'
import { createBrowserHistory } from 'history';
import { CookiesProvider } from 'react-cookie';

import { LocaleProvider } from 'antd';
import enUS from 'antd/lib/locale-provider/en_US';

import App from './App';
import Users from './Users';
import UiStore from '../stores/UiStore';
import ChangePassword from './ChangePassword';
import Operations from './Operations';
import { EditUser } from './Users/edit';
import UserStore from '../stores/UserStore';

const uiStore = UiStore.getInstance();
const userStore = UserStore.getInstance();

@observer
export class Root extends React.Component {
    render() {
        return (
            <CookiesProvider>
                <LocaleProvider locale={enUS}>
                    <Provider
                        uiStore={ uiStore }
                        >
                        <Router history={ createBrowserHistory() }>
                            <App>
                                <Route exact path="/" component={Operations} />
                                <Route exact path="/operations" component={Operations} />
                                <Route path="/changePassword" component={ChangePassword} />
                                <Route path="/users" exact component={Users} />
                                <Route path="/users/new" exact component={EditUser} />
                                <Route path="/users/:id" render={( { match } ) => {
                                    if ( match.params.id !== "new" ) {
                                        return <EditUser userID={parseInt(match.params.id)} />
                                    } else {
                                        return null;
                                    }
                                } } />
                            </App>
                        </Router>
                    </Provider>
                </LocaleProvider>
            </CookiesProvider>
        );
    }
}

export default Root;
