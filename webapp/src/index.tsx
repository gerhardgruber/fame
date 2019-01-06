import * as React from 'react';
import * as ReactDOM from 'react-dom';
import Root from './pages/Root';
import UiStore from './stores/UiStore'
import { when } from 'mobx'
import 'antd/dist/antd.css';
import '../static/ant-theme-vars.less'

const uiStore = UiStore.getInstance();

function bootstrapped() : boolean {
    console.log( "check if bootstrapped..." );
    return uiStore.bootstrapped;
}

function start() : void {
    const element = (
        <Root />
    )
    ReactDOM.render( element, document.getElementById('root'));
}

when( bootstrapped, start );
