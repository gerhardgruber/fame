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

const uiStore = UiStore.getInstance();

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
                                <Route exact path="/" component={Users} />
                                <Route path="/changePassword" component={ChangePassword} />
                                <Route path="/users" component={Users} />
                            </App>
                        </Router>
                    </Provider>
                </LocaleProvider>
            </CookiesProvider>
        );
    }
}

export default Root;
