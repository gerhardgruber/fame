import * as React from 'react';
import {observer} from 'mobx-react';
import Page from "../../components/Page";
import User from "../../stores/User";
import UiStore from "../../stores/UiStore";
import { Form, Input, Button } from "antd";
import { WrappedFormUtils } from "antd/lib/form/Form";
import FormItem from "antd/lib/form/FormItem";
import { Link } from 'react-router-dom';
import UserStore from '../../stores/UserStore';
import { UserForm } from '../../components/UserForm';

interface EditUserProps {
  userID?: number;
  form: WrappedFormUtils;
}

interface EditUserState {
  user: User;
}

const uiStore = UiStore.getInstance();
const userStore = UserStore.getInstance();

@observer
class _EditUser extends Page<EditUserProps, EditUserState> {
  state = {
    user: null
  }

  componentWillMount() {
    if ( this.props.userID) {
      userStore.loadUser(this.props.userID).then((u: User) => {
        this.setState({
          user: u
        });
      })
    }
  }

  pageTitle(): string {
    if (this.state.user) {
      return 'USERS_EDIT_USER';
    } else {
      return 'USERS_NEW_USER';
    }
  }

  renderContent(): JSX.Element {
    const { getFieldDecorator } = this.props.form;

    let passwordField = null;
    if (this.props.userID && this.state.user) {
      return <UserForm user={this.state.user} />;
    } else if (!this.props.userID) {
      return <UserForm />
    } else {
      return null;
    }
  }
}

const EditUser = Form.create()(_EditUser);
export {EditUser};