import * as React from 'react';
import { observer } from 'mobx-react';
import PageHeader from '../../components/PageHeader'
import { Layout, Row, Col, Table, Form, Input, Button } from 'antd'
import Page from '../../components/Page';
import UserStore from '../../stores/UserStore';
import UiStore from '../../stores/UiStore';
import { FormComponentProps } from 'antd/lib/form';
import FormItem from 'antd/lib/form/FormItem';

const userStore = UserStore.getInstance( );
const uiStore = UiStore.getInstance( );

@observer
class _ChangePassword extends Page<FormComponentProps> {
  state = {
    confirmDirty: false,
  }

  pageTitle(): string {
    return "CHANGE_PASSWORD";
  }

  changePassword(e) {
    e.preventDefault();

    this.props.form.validateFields((err, values) => {
        if (err) {
            return;
        }

        uiStore.changePassword(values["currentPassword"], values["newPassword"]).then( () => {
          this.setState({
            gotoOrders: true
          })
        });
    });
  }

  compareToFirstPassword = (rule, value, callback) => {
    const form = this.props.form;
    if (value && value !== form.getFieldValue('newPassword')) {
      callback(uiStore.T('CHANGE_PASSWORD_NOT_SAME'));
    } else {
      callback();
    }
  }

  validateToNextPassword = (rule, value, callback) => {
    const form = this.props.form;
    if (value && this.state.confirmDirty) {
      form.validateFields(['confirm'], { force: true }, null);
    }
    callback();
  }

  handleConfirmBlur = (e) => {
    const value = e.target.value;
    this.setState({ confirmDirty: this.state.confirmDirty || !!value });
  }

  renderContent(): JSX.Element {
    const { getFieldDecorator } = this.props.form;

    return <div>
             <Form onSubmit={this.changePassword.bind( this )} layout={'vertical'} hideRequiredMark={true}>
                <FormItem {...uiStore.formItemLayout} label={uiStore.T("CHANGE_PASSWORD_CURRENT_PASSWORD")}>
                  {getFieldDecorator('currentPassword', {
                    rules: [{ required: true, message: uiStore.T("CHANGE_PASSWORD_CURRENT_PASSWORD_NOT_GIVEN") }]
                  })(
                    <Input placeholder={uiStore.T("CHANGE_PASSWORD_CURRENT_PASSWORD_PLACEHOLDER")} type="password" />
                  )}
                </FormItem>
                <FormItem {...uiStore.formItemLayout} label={uiStore.T("CHANGE_PASSWORD_NEW_PASSWORD")}>
                  {getFieldDecorator('newPassword', {
                    rules: [{
                      required: true, message: uiStore.T("CHANGE_PASSWORD_NEW_PASSWORD_NOT_GIVEN")
                    },
                    {
                      validator: this.validateToNextPassword
                    }],
                  })(
                    <Input placeholder={uiStore.T("CHANGE_PASSWORD_NEW_PASSWORD_PLACEHOLDER")} type="password" />
                  )}
                </FormItem>
                <FormItem {...uiStore.formItemLayout} label={uiStore.T("CHANGE_PASSWORD_REPEAT_PASSWORD")}>
                  {getFieldDecorator('repeatPassword', {
                    rules: [{
                      required: true, message: uiStore.T("CHANGE_PASSWORD_REPEAT_PASSWORD_NOT_GIVEN")
                    }, {
                      validator: this.compareToFirstPassword
                    }],
                  })(
                    <Input placeholder={uiStore.T("CHANGE_PASSWORD_REPEAT_PASSWORD_PLACEHOLDER")} type="password" onBlur={this.handleConfirmBlur} />
                  )}
                </FormItem>

                <div style={{"display": "inline-block", "marginRight": "1rem"}}>
                  <Button onClick={() => {
                    this.setState({
                      gotoOrders: true
                    });
                  }}>
                      {uiStore.T('CANCEL')}
                  </Button>
                </div>
                <div style={{"display": "inline-block"}}>
                  <Button htmlType="submit" type="primary">
                      {uiStore.T('SAVE')}
                  </Button>
                </div>
             </Form>
           </div>;
  }
}

const ChangePassword = Form.create()(_ChangePassword);

export default ChangePassword;
