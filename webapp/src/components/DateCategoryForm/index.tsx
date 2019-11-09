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
import DateCategoryStore from '../../stores/DateCategoryStore';
import { DateCategory } from '../../stores/DateCategoryStore/DateCategory';

interface DateCategoryFormProps {
  dateCategory?: DateCategory;
  form: WrappedFormUtils;
}

const uiStore = UiStore.getInstance();
const dateCategoryStore = DateCategoryStore.getInstance();

@observer
class _DateCategoryForm extends React.Component<DateCategoryFormProps> {
  state = {
    gotoDateCategories: false
  };

  save = (e) => {
    e.preventDefault();

    this.props.form.validateFields((err, data) => {
      if (err) {
        return
      }

      if (this.props.dateCategory) {
        this.props.dateCategory.setData(data);
        dateCategoryStore.saveDateCategory(this.props.dateCategory).then( () => {
          dateCategoryStore.loadDateCategories().then( () => {
            this.setState({
              gotoDateCategories: true
            });
          } );
        });
      } else {
        dateCategoryStore.createDateCategory(new DateCategory(data)).then( () => {
          dateCategoryStore.loadDateCategories().then( () => {
            this.setState({
              gotoDateCategories: true
            });
          } );
        });
      }
    });
  }

  deleteDateCategory = (e) => {
    e.preventDefault();
    dateCategoryStore.deleteDateCategory(this.props.dateCategory).then( () => {
      dateCategoryStore.loadDateCategories().then( () => {
        this.setState({
          gotoDateCategories: true
        });
      } );
    });
  }

  render(): JSX.Element {
    const { getFieldDecorator } = this.props.form;

    let deleteButton = null;
    if (this.props.dateCategory) {
      deleteButton = <div style={{"display": "inline-block", "marginRight": "1rem"}}>
        <Button onClick={this.deleteDateCategory} type="danger">
          {uiStore.T('DELETE')}
        </Button>
      </div>
    }

    let gotoDateCategories = null;
    if (this.state.gotoDateCategories) {
      gotoDateCategories = <Redirect to="/date_categories" />;
    }

    return  <Form onSubmit={this.save}>
              <FormItem {...uiStore.formItemLayout} label={uiStore.T("DATE_CATEGORY_NAME")} hasFeedback>
                    {getFieldDecorator('Name', {
                        rules: [{ required: true, message: uiStore.T("DATE_CATEGORY_NAME_NOT_GIVEN") }]
                    })(
                        <Input placeholder={uiStore.T("DATE_CATEGORY_NAME_PLACEHOLDER")} />
                    )}
              </FormItem>

              {deleteButton}
              <div style={{"display": "inline-block", "marginRight": "1rem"}}>
                <Link to="/date_categories"><Button>
                  {uiStore.T('CANCEL')}
                </Button></Link>
              </div>
              <div style={{"display": "inline-block"}}>
                <Button htmlType="submit" type="primary">
                  {uiStore.T('SAVE')}
                </Button>
              </div>
              {gotoDateCategories}
            </Form>
  }
}

const DateCategoryForm = Form.create({
  mapPropsToFields(props: DateCategoryFormProps) {
    const dc = props.dateCategory;
    if(!dc) return {};

    return {
      Name: Form.createFormField({value: dc.Name}),
    }
  }
})(_DateCategoryForm);
export {DateCategoryForm};