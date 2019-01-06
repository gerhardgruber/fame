import * as React from 'react';
import { observer } from 'mobx-react';
import PageHeader from '../../components/PageHeader'
import { Layout, Row, Col, Table } from 'antd'
import Page from '../../components/Page';
import UserStore from '../../stores/UserStore';
import UiStore from '../../stores/UiStore';

const userStore = UserStore.getInstance( );
const uiStore = UiStore.getInstance( );

@observer
export default class Users extends Page {
  columns = [ {
    title: uiStore.T( 'USERS_USER_NAME' ),
    dataIndex: 'Name'
  }, {
    title: uiStore.T( 'USERS_USER_LAST_NAME' ),
    dataIndex: 'FirstName'
  }, {
    title: uiStore.T( 'USERS_USER_FIRST_NAME' ),
    dataIndex: 'LastName'
  } ];

  componentDidMount( ) {
    userStore.loadUsers();
  }

  pageTitle(): string {
    return "USERS";
  }

  renderContent(): JSX.Element {
    return <div>
             <Table
               columns={this.columns}
               dataSource={userStore.users}
               size={"small"}
               pagination={false} />
           </div>;
  }
}