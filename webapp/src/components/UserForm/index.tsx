import * as React from 'react';
import {observer} from 'mobx-react';
import Page from "../../components/Page";
import User, { RightType } from "../../stores/User";
import UiStore from "../../stores/UiStore";
import { Form, Input, Button, Modal, DatePicker, Row } from "antd";
import { WrappedFormUtils } from "antd/lib/form/Form";
import FormItem from "antd/lib/form/FormItem";
import { Link, Redirect } from 'react-router-dom';
import UserStore from '../../stores/UserStore';
import RightTypeSelect from '../RightTypeSelect';
import moment from 'moment';
import { PauseAction, PauseType } from '../../stores/PauseAction';

interface UserFormProps {
  user?: User;
  form: WrappedFormUtils;
}

const uiStore = UiStore.getInstance();
const userStore = UserStore.getInstance();

@observer
class _UserForm extends React.Component<UserFormProps> {
  state = {
    gotoUsers: false,
    pauseDateModal: null,
    pauseDate: null
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

  getPauseButton(pauseData: PauseAction, type: PauseType) {
    if (!pauseData || pauseData.EndTime) {
      return <Button onClick={() => {
        this.setState({
          pauseDateModal: {
            caption: "USER_START_PAUSE_DATE",
            start: true,
            type: type
          },
          pauseDate: new Date()
        })
      }}>{uiStore.T('USER_START_PAUSE')}</Button>
    } else {
      return <Button onClick={() => {
        this.setState({
          pauseDateModal: {
            caption: "USER_STOP_PAUSE_DATE",
            start: false,
            type: type
          },
          pauseDate: new Date()
        })
      }}>{uiStore.T('USER_STOP_PAUSE')}</Button>
    }
  }

  render(): JSX.Element {
    const { getFieldDecorator } = this.props.form;

    let passwordField = null;
    let deleteButton = null;
    let pauseFields = null;
    if (!this.props.user) {
      passwordField = <FormItem {...uiStore.formItemLayout} label={uiStore.T("USER_PW")} hasFeedback>
        {getFieldDecorator('PW', {
          rules: [{ required: true, message: uiStore.T("USER_PW_NOT_GIVEN") }]
        })(
          <Input placeholder={uiStore.T("USER_PW_PLACEHOLDER")} type="password"/>
        )}
      </FormItem>
    } else {
      pauseFields = [
        <FormItem {...uiStore.formItemLayout} label={uiStore.T("USER_TRAINING_PAUSE")} hasFeedback>
          {this.getPauseButton(this.props.user.TrainingPause, PauseType.TrainingPause)}
        </FormItem>,
        <FormItem {...uiStore.formItemLayout} label={uiStore.T("USER_OPERATION_PAUSE")} hasFeedback>
          {this.getPauseButton(this.props.user.OperationPause, PauseType.OperationPause)}
        </FormItem>
      ]
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

    let pauseDateModal = null;
    if (this.state.pauseDateModal) {
      pauseDateModal = <Modal visible={true} onCancel={
        () => {
          this.setState({
            pauseDateModal: null
          })
        }
      } onOk={() => {
        if ( this.state.pauseDateModal.start ) {
          this.props.user.startPause(
            this.state.pauseDateModal.type,
            this.state.pauseDate
          )
        } else {
          this.props.user.stopPause(
            this.state.pauseDateModal.type,
            this.state.pauseDate
          )
        }
        this.setState({
          pauseDateModal: null
        })
      }}>
        <Row style={{marginBottom: '0.5rem'}}>
          {uiStore.T(this.state.pauseDateModal.caption)}
        </Row>
        <Row>
          <DatePicker
            defaultValue={moment(this.state.pauseDate)}
            onChange={(date) => {
              this.setState({
                pauseDate: date.toDate()
              })
            }}
            />
        </Row>
      </Modal>
    }

    return  <Form onSubmit={this.save}>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("USER_NAME")} hasFeedback>
                    {getFieldDecorator('Name', {
                        rules: [{ required: true, message: uiStore.T("USER_NAME_NOT_GIVEN") }]
                    })(
                        <Input disabled={this.props.user ? true : false} placeholder={uiStore.T("USER_NAME_PLACEHOLDER")} />
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
              {pauseFields}

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
              {pauseDateModal}
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