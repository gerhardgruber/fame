import * as React from 'react';
import { observer } from 'mobx-react';
import PageHeader from '../../components/PageHeader'
import { Layout, Row, Col, Table, Button } from 'antd'
import Page from '../../components/Page';
import UserStore from '../../stores/UserStore';
import UiStore from '../../stores/UiStore';
import { Link, Redirect } from 'react-router-dom';
import User from '../../stores/User';

const userStore = UserStore.getInstance( );
const uiStore = UiStore.getInstance( );

@observer
export default class Users extends Page {
  state = {
    navigateTo: null
  };

  columns = [ {
    title: uiStore.T( 'USERS_USER_NAME' ),
    dataIndex: 'Name',
    sortDirections: []
  }, {
    title: uiStore.T( 'USERS_USER_LAST_NAME' ),
    dataIndex: 'FirstName',
    sortDirections: []
  }, {
    title: uiStore.T( 'USERS_USER_FIRST_NAME' ),
    dataIndex: 'LastName',
    sortDirections: []
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