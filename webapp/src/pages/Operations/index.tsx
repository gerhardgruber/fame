import * as React from 'react';
import { observer } from 'mobx-react';
import PageHeader from '../../components/PageHeader'
import { Layout, Row, Col, Table } from 'antd'
import Page from '../../components/Page';
import UiStore from '../../stores/UiStore';
import OperationStore from '../../stores/OperationStore';

const operationStore = OperationStore.getInstance( );
const uiStore = UiStore.getInstance( );

@observer
export default class Operations extends Page {
  columns = [ {
    title: uiStore.T( 'OPERATIONS_TITLE' ),
    dataIndex: 'Name'
  } ];

  componentDidMount( ) {
    operationStore.loadOperations();
  }

  pageTitle(): string {
    return "OPERATIONS";
  }

  renderContent(): JSX.Element {
    return <div>
             <Table
               columns={this.columns}
               dataSource={operationStore.operations}
               size={"small"}
               pagination={false} />
           </div>;
  }
}