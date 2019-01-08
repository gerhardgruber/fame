import * as React from 'react';
import { observer } from 'mobx-react';
import PageHeader from '../../components/PageHeader'
import { Layout, Row, Col, Table, Button } from 'antd'
import Page from '../../components/Page';
import UiStore from '../../stores/UiStore';
import OperationStore from '../../stores/OperationStore';
import {Link, Redirect} from 'react-router-dom';
import Operation from '../../stores/Operation';

const operationStore = OperationStore.getInstance( );
const uiStore = UiStore.getInstance( );

@observer
export default class Operations extends Page {
  state = {
    navigateTo: null
  };

  columns = [ {
    title: uiStore.T( 'OPERATIONS_TITLE' ),
    dataIndex: 'Title'
  } ];

  componentDidMount( ) {
    operationStore.loadOperations();
  }

  pageTitle(): string {
    return "OPERATIONS";
  }

  rowClicked = (record: Operation) => {
    this.setState({
      navigateTo: record.ID
    });
  }

  renderButtons(): JSX.Element {
    return <Link to="/operations/new"><Button>
      {uiStore.T( 'OPERATIONS_ADD_OPERATION' )}
    </Button></Link>
  }

  renderContent(): JSX.Element {
    return <div>
             <Table
               columns={this.columns}
               dataSource={operationStore.operations}
               size={"small"}
               pagination={false}
               onRowClick={this.rowClicked} />
              {this.state.navigateTo ? <Redirect push to={"/operations/" + this.state.navigateTo} /> : null}
           </div>;
  }
}