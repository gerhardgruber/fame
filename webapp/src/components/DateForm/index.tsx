import * as React from 'react';
import {observer} from 'mobx-react';
import UiStore from "../../stores/UiStore";
import { Form, Input, Button, DatePicker, Select, Spin, List, Icon } from "antd";
import { WrappedFormUtils } from "antd/lib/form/Form";
import FormItem from "antd/lib/form/FormItem";
import { Link, Redirect } from 'react-router-dom';
import DateStore from '../../stores/DateStore';
import { DateModel, DateFeedback } from '../../stores/DateStore/Date';
import TextArea from 'antd/lib/input/TextArea';
import moment from 'moment';
import ButtonGroup from 'antd/lib/button/button-group';
import { ButtonProps } from 'antd/lib/button';
import DateCategoryStore from '../../stores/DateCategoryStore';

interface DateFormProps {
  date?: DateModel;
  form: WrappedFormUtils;
}

const uiStore = UiStore.getInstance();
const dateStore = DateStore.getInstance();
const dateCategoryStore = DateCategoryStore.getInstance();

@observer
class _DateForm extends React.Component<DateFormProps> {
  state = {
    gotoDates: false
  };

  save = (e) => {
    e.preventDefault();

    this.props.form.validateFields((err, data) => {
      if (err) {
        return
      }

      if (this.props.date) {
        this.props.date.setData(data);
        dateStore.saveDate(this.props.date).then( () => {
          dateStore.loadDates().then( () => {
            this.setState({
              gotoDates: true
            });
          } );
        });
      } else {
        dateStore.createDate(new DateModel(data)).then( () => {
          dateStore.loadDates().then( () => {
            this.setState({
              gotoDates: true
            });
          } );
        });
      }
    });
  }

  deleteDate = (e) => {
    e.preventDefault();
    dateStore.deleteDate(this.props.date).then( () => {
      dateStore.loadDates().then( () => {
        this.setState({
          gotoDates: true
        });
      } );
    });
  }

  feedbackYesClicked = (item: DateFeedback) => {
    item.Feedback = uiStore.dateFeedbackTypes.Yes;
    dateStore.sendFeedback(this.props.date.ID, item).then( (response) => {
      item.UpdatedAt = response.UpdatedAt;
    });
    this.forceUpdate();
  };

  feedbackNoClicked = (item: DateFeedback) => {
    item.Feedback = uiStore.dateFeedbackTypes.No;
    dateStore.sendFeedback(this.props.date.ID, item).then( (response) => {
      item.UpdatedAt = response.UpdatedAt;
    });
    this.forceUpdate();
  };

  renderAnswerButton = (item: DateFeedback): JSX.Element => {
    if (item.User.ID !== uiStore.currentUser.ID) {
      return null;
    }

    return <ButtonGroup style={{marginLeft: '1rem'}}>
      <Button onClick={() => this.feedbackYesClicked(item)} style={{backgroundColor: item.Feedback === uiStore.dateFeedbackTypes.Yes ? '#76FF03' : '#CCFF90', color: 'black'}}>
        {uiStore.T( 'DATE_YES' )}
      </Button>
      <Button onClick={() => this.feedbackNoClicked(item)} style={{backgroundColor: item.Feedback === uiStore.dateFeedbackTypes.No ? '#FF1744' : '#FF8A80', color: 'black'}}>
        {uiStore.T( 'DATE_NO' )}
      </Button>
    </ButtonGroup>
  };

  renderFeedbacks = (): JSX.Element  => {
    if (this.props.date) {
      return <div style={{marginBottom: '2rem'}}>
        <h1>{uiStore.T('DATE_FEEDBACKS')}</h1>
        <List
          dataSource={this.props.date.Feedbacks}
          bordered={false}
          renderItem={(item: DateFeedback) => {
            return <List.Item>
              <span style={{marginRight: '1rem'}}>
                {item.Feedback === uiStore.dateFeedbackTypes["Yes"] ? <Icon style={{color: 'green'}} type="check-circle" /> : <Icon style={{color: 'red'}} type="close-circle" />}
              </span>
              {item.User.FirstName} {item.User.LastName}
              {this.renderAnswerButton(item)}
            </List.Item>
          }} />
      </div>
    }

    return null;
  }

