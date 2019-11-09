import * as React from 'react';
import {observer} from 'mobx-react';
import Page from "../../components/Page";
import { Form } from "antd";
import { WrappedFormUtils } from "antd/lib/form/Form";
import DateStore from '../../stores/DateStore';
import { DateForm } from '../../components/DateForm';
import { DateModel } from '../../stores/DateStore/Date';

interface EditDateProps {
  dateID?: number;
  form: WrappedFormUtils;
}

interface EditDateState {
  date: DateModel;
}

const dateStore = DateStore.getInstance();

@observer
class _EditDate extends Page<EditDateProps, EditDateState> {
  state = {
    date: null
  }

  componentWillMount() {
    if ( this.props.dateID) {
      dateStore.loadDate(this.props.dateID).then((dt: DateModel) => {
        this.setState({
          date: dt
        });
      })
    }
  }

  pageTitle(): string {
    if (this.state.date) {
      return 'DATES_EDIT_DATE';
    } else {
      return 'DATES_NEW_DATE';
    }
  }

  renderContent(): JSX.Element {
    if (this.props.dateID && this.state.date) {
      return <DateForm date={this.state.date} />;
    } else if (!this.props.dateID) {
      return <DateForm />
    } else {
      return null;
    }
  }
}

const EditDate = Form.create<EditDateProps>()(_EditDate);
export {EditDate};