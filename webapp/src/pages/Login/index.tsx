import * as React from 'react';
import { Form, Input, Button, Col, Row } from 'antd';
import UiStore from '../../stores/UiStore';
import { FormComponentProps } from 'antd/lib/form';
import User from '../../stores/User';
import { ButtonProps } from 'antd/lib/button';
import './theme.scss';

const FormItem = Form.Item;

const uiStore = UiStore.getInstance()

const tailFormItemLayout = {
    wrapperCol: {
        xs: {
            span: 16,
            offset: uiStore.formItemLayout.labelCol.sm.span,
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

    _forgotPassword: ButtonProps["onClick"] = (e) => {
        e.preventDefault();
    }

    renderFields() {
        const { getFieldDecorator } = this.props.form;
        return <div>
                                    <h1 style={{ marginTop: "0.5em", color: 'white' }}>
                                FAME - {uiStore.T("LOGIN_TITLE")}
                            </h1>
                            <FormItem {...uiStore.formItemLayout} hasFeedback>
                                {getFieldDecorator('username', {
                                    rules: [{ required: true, message: uiStore.T("LOGIN_USERNAME_NOT_GIVEN") }]
                                })(
                                    <Input placeholder={uiStore.T("LOGIN_USERNAME_PLACEHOLDER")} />
                                )}
                            </FormItem>
                            <FormItem {...uiStore.formItemLayout} hasFeedback >
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
                            {/*<FormItem {...tailFormItemLayout}>
                                <Button onClick={this._forgotPassword}>
                                    {uiStore.T("FORGOT_PASSWORD_BUTTON")}
                                </Button>
                                </FormItem>*/}
        </div>;
    }

    render() {


        return (
            <Form style={{backgroundColor: 'none'}} onSubmit={this.handleLogin.bind( this )} >
            <Row>
                <Col xs={0} sm={0} md={24} lg={24} xl={24} xxl={24}>
                    <div style={{
                        backgroundImage: 'url(/static/background.jpg)',
                        position: 'absolute',
                        top: '0',
                        left: '0',
                        width: '100%',
                        height: '100vh'
                    }}>
                        <div style={{
                            position: 'absolute',
                            backgroundColor: '#343d46',
                            width: '40%',
                            height: '100vh',
                            paddingLeft: '5rem',
                            paddingTop: '6rem'
                        }} className="slanted">
                            {this.renderFields()}


                        <div style={{
                            position: 'absolute',
                            top: '95%',
                            fontSize: '8pt',
                            color: 'whitesmoke'
                        }}>
                            Rettungshunde Niederösterreich
                        </div>
                        </div>
                    </div>
                </Col>
                <Col xs={24} sm={24} md={0} lg={0} xl={0} xxl={0}>
                    <div style={{
                            position: 'absolute',
                            backgroundColor: '#343d46',
                            width: '100%',
                            height: '100vh',
                            paddingLeft: '5rem',
                            paddingTop: '6rem'
                        }}>
                            {this.renderFields()}

                        </div>
                    </Col>
            </Row>
            </Form>
        );
    }
}

const Login = Form.create()(_Login);

export default Login;