  render(): JSX.Element {
    if (!uiStore.dateTypes) {
      return <Spin />;
    }

    const { getFieldDecorator } = this.props.form;

    let deleteButton = null;
    if (this.props.date) {
      deleteButton = <div style={{"display": "inline-block", "marginRight": "1rem"}}>
        <Button onClick={this.deleteDate} type="danger">
          {uiStore.T('DELETE')}
        </Button>
      </div>
    }

    let gotoDates = null;
    if (this.state.gotoDates) {
      gotoDates = <Redirect to="/dates" />;
    }

    const dateCategories = ( dateCategoryStore.dateCategories || [] ).map((dc) => {
      return <Select.Option key={dc.ID} value={dc.ID}>{dc.Name}</Select.Option>
    })

    return  <Form onSubmit={this.save}>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("DATE_TITLE")} hasFeedback>
                    {getFieldDecorator('Title', {
                        rules: [{ required: true, message: uiStore.T("DATE_TITLE_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("DATE_TITLE_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("DATE_DESCRIPTION")} hasFeedback>
                    {getFieldDecorator('Description', {})(
                        <TextArea placeholder={uiStore.T("DATE_DESCRIPTION_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("DATE_LOCATION")} hasFeedback>
                    {getFieldDecorator('Location', {
                        rules: [{ required: true, message: uiStore.T("DATE_LOCATION_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("DATE_LOCATION_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("DATE_START_TIME")} hasFeedback>
                    {getFieldDecorator('StartTime', {
                        rules: [{ required: true, message: uiStore.T("DATE_START_TIME_NOT_GIVEN") }]
                    })(
                        <DatePicker showTime={true} placeholder={uiStore.T("DATE_START_TIME_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("DATE_END_TIME")} hasFeedback>
                    {getFieldDecorator('EndTime', {
                        rules: [{ required: true, message: uiStore.T("DATE_END_TIME_NOT_GIVEN") }]
                    })(
                        <DatePicker showTime={true} placeholder={uiStore.T("DATE_END_TIME_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("DATE_CATEGORY")} hasFeedback>
                    {getFieldDecorator('CategoryID', {})(
                        <Select>
                          {dateCategories}
                        </Select>
                    )}
              </FormItem>

              {this.renderFeedbacks()}

              {deleteButton}
              <div style={{"display": "inline-block", "marginRight": "1rem"}}>
                <Link to="/dates"><Button>
                  {uiStore.T('CANCEL')}
                </Button></Link>
              </div>
              <div style={{"display": "inline-block"}}>
                <Button htmlType="submit" type="primary">
                  {uiStore.T('SAVE')}
                </Button>
              </div>
              {gotoDates}
            </Form>
  }
}

const DateForm = Form.create({
  mapPropsToFields(props: DateFormProps) {
    dateCategoryStore.loadDateCategories();

    let fallbackID = 0;
    if (dateCategoryStore.dateCategories) {
      fallbackID = dateCategoryStore.dateCategories[ 0 ].ID;
    }

    const dt = props.date;
    if(!dt) return { CategoryID: fallbackID };

    return {
      Title: Form.createFormField({value: dt.Title}),
      Description: Form.createFormField({value: dt.Description}),
      Location: Form.createFormField({value: dt.Location}),
      StartTime: Form.createFormField({value: moment(dt.StartTime || new Date())}),
      EndTime: Form.createFormField({value: moment(dt.EndTime || new Date())}),
      CategoryID: Form.createFormField({value: dt.CategoryID || fallbackID})
    }
  }
})(_DateForm);
export {DateForm};