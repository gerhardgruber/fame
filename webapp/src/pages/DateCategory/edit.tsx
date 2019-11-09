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
import { DateCategory } from '../../stores/DateCategoryStore/DateCategory';
import DateCategoryStore from '../../stores/DateCategoryStore';
import { DateCategoryForm } from '../../components/DateCategoryForm';

interface EditDateCategoryProps {
  dateCategoryID?: number;
  form: WrappedFormUtils;
}

interface EditDateCategoryState {
  dateCategory: DateCategory;
}

const dateCategoryStore = DateCategoryStore.getInstance();

@observer
class _EditDateCategory extends Page<EditDateCategoryProps, EditDateCategoryState> {
  state = {
    dateCategory: null
  }

  componentWillMount() {
    if ( this.props.dateCategoryID) {
      dateCategoryStore.loadDateCategory(this.props.dateCategoryID).then((dc: DateCategory) => {
        this.setState({
          dateCategory: dc
        });
      })
    }
  }

  pageTitle(): string {
    if (this.state.dateCategory) {
      return 'DATE_CATEGORIES_EDIT_DATE_CATEGORY';
    } else {
      return 'DATE_CATEGORIES_NEW_DATE_CATEGORY';
    }
  }

  renderContent(): JSX.Element {
    if (this.props.dateCategoryID && this.state.dateCategory) {
      return <DateCategoryForm dateCategory={this.state.dateCategory} />;
    } else if (!this.props.dateCategoryID) {
      return <DateCategoryForm />
    } else {
      return null;
    }
  }
}

const EditDateCategory = Form.create<EditDateCategoryProps>()(_EditDateCategory);
export {EditDateCategory};