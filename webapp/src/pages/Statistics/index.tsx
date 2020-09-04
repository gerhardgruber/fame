import * as React from 'react';
import { observer } from 'mobx-react';
import PageHeader from '../../components/PageHeader'
import { Layout, Row, Col, Table, Form, Input, Button, DatePicker } from 'antd'
import Page from '../../components/Page';
import UserStore from '../../stores/UserStore';
import UiStore from '../../stores/UiStore';
import { FormComponentProps } from 'antd/lib/form';
import FormItem from 'antd/lib/form/FormItem';
import Api from '../../core/Api';
import { API_ROOT } from '../../../config/fame';

const userStore = UserStore.getInstance( );
const uiStore = UiStore.getInstance( );

@observer
export class _Statistics extends Page<FormComponentProps> {
  pageTitle(): string {
    return "STATISTICS";
  }

  createStatistic(e) {
    e.preventDefault();

    this.props.form.validateFields((err, values) => {
        if (err) {
            return;
        }

        window.open( Api.buildURL( "/statistics/attendance", {
          fromDate: values[ "fromDate" ],
          toDate: values[ "toDate" ]
        } ) );
    });
  }

  renderContent(): JSX.Element {
    const { getFieldDecorator } = this.props.form;

    return <div>
             <Form onSubmit={this.createStatistic.bind( this )} layout={'vertical'} hideRequiredMark={true}>
                <FormItem {...uiStore.formItemLayout} label={uiStore.T("STATISTICS_FROM_DATE")}>
                  {getFieldDecorator('fromDate', {
                    rules: [{ required: true, message: uiStore.T("STATISTICS_FROM_DATE_NOT_GIVEN") }]
                  })(
                    <DatePicker placeholder={uiStore.T("STATISTICS_FROM_DATE_PLACEHOLDER")} />
                  )}
                </FormItem>
                <FormItem {...uiStore.formItemLayout} label={uiStore.T("STATISTICS_TO_DATE")}>
                  {getFieldDecorator('toDate', {
                    rules: [{ required: true, message: uiStore.T("STATISTICS_TO_DATE_NOT_GIVEN") }]
                  })(
                    <DatePicker placeholder={uiStore.T("STATISTICS_TO_DATE_PLACEHOLDER")} />
                  )}
                </FormItem>

                <div style={{"display": "inline-block"}}>
                  <Button htmlType="submit" type="primary">
                      {uiStore.T('OK')}
                  </Button>
                </div>
             </Form>
           </div>;
  }
}

const Statistics = Form.create()(_Statistics);

export default Statistics;
