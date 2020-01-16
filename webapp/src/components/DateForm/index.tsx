import * as React from 'react';
import {observer} from 'mobx-react';
import UiStore from "../../stores/UiStore";
import { Form, Input, Button, DatePicker, Select, Spin, List, Icon, TimePicker, Row, Col } from "antd";
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
import { computed } from 'mobx';
import { RightType } from '../../stores/User';

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

      data.StartTime.year(data.StartDate.year())
      data.StartTime.month(data.StartDate.month())
      data.StartTime.date(data.StartDate.date())

      data.EndTime.year(data.EndDate.year())
      data.EndTime.month(data.EndDate.month())
      data.EndTime.date(data.EndDate.date())

      if (this.props.date) {
        this.props.date.setData(data);
        dateStore.saveDate(this.props.date).then( () => {
          dateStore.loadDates(false).then( () => {
            this.setState({
              gotoDates: true
            });
          } );
        });
      } else {
        dateStore.createDate(new DateModel(data)).then( () => {
          dateStore.loadDates(false).then( () => {
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
      dateStore.loadDates(false).then( () => {
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

  renderFeedbackIcon(feedback: number) {
    if (feedback === uiStore.dateFeedbackTypes["Yes"]) {
      return <Icon style={{color: 'green'}} type="check-circle" />;
    } else if (feedback === uiStore.dateFeedbackTypes["No"]) {
      return <Icon style={{color: 'red'}} type="close-circle" />;
    } else {
      return <Icon style={{color: 'orange'}} type="warning" />;
    }
  }

  renderFeedbacks = (): JSX.Element  => {
    if (this.props.date) {
      return <div style={{marginBottom: '2rem'}}>
        <h1>{uiStore.T('DATE_FEEDBACKS')}</h1>
        <List
          dataSource={this.props.date.orderedFeedbacksWithHeaders}
          bordered={false}
          renderItem={(item) => {
            if (item instanceof DateFeedback) {
              return <List.Item>
                <span style={{marginRight: '1rem'}}>
                  {this.renderFeedbackIcon(item.Feedback)}
                </span>
                {item.User.FirstName} {item.User.LastName}
                {this.renderAnswerButton(item)}
              </List.Item>
            } else {
              return <List.Item>
                <b>{uiStore.T( `DATE_FEEDBACK_${item.Feedback}` ) + ` (${item.count})` }</b>
              </List.Item>
            }
          }} />
      </div>
    }

    return null;
  }

  @computed get editable(): boolean {
    return uiStore.isAdmin();
  }

  renderButtons(): JSX.Element[] {
    if (!this.editable) {
      return null;
    }

    let deleteButton = null;
    if (this.props.date) {
      deleteButton = <div style={{"display": "inline-block", "marginRight": "1rem"}}>
        <Button onClick={this.deleteDate} type="danger">
          {uiStore.T('DELETE')}
        </Button>
      </div>
    }

    return [
      deleteButton,
      <div style={{"display": "inline-block", "marginRight": "1rem"}}>
        <Link to="/dates"><Button>
          {uiStore.T('CANCEL')}
        </Button></Link>
      </div>,
      <div style={{"display": "inline-block"}}>
        <Button htmlType="submit" type="primary">
          {uiStore.T('SAVE')}
        </Button>
      </div>
    ]
  }

  render(): JSX.Element {
    if (!uiStore.dateTypes) {
      return <Spin />;
    }

    const { getFieldDecorator } = this.props.form;

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
                        <Input disabled={!this.editable} placeholder={uiStore.T("DATE_TITLE_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("DATE_DESCRIPTION")} hasFeedback>
                    {getFieldDecorator('Description', {})(
                        <TextArea disabled={!this.editable} placeholder={uiStore.T("DATE_DESCRIPTION_PLACEHOLDER")} />
                    )}
              </FormItem>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("DATE_LOCATION")} hasFeedback>
                    {getFieldDecorator('Location', {
                        rules: [{ required: true, message: uiStore.T("DATE_LOCATION_NOT_GIVEN") }]
                    })(
                        <Input disabled={!this.editable} placeholder={uiStore.T("DATE_LOCATION_PLACEHOLDER")} />
                    )}
              </FormItem>
              <Row>
                <Col xs={10}>
                  <FormItem labelCol={{xs:12}} wrapperCol={{xs:12}} label={uiStore.T("DATE_START_TIME")} hasFeedback>
                        {getFieldDecorator('StartDate', {
                            rules: [{ required: true, message: uiStore.T("DATE_START_DATE_NOT_GIVEN") }]
                        })(
                            <DatePicker
                              disabled={!this.editable}
                              placeholder={uiStore.T("DATE_START_DATE_PLACEHOLDER")}
                              />
                        )}
                  </FormItem>
                </Col>
                <Col xs={14}>
                  <FormItem wrapperCol={{xs:24}} hasFeedback>
                        {getFieldDecorator('StartTime', {
                            rules: [{ required: true, message: uiStore.T("DATE_START_TIME_NOT_GIVEN") }]
                        })(
                            <TimePicker
                              disabled={!this.editable}
                              placeholder={uiStore.T("DATE_START_TIME_PLACEHOLDER")}
                              format={"HH:mm"}
                              />
                        )}
                  </FormItem>
                </Col>
              </Row>
              <Row>
                <Col xs={10}>
                  <FormItem labelCol={{xs:12}} wrapperCol={{xs:12}} label={uiStore.T("DATE_END_TIME")} hasFeedback>
                        {getFieldDecorator('EndDate', {
                            rules: [{ required: true, message: uiStore.T("DATE_END_DATE_NOT_GIVEN") }]
                        })(
                            <DatePicker
                              disabled={!this.editable}
                              placeholder={uiStore.T("DATE_END_DATE_PLACEHOLDER")}
                              />
                        )}
                  </FormItem>
                </Col>
                <Col xs={14}>
                  <FormItem wrapperCol={{xs:24}} hasFeedback>
                        {getFieldDecorator('EndTime', {
                            rules: [{ required: true, message: uiStore.T("DATE_END_TIME_NOT_GIVEN") }]
                        })(
                            <TimePicker
                              disabled={!this.editable}
                              placeholder={uiStore.T("DATE_END_TIME_PLACEHOLDER")}
                              format={"HH:mm"}
                              />
                        )}
                  </FormItem>
                </Col>
              </Row>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("DATE_CATEGORY")} hasFeedback>
                    {getFieldDecorator('CategoryID', {})(
                        <Select disabled={!this.editable}>
                          {dateCategories}
                        </Select>
                    )}
              </FormItem>

              {this.renderFeedbacks()}

              {this.renderButtons()}
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
      StartDate: Form.createFormField({value: moment(dt.StartTime || new Date(), 'YYYY-MM-DD')}),
      StartTime: Form.createFormField({value: moment(dt.StartTime || new Date(), 'HH:mm')}),
      EndDate: Form.createFormField({value: moment(dt.EndTime || new Date(), 'YYYY-MM-DD')}),
      EndTime: Form.createFormField({value: moment(dt.EndTime || new Date(), 'HH:mm')}),
      CategoryID: Form.createFormField({value: dt.CategoryID || fallbackID})
    }
  }
})(_DateForm);
export {DateForm};