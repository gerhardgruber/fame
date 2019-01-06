import * as React from 'react';
import DevTools from 'mobx-react-devtools';
import {inject, observer} from 'mobx-react';
import {Layout} from 'antd';
import FameMenu from './FameMenu';
import Login from './Login';
import { withRouter } from 'react-router';
const { Content } = Layout;
import { Route, Router } from 'react-router';
import { createBrowserHistory } from 'history';
import queryString from 'query-string'

@inject( "uiStore" )
@observer
class App extends React.Component<any,any> {
    state = {
        sidebarCollapsed: false,
    }

    constructor( props ) {
        super( props )
    }

    onCollapse = (collapsed) => {
        this.setState({ sidebarCollapsed: collapsed })
    }

    renderMenu() {
        return (
            <FameMenu
                collapsed={this.state.sidebarCollapsed}
                onCollapse={this.onCollapse}
                location={this.props.location}
                />
        );
    }

    render() {
        const params = queryString.parse( location.search );

        if ( this.props.uiStore.loggedIn ) {
            return (
                <Layout>
                    {this.renderMenu()}
                    <Layout style={this.state.sidebarCollapsed ?
                            { marginLeft: 64 } : { marginLeft: 200}} >
                        <Content style={{ height: '100%' }}>
                            {this.props.children}
                        </Content>
                    </Layout>
                    {<DevTools />}
                </Layout>
            );

        } else {
            return <Login />;
        }
    }
}

export default withRouter(App);