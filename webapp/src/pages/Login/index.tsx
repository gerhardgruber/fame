import * as React from 'react';
import { Form, Input, Button } from 'antd';
import UiStore from '../../stores/UiStore';
import { FormComponentProps } from 'antd/lib/form';
import User from '../../stores/User';

const FormItem = Form.Item;

const uiStore = UiStore.getInstance()

const tailFormItemLayout = {
    wrapperCol: {
        xs: {
            span: 24,
            offset: 0,
        },
        sm: {
            span: 16,
            offset: uiStore.formItemLayout.labelCol.sm.span,
        },
    },
};

class _Login extends React.Component<FormComponentProps> {
    handleLogin(e) {
        e.preventDefault();

        this.props.form.validateFields((err, values) => {
            if (err) {
                return;
            }

            uiStore.login( values["username"], values["password"] );
        });
    }

    render() {
        const { getFieldDecorator } = this.props.form;

        return (
            <Form onSubmit={this.handleLogin.bind( this )} >
                <h1 style={{ marginLeft: "2em", marginTop: "0.5em" }}>
                    FAME - {uiStore.T("LOGIN_TITLE")}
                </h1>
                <FormItem {...uiStore.formItemLayout} label={uiStore.T("LOGIN_USERNAME")} hasFeedback>
                    {getFieldDecorator('username', {
                        rules: [{ required: true, message: uiStore.T("LOGIN_USERNAME_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("LOGIN_USERNAME_PLACEHOLDER")} />
                    )}
                </FormItem>
                <FormItem {...uiStore.formItemLayout} label={uiStore.T("LOGIN_PASSWORD")} hasFeedback >
                    {getFieldDecorator('password', {
                        rules: [{ required: true, message: uiStore.T("LOGIN_PASSWORD_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("LOGIN_PASSWORD_PLACEHOLDER")} type="password" />
                    )}
                </FormItem>

                <FormItem {...tailFormItemLayout}>
                    <Button type="primary" htmlType="submit" className="login-form-button">
                        {uiStore.T("LOGIN_LOGIN_BUTTON")}
                    </Button>
                </FormItem>
            </Form>
        );
    }
}

const Login = Form.create()(_Login);

export default Login;
