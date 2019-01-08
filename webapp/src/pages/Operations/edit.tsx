import * as React from 'react';
import {observer} from 'mobx-react';
import Page from "../../components/Page";
import Operation from "../../stores/Operation";
import UiStore from "../../stores/UiStore";
import { Form, Input, Button } from "antd";
import { WrappedFormUtils } from "antd/lib/form/Form";
import FormItem from "antd/lib/form/FormItem";
import { Link } from 'react-router-dom';
import OperationStore from '../../stores/OperationStore';
import { OperationForm } from '../../components/OperationForm';

interface EditOperationProps {
  operationID?: number;
  form: WrappedFormUtils;
}

interface EditOperationState {
  operation: Operation;
}

const uiStore = UiStore.getInstance();
const operationStore = OperationStore.getInstance();

@observer
class _EditOperation extends Page<EditOperationProps, EditOperationState> {
  state = {
    operation: null
  }

  componentWillMount() {
    if ( this.props.operationID) {
      operationStore.loadOperation(this.props.operationID).then((o: Operation) => {
        this.setState({
          operation: o
        });
      })
    }
  }

  pageTitle(): string {
    if (this.state.operation) {
      return 'OPERATIONS_EDIT_OPERATION';
    } else {
      return 'OPERATIONS_NEW_OPERATION';
    }
  }

  renderContent(): JSX.Element {
    const { getFieldDecorator } = this.props.form;

    let passwordField = null;
    if (this.props.operationID && this.state.operation) {
      return <OperationForm operation={this.state.operation} />;
    } else if (!this.props.operationID) {
      return <OperationForm />
    } else {
      return null;
    }
  }
}

const EditOperation = Form.create()(_EditOperation);
export {EditOperation};