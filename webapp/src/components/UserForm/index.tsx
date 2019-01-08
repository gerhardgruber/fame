import * as React from 'react';
import {observer} from 'mobx-react';
import Page from "../../components/Page";
import User, { RightType } from "../../stores/User";
import UiStore from "../../stores/UiStore";
import { Form, Input, Button } from "antd";
import { WrappedFormUtils } from "antd/lib/form/Form";
import FormItem from "antd/lib/form/FormItem";
import { Link, Redirect } from 'react-router-dom';
import UserStore from '../../stores/UserStore';
import RightTypeSelect from '../RightTypeSelect';

interface UserFormProps {
  user?: User;
  form: WrappedFormUtils;
}

const uiStore = UiStore.getInstance();
const userStore = UserStore.getInstance();

@observer
class _UserForm extends React.Component<UserFormProps> {
  state = {
    gotoUsers: false
  };

  save = (e) => {
    e.preventDefault();

    this.props.form.validateFields((err, data) => {
      if (err) {
        return
      }

      if (this.props.user) {
        this.props.user.setData(data);
        userStore.saveUser(this.props.user).then( () => {
          userStore.loadUsers().then( () => {
            this.setState({
              gotoUsers: true
            });
          } );
        });
      } else {
        userStore.createUser(new User(data)).then( () => {
          userStore.loadUsers().then( () => {
            this.setState({
              gotoUsers: true
            });
          } );
        });
      }
    });
  }

  deleteUser = (e) => {
    e.preventDefault();
    userStore.deleteUser(this.props.user).then( () => {
      userStore.loadUsers().then( () => {
        this.setState({
          gotoUsers: true
        });
      } );
    });
  }

  render(): JSX.Element {
    const { getFieldDecorator } = this.props.form;

    let passwordField = null;
    let deleteButton = null;
    if (!this.props.user) {
      passwordField = <FormItem {...uiStore.formItemLayout} label={uiStore.T("USER_PW")} hasFeedback>
        {getFieldDecorator('PW', {
          rules: [{ required: true, message: uiStore.T("USER_PW_NOT_GIVEN") }]
        })(
          <Input placeholder={uiStore.T("USER_PW_PLACEHOLDER")} type="password"/>
        )}
      </FormItem>
    } else {
      deleteButton = <div style={{"display": "inline-block", "marginRight": "1rem"}}>
        <Button onClick={this.deleteUser} type="danger">
          {uiStore.T('DELETE')}
        </Button>
      </div>
    }

    let gotoUsers = null;
    if (this.state.gotoUsers) {
      gotoUsers = <Redirect to="/users" />;
    }

    return  <Form onSubmit={this.save}>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("USER_NAME")} hasFeedback>
                    {getFieldDecorator('Name', {
                        rules: [{ required: true, message: uiStore.T("USER_NAME_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("USER_NAME_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("USER_FIRST_NAME")} hasFeedback>
                    {getFieldDecorator('FirstName', {
                        rules: [{ required: true, message: uiStore.T("USER_FIRST_NAME_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("USER_FIRST_NAME_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("USER_LAST_NAME")} hasFeedback>
                    {getFieldDecorator('LastName', {
                        rules: [{ required: true, message: uiStore.T("USER_LAST_NAME_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("USER_LAST_NAME_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("USER_EMAIL")} hasFeedback>
                    {getFieldDecorator('EMail', {
                        rules: [{ required: true, message: uiStore.T("USER_EMAIL_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("USER_EMAIL_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("USER_RIGHT_TYPE")} hasFeedback>
                    {getFieldDecorator('RightType', {
                    })(
                        <RightTypeSelect />
                    )}
              </FormItem>
              {passwordField}

              {deleteButton}
              <div style={{"display": "inline-block", "marginRight": "1rem"}}>
                <Link to="/users"><Button>
                  {uiStore.T('CANCEL')}
                </Button></Link>
              </div>
              <div style={{"display": "inline-block"}}>
                <Button htmlType="submit" type="primary">
                  {uiStore.T('SAVE')}
                </Button>
              </div>
              {gotoUsers}
            </Form>
  }
}

const UserForm = Form.create({
  mapPropsToFields(props: UserFormProps) {
    const u = props.user;
    if(!u) return {
      RightType: Form.createFormField({value: 0})
    };

    return {
      Name: Form.createFormField({value: u.Name}),
      FirstName: Form.createFormField({value: u.FirstName}),
      LastName: Form.createFormField({value: u.LastName}),
      EMail: Form.createFormField({value: u.EMail}),
      RightType: Form.createFormField({value: u.RightType})
    }
  }
})(_UserForm);
export {UserForm};