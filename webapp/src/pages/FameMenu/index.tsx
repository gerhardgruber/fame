import * as React from 'react'
import { observer } from 'mobx-react';
import { Menu, Layout, Icon, Button } from 'antd'
import { Link } from 'react-router-dom'
import UiStore from '../../stores/UiStore'
import StateStore from '../../stores/StateStore'
import Scroll from 'react-scrollbar'
import Api from '../../core/Api';
import { computed } from 'mobx';
import { RightType } from '../../stores/User';

const { Header, Footer, Sider, Content } = Layout;
const { SubMenu } = Menu
@observer
class FameMenu extends React.Component<any,any> {
    private uiStore : UiStore;

    public state = {
        collapsed: false,
    }

    constructor( props ) {
        super( props )

        this.uiStore = UiStore.getInstance()
    }

    handleLogout = (e) => {
        this.uiStore.logout()
    }

    onCollapse = () => {
        var collapsed = !this.state.collapsed
        this.setState({ collapsed: !this.state.collapsed })
        this.props.onCollapse(collapsed)
    }

    currentMenu( ) : string[] {
        const key = this.props.location.pathname.split( "/" ).slice( 1, 2 );

        if ( key.length == 0 ) {
            return ["users"];
        }

        return key;
    }

    makeEntry(name:string, icon:string, onClick) : JSX.Element {
        return (
            <Menu.Item key={name}>
                <Link to={"/" + name} onClick={onClick}><Icon type={icon} /><span>{this.uiStore.T("MENU_" + name.toUpperCase())}</span></Link>
            </Menu.Item>
        )
    }

    render() {
        let userdescr = (<div style={{ display: "inline-block", textAlign: "left", marginLeft: 10, width: 136 }} >
                            <span style={{ fontSize: 16 }} ><b>{this.uiStore.username}</b></span><br/>
                            <span>{this.uiStore.DateTime(this.uiStore.loginTime)}</span>
                         </div>)

        let copyright = (<div style={{ display: "inline-block" }}>
                            <Icon type="copyright" /> {new Date().getFullYear()} Fame Inc.<br />
                        </div>)

        const entries = [];

        if ( this.uiStore.currentUser) {
            entries.push(this.makeEntry("operations", "notification", null));
            if (this.uiStore.currentUser.RightType === RightType.ADMIN) {
                entries.push(this.makeEntry("users", "user", null));
            }
            entries.push(this.makeEntry("changePassword", "setting", null));
        }

        return (
            <Sider theme={"light"} style={{ overflow: 'hidden', height: '100vh', position: 'fixed', left: 0 }}
                trigger={null}
                collapsible
                collapsed={this.state.collapsed}>
                <div style={{ textAlign: "center" }}>
                    <div style={{marginTop: '10px'}}>
                        <div>
                            <div style={{ display: "inline-block" }} >
                                <Icon type="user" style={{ fontSize: 40 }} />
                            </div>
                            {this.state.collapsed || userdescr}
                        </div>
                    </div>
                </div>
                <br/>
                <Scroll style={ {maxHeight: 'calc( 100vh - 200px )'}}>
                    <Menu style={{ overflow: 'auto', overflowX: 'hidden'}} mode="inline" defaultSelectedKeys={this.currentMenu()}>
                        {entries}
                        <Menu.Item>
                            <Link to="/" onClick={this.handleLogout}><Icon type={"logout"} /><span>{this.uiStore.T("MENU_LOGOUT")}</span></Link>
                        </Menu.Item>
                    </Menu>
                </Scroll>
            </Sider>
        );
    }
}

export default FameMenu;