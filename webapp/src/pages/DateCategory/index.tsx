import * as React from 'react';
import { observer } from 'mobx-react';
import PageHeader from '../../components/PageHeader'
import { Layout, Row, Col, Table, Button } from 'antd'
import Page from '../../components/Page';
import UserStore from '../../stores/UserStore';
import UiStore from '../../stores/UiStore';
import { Link, Redirect } from 'react-router-dom';
import User from '../../stores/User';
import DateCategoryStore from '../../stores/DateCategoryStore';

const dateCategoryStore = DateCategoryStore.getInstance( );
const uiStore = UiStore.getInstance( );

@observer
export default class DateCategories extends Page {
  state = {
    navigateTo: null
  };

  columns = [ {
    title: uiStore.T( 'DATE_CATEGORIES_NAME' ),
    dataIndex: 'Name',
    sortDirections: []
  } ];

  componentDidMount( ) {
    dateCategoryStore.loadDateCategories();
  }

  pageTitle(): string {
    return "DATE_CATEGORIES";
  }

  renderButtons(): JSX.Element {
    return <Link to="/date_categories/new"><Button>
      {uiStore.T( 'DATE_CATEGORIES_ADD' )}
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
               dataSource={dateCategoryStore.dateCategories}
               size={"small"}
               pagination={false}
               onRowClick={this.rowClicked} />
             {this.state.navigateTo ? <Redirect push to={"/date_categories/" + this.state.navigateTo} /> : null}
           </div>;
  }
}