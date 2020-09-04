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
import { RightType } from '../stores/User';
import { EditOperation } from './Operations/edit';
import Dates from './Dates';
import { EditDate } from './Dates/edit';
import DateCategories from './DateCategory';
import { EditDateCategory } from './DateCategory/edit';
import Statistics from './Statistics';

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
                                <Route exact path="/" component={Dates} />
                                <Route exact path="/date_categories" component={DateCategories} />
                                <Route path="/date_categories/new" exact component={EditDateCategory} />
                                <Route path="/date_categories/:id" render={( { match } ) => {
                                    if ( match.params.id !== "new" ) {
                                        return <EditDateCategory dateCategoryID={parseInt(match.params.id)} />
                                    } else {
                                        return null;
                                    }
                                } } />
                                <Route exact path="/dates" component={Dates} />
                                <Route path="/dates/new" exact component={EditDate} />
                                <Route path="/dates/:id" render={( { match } ) => {
                                    if ( match.params.id !== "new" ) {
                                        return <EditDate dateID={parseInt(match.params.id)} />
                                    } else {
                                        return null;
                                    }
                                } } />
                                <Route exact path="/operations" component={Operations} />
                                <Route path="/operations/new" exact component={EditOperation} />
                                <Route path="/operations/:id" render={( { match } ) => {
                                    if ( match.params.id !== "new" ) {
                                        return <EditOperation operationID={parseInt(match.params.id)} />
                                    } else {
                                        return null;
                                    }
                                } } />
                                <Route path="/changePassword" component={ChangePassword} />
                                <Route path="/statistics" component={Statistics} />
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
