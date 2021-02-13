import * as React from 'react';
import { observer } from 'mobx-react';
import PageHeader from '../../components/PageHeader'
import { Layout, Row, Col, Table, Button, Icon } from 'antd'
import Page from '../../components/Page';
import UserStore from '../../stores/UserStore';
import UiStore from '../../stores/UiStore';
import { Link, Redirect } from 'react-router-dom';
import User from '../../stores/User';
import Column from 'antd/lib/table/Column';
import { ColumnProps } from 'antd/lib/table';
import Dashboard from '../Dashboard';

const userStore = UserStore.getInstance( );
const uiStore = UiStore.getInstance( );

@observer
export default class Users extends Page {
  state = {
    navigateTo: null
  };

  renderResult(value: number, total: number) {
    let factor = 0;
    if (total !== 0 ) {
      factor = value / total;
    }

    let icon: string;
    let color: string;
    if ( factor >= 0.75 || total === 0 ) {
      icon = "check-circle";
      color = Dashboard.green;
    } else if ( factor >= 0.5 ) {
      icon = "warning";
      color = Dashboard.yellow;
    } else if ( factor >= 0.25 ) {
      icon = "warning";
      color = Dashboard.orange;
    } else {
      icon = "close-circle";
      color = Dashboard.red;
    }

    return <Icon type={icon} twoToneColor={color} theme="twoTone" />
  }

  columns: ColumnProps<any>[] = [ {
    title: uiStore.T( 'USERS_USER_NAME' ),
    dataIndex: 'Name',
    sortDirections: []
  }, {
    title: uiStore.T( 'USERS_USER_FIRST_NAME' ),
    dataIndex: 'FirstName',
    sortDirections: []
  }, {
    title: uiStore.T( 'USERS_USER_LAST_NAME' ),
    dataIndex: 'LastName',
    sortDirections: []
  }, {
    title: uiStore.T('DASHBOARD_FEEDBACK'),
    width: '20%',
    render: (txt, record) => {
      if (userStore.stati && userStore.stati[record.ID]) {
        return this.renderResult(userStore.stati[record.ID][1], userStore.stati[record.ID][0])
      }
      return null;
    }
  }, {
    title: uiStore.T('DASHBOARD_PRESENT'),
    width: '20%',
    render: (txt, record) => {
      if (userStore.stati && userStore.stati[record.ID]) {
        return this.renderResult(userStore.stati[record.ID][2], userStore.stati[record.ID][0])
      }
      return null;
    }
  } ];

  componentDidMount( ) {
    userStore.loadUsers();
  }

  pageTitle(): string {
    return "USERS";
  }

  renderButtons(): JSX.Element {
    return <Link to="/users/new"><Button>
      {uiStore.T( 'USERS_ADD_USER' )}
    </Button></Link>
  }

  rowClicked = (record: User) => {
    this.setState({
      navigateTo: record.ID
    });
  }

  renderContent(): JSX.Element {
    return <div>
             <Table
               columns={this.columns}
               dataSource={userStore.users}
               size={"small"}
               pagination={false}
               onRowClick={this.rowClicked} />
             {this.state.navigateTo ? <Redirect push to={"/users/" + this.state.navigateTo} /> : null}
           </div>;
  }
}