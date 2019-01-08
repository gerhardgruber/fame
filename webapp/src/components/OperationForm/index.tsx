import * as React from 'react';
import {observer} from 'mobx-react';
import Page from "../../components/Page";
import Operation from "../../stores/Operation";
import UiStore from "../../stores/UiStore";
import { Form, Input, Button } from "antd";
import { WrappedFormUtils } from "antd/lib/form/Form";
import FormItem from "antd/lib/form/FormItem";
import { Link, Redirect } from 'react-router-dom';
import OperationStore from '../../stores/OperationStore';

interface OperationFormProps {
  operation?: Operation;
  form: WrappedFormUtils;
}

const uiStore = UiStore.getInstance();
const operationStore = OperationStore.getInstance();

@observer
class _OperationForm extends React.Component<OperationFormProps> {
  state = {
    gotoOperations: false
  };

  save = (e) => {
    e.preventDefault();

    this.props.form.validateFields((err, data) => {
      if (err) {
        return
      }

      if (this.props.operation) {
        this.props.operation.setData(data);
        operationStore.saveOperation(this.props.operation).then( () => {
          operationStore.loadOperations().then( () => {
            this.setState({
              gotoOperations: true
            });
          } );
        });
      } else {
        operationStore.createOperation(new Operation(data)).then( () => {
          operationStore.loadOperations().then( () => {
            this.setState({
              gotoOperations: true
            });
          } );
        });
      }
    });
  }

  deleteOperation = (e) => {
    e.preventDefault();
    operationStore.deleteOperation(this.props.operation).then( () => {
      operationStore.loadOperations().then( () => {
        this.setState({
          gotoOperations: true
        });
      } );
    });
  }

  render(): JSX.Element {
    const { getFieldDecorator } = this.props.form;

    let passwordField = null;
    let deleteButton = null;
    if (this.props.operation) {
      deleteButton = <div style={{"display": "inline-block", "marginRight": "1rem"}}>
        <Button onClick={this.deleteOperation} type="danger">
          {uiStore.T('DELETE')}
        </Button>
      </div>
    }

    let gotoOperations = null;
    if (this.state.gotoOperations) {
      gotoOperations = <Redirect to="/operations" />;
    }

    return  <Form onSubmit={this.save}>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("OPERATION_TITLE")} hasFeedback>
                    {getFieldDecorator('Title', {
                        rules: [{ required: true, message: uiStore.T("OPERATION_TITLE_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("OPERATION_TITLE_PLACEHOLDER")} />
                    )}
              </FormItem>
              <h2>{uiStore.T("OPERATION_MISSING_PERSON")}</h2>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("OPERATION_FIRST_NAME")} hasFeedback>
                    {getFieldDecorator('FirstName', {
                        rules: [{ required: true, message: uiStore.T("OPERATION_FIRST_NAME_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("OPERATION_FIRST_NAME_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("OPERATION_LAST_NAME")} hasFeedback>
                    {getFieldDecorator('LastName', {
                        rules: [{ required: true, message: uiStore.T("OPERATION_LAST_NAME_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("OPERATION_LAST_NAME_PLACEHOLDER")} />
                    )}
              </FormItem>

              {deleteButton}
              <div style={{"display": "inline-block", "marginRight": "1rem"}}>
                <Link to="/operations"><Button>
                  {uiStore.T('CANCEL')}
                </Button></Link>
              </div>
              <div style={{"display": "inline-block"}}>
                <Button htmlType="submit" type="primary">
                  {uiStore.T('SAVE')}
                </Button>
              </div>
              {gotoOperations}
            </Form>
  }
}

const OperationForm = Form.create({
  mapPropsToFields(props: OperationFormProps) {
    const o = props.operation;
    if(!o) return {};

    return {
      Title: Form.createFormField({value: o.Title}),
      FirstName: Form.createFormField({value: o.FirstName}),
      LastName: Form.createFormField({value: o.LastName})
    }
  }
})(_OperationForm);
export {OperationForm};